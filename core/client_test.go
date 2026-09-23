package core

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClientRequiresURLAndAuth(t *testing.T) {
	t.Parallel()
	if _, err := NewClient(WithAppID("a"), WithAuthSecret("s")); err == nil {
		t.Fatal("缺 URL 应失败")
	}
	if _, err := NewClient(WithURL("http://127.0.0.1"), WithAppID("a")); err == nil {
		t.Fatal("缺鉴权应失败")
	}
	if _, err := NewClient(WithURL("http://127.0.0.1"), WithAppID("a"), WithAuthSecret("s"), WithUserToken("t")); err == nil {
		t.Fatal("双鉴权应失败")
	}
}

func TestDoJSONBusinessCodeIsAPIErrorNotHTTPError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":14003,"message":"账号已被禁言"}`))
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithURL(srv.URL), WithAppID("shim_test"), WithAuthSecret("secret"))
	if err != nil {
		t.Fatal(err)
	}
	err = c.DoJSON(context.Background(), http.MethodPost, "/v1/chat/messages", map[string]string{"content": "x"}, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Code != 14003 {
		t.Fatalf("want APIError 14003 got %v", err)
	}
	if _, ok := AsHTTPError(err); ok {
		t.Fatal("业务 code 不得映射成 HTTPError")
	}
}

func TestDoJSONNon200IsHTTPError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusBadGateway)
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(WithURL(srv.URL), WithAppID("shim_test"), WithUserToken("tok"))
	if err != nil {
		t.Fatal(err)
	}
	err = c.DoJSON(context.Background(), http.MethodGet, "/v1/chat/messages", nil, nil)
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("want HTTPError 502 got %v", err)
	}
	if _, ok := AsAPIError(err); ok {
		t.Fatal("非 200 不得当成 chat APIError")
	}
}

func TestDoJSONGETAndPOSTAndHeaders(t *testing.T) {
	t.Parallel()
	var gotMethod, gotAuth, gotToken, gotVer, gotApp string
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		gotToken = r.Header.Get("X-User-Token")
		gotVer = r.Header.Get("X-Server-Version")
		gotApp = r.Header.Get("X-App-ID")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"message":"success","data":{"ok":true}}`))
	}))
	t.Cleanup(srv.Close)

	fixedNow := time.Unix(1735689600, 0)
	gs, err := NewClient(
		WithURL(srv.URL),
		WithAppID("shim_test"),
		WithOrganizationID("10000000000"),
		WithAuthSecret("test-auth-secret"),
		withNow(func() time.Time { return fixedNow }),
		withNonce(func() (string, error) { return "n1", nil }),
	)
	if err != nil {
		t.Fatal(err)
	}
	var out struct {
		OK bool `json:"ok"`
	}
	body := map[string]string{"content": "x"}
	if err := gs.DoJSON(context.Background(), http.MethodPost, "/v1/chat/messages", body, &out); err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method=%s", gotMethod)
	}
	if gotToken != "" {
		t.Fatal("GameServer 不应带 X-User-Token")
	}
	if gotApp != "shim_test" || gotVer != DefaultServerVersion {
		t.Fatalf("app=%s ver=%s", gotApp, gotVer)
	}
	if !strings.Contains(gotAuth, `signature="`) || strings.Contains(strings.ToLower(gotAuth), "rsa") {
		t.Fatalf("bad authorization %q", gotAuth)
	}
	if !strings.HasSuffix(gotBody, "\n") {
		t.Fatalf("POST body 应含 Encoder 换行: %q", gotBody)
	}
	if !out.OK {
		t.Fatal("data 未解码")
	}

	player, err := NewClient(WithURL(srv.URL), WithAppID("shim_test"), WithUserToken("player-tok"))
	if err != nil {
		t.Fatal(err)
	}
	if err := player.DoJSON(context.Background(), http.MethodGet, "/v1/chat/messages?peer_id=b", nil, nil); err != nil {
		t.Fatal(err)
	}
	if gotMethod != http.MethodGet {
		t.Fatalf("GET method=%s", gotMethod)
	}
	if gotAuth != "" {
		t.Fatal("Player 不应带 Authorization")
	}
	if gotToken != "player-tok" {
		t.Fatalf("token=%s", gotToken)
	}
}

func TestDoJSONSuccessData(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success", "data": map[string]any{"seq": 12},
		})
	}))
	t.Cleanup(srv.Close)
	c, err := NewClient(WithURL(srv.URL), WithAppID("a"), WithUserToken("t"))
	if err != nil {
		t.Fatal(err)
	}
	var out struct {
		Seq int64 `json:"seq"`
	}
	if err := c.DoJSON(context.Background(), http.MethodGet, "/v1/chat/messages", nil, &out); err != nil {
		t.Fatal(err)
	}
	if out.Seq != 12 {
		t.Fatalf("seq=%d", out.Seq)
	}
}

func TestWithHTTPClientIsHonored(t *testing.T) {
	t.Parallel()
	called := false
	rt := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"code":0,"data":{}}`)),
			Header:     make(http.Header),
			Request:    r,
		}, nil
	})
	c, err := NewClient(
		WithURL("http://example.invalid"),
		WithAppID("a"),
		WithUserToken("t"),
		WithHTTPClient(&http.Client{Transport: rt}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.DoJSON(context.Background(), http.MethodGet, "/v1/x", nil, nil); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("WithHTTPClient 未被使用（旧 sdk-go WithHttpClient 有此 bug）")
	}
}

func TestPathPrefixUsedInSignURINotRequestPath(t *testing.T) {
	t.Parallel()
	var gotPath, gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{}}`))
	}))
	t.Cleanup(srv.Close)
	fixed := time.Unix(1735689600, 0)
	c, err := NewClient(
		WithURL(srv.URL),
		WithAppID("shim_test"),
		WithOrganizationID("10000000000"),
		WithAuthSecret("test-auth-secret"),
		WithPathPrefix("/chat"),
		withNow(func() time.Time { return fixed }),
		withNonce(func() (string, error) { return "n1", nil }),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.DoJSON(context.Background(), http.MethodGet, "/v1/chat/messages", nil, nil); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/v1/chat/messages" {
		t.Fatalf("HTTP path 不应被 PathPrefix 改写: %s", gotPath)
	}
	want := SignAuthorization("10000000000", "shim_test", "test-auth-secret", "/chat/v1/chat/messages", "", 1735689600, "n1")
	if gotAuth != want {
		t.Fatalf("sign URI 应对齐 PathPrefix\n got %s\nwant %s", gotAuth, want)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
