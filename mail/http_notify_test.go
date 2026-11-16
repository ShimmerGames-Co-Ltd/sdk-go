package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const notifyTestSecret = "test-notify-auth-secret"

func signedNotifyRequest(body []byte) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/project/mail-callback", bytes.NewReader(body))
	request.Header.Set(HeaderNotifyAppID, "app")
	request.Header.Set(HeaderNotifySignature, core.SignMessage(notifyTestSecret, string(body)))
	return request
}

func TestParseHTTPNotifyPreservesRawAndInt64(t *testing.T) {
	raw := []byte("{\n\"player_mail\":{\"server_id\":\"S1\",\"player_id\":\"player\",\"headers\":[{\"id\":9007199254740993,\"title\":\"邮件\",\"type\":3}]},\"server_mail\":null\n}")
	parsed, err := ParseHTTPNotify(signedNotifyRequest(raw), "app", notifyTestSecret)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.AppID != "app" || !bytes.Equal(parsed.RawBody, raw) || parsed.PlayerMail == nil || parsed.ServerMail != nil || parsed.PlayerMail.Headers[0].ID != 9007199254740993 {
		t.Fatalf("notification changed: %+v", parsed)
	}
}

func TestHTTPNotifyHandlerSuccessAndBusinessFailure(t *testing.T) {
	raw := []byte(`{"server_mail":{"server_id":"S1","headers":[{"id":7}]}}`)
	calls := 0
	receiver, err := NewHTTPNotifyHandler("app", notifyTestSecret, func(ctx context.Context, r *http.Request, notification HTTPNotifyRequest) error {
		calls++
		if notification.ServerMail.ServerID != "S1" || r.URL.Path != "/project/mail-callback" || !bytes.Equal(notification.RawBody, raw) {
			t.Fatal("callback lost payload or request")
		}
		return ctx.Err()
	})
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	receiver.ServeHTTP(response, signedNotifyRequest(raw))
	if calls != 1 || response.Code != http.StatusNoContent || response.Body.Len() != 0 {
		t.Fatalf("expected one callback and empty204, calls=%d status=%d", calls, response.Code)
	}
	failed, err := NewHTTPNotifyHandler("app", notifyTestSecret, func(context.Context, *http.Request, HTTPNotifyRequest) error {
		return errors.New("password=must-not-leak")
	})
	if err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	failed.ServeHTTP(response, signedNotifyRequest(raw))
	if response.Code != http.StatusServiceUnavailable || strings.Contains(response.Body.String(), "password") {
		t.Fatalf("unsafe business error response: %d %s", response.Code, response.Body.String())
	}
}

func TestHTTPNotifyRejectsBeforeBusinessCallback(t *testing.T) {
	server := `{"server_mail":{"server_id":"S1","headers":[{"id":1}]}}`
	tests := []struct {
		name      string
		body      string
		method    string
		appID     string
		signature string
		status    int
	}{
		{
			name:   "wrong method",
			body:   server,
			method: http.MethodGet,
			appID:  "app",
			status: http.StatusMethodNotAllowed,
		},
		{
			name:   "wrong app",
			body:   server,
			appID:  "other",
			status: http.StatusUnauthorized,
		},
		{
			name:      "missing signature",
			body:      server,
			appID:     "app",
			signature: "missing",
			status:    http.StatusUnauthorized,
		},
		{
			name:      "bad signature",
			body:      server,
			appID:     "app",
			signature: "bad",
			status:    http.StatusUnauthorized,
		},
		{
			name:   "bad JSON with valid signature",
			body:   `{`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "trailing JSON with valid signature",
			body:   server + ` {}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "neither payload",
			body:   `{"player_mail":null,"server_mail":null}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "both payloads",
			body:   `{"player_mail":{"server_id":"S1","player_id":"p","headers":[{"id":1}]},"server_mail":{"server_id":"S1","headers":[{"id":1}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "wildcard server",
			body:   `{"server_mail":{"server_id":"*","headers":[{"id":1}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "missing player",
			body:   `{"player_mail":{"server_id":"S1","headers":[{"id":1}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "empty headers",
			body:   `{"server_mail":{"server_id":"S1","headers":[]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "multiple player mail headers",
			body:   `{"player_mail":{"server_id":"S1","player_id":"p","headers":[{"id":1},{"id":2}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "multiple server mail headers",
			body:   `{"server_mail":{"server_id":"S1","headers":[{"id":1},{"id":2}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "blank server",
			body:   `{"server_mail":{"server_id":" ","headers":[{"id":1}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "blank player",
			body:   `{"player_mail":{"server_id":"S1","player_id":" ","headers":[{"id":1}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "null server mail header",
			body:   `{"server_mail":{"server_id":"S1","headers":[null]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "empty player mail header",
			body:   `{"player_mail":{"server_id":"S1","player_id":"p","headers":[{}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "negative mail header id",
			body:   `{"server_mail":{"server_id":"S1","headers":[{"id":-1}]}}`,
			appID:  "app",
			status: http.StatusBadRequest,
		},
		{
			name:   "over size limit",
			body:   strings.Repeat("x", HTTPNotifyMaxBodyBytes+1),
			appID:  "app",
			status: http.StatusRequestEntityTooLarge,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			receiver, err := NewHTTPNotifyHandler("app", notifyTestSecret, func(context.Context, *http.Request, HTTPNotifyRequest) error {
				calls++
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
			request := signedNotifyRequest([]byte(test.body))
			if test.method != "" {
				request.Method = test.method
			}
			request.Header.Set(HeaderNotifyAppID, test.appID)
			if test.signature == "missing" {
				request.Header.Del(HeaderNotifySignature)
			} else if test.signature != "" {
				request.Header.Set(HeaderNotifySignature, test.signature)
			}
			response := httptest.NewRecorder()
			receiver.ServeHTTP(response, request)
			if response.Code != test.status || calls != 0 {
				t.Fatalf("status=%d want=%d calls=%d", response.Code, test.status, calls)
			}
			if test.status == http.StatusMethodNotAllowed && response.Header().Get("Allow") != http.MethodPost {
				t.Fatal("missing Allow: POST")
			}
		})
	}
	raw := []byte(server)
	request := signedNotifyRequest(raw)
	request.Body = io.NopCloser(strings.NewReader(server + " "))
	if _, err := ParseHTTPNotify(request, "app", notifyTestSecret); err == nil {
		t.Fatal("whitespace tampering accepted")
	}
}

func TestHTTPNotifyConfigurationAndTypedError(t *testing.T) {
	callback := func(context.Context, *http.Request, HTTPNotifyRequest) error {
		return nil
	}
	for _, values := range [][2]string{{"", notifyTestSecret}, {"app", ""}} {
		if _, err := NewHTTPNotifyHandler(values[0], values[1], callback); err == nil {
			t.Fatal("missing credentials accepted")
		}
	}
	if _, err := NewHTTPNotifyHandler("app", notifyTestSecret, nil); err == nil {
		t.Fatal("nil callback accepted")
	}
	_, err := ParseHTTPNotify(signedNotifyRequest([]byte(`{}`)), "app", notifyTestSecret)
	var protocol *HTTPNotifyError
	if !errors.As(err, &protocol) || protocol.StatusCode != http.StatusBadRequest {
		t.Fatalf("typed request error missing: %v", err)
	}
}

func TestSingleHTTPEventKeepsLegacyNotificationArrayDecoding(t *testing.T) {
	raw := []byte(`{"server_mail":{"server_id":"S1","headers":[{"id":1},{"id":2}]}}`)
	var legacy NewMailNotify
	if err := json.Unmarshal(raw, &legacy); err != nil {
		t.Fatal(err)
	}
	if legacy.ServerMail == nil || len(legacy.ServerMail.Headers) != 2 {
		t.Fatal("legacy notification array decoding changed")
	}
	_, err := ParseHTTPNotify(signedNotifyRequest(raw), "app", notifyTestSecret)
	var protocol *HTTPNotifyError
	if !errors.As(err, &protocol) || protocol.StatusCode != http.StatusBadRequest {
		t.Fatalf("HTTP batch was not rejected: %v", err)
	}
}
