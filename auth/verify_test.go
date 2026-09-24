package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

func TestVerifySHA1BodySignAndAuthorizationHeader(t *testing.T) {
	t.Parallel()
	var gotBody verifyBody
	var gotAuth, gotApp string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v1/server/verify" || r.Method != http.MethodPost {
			t.Errorf("path=%s method=%s", r.URL.Path, r.Method)
		}
		gotAuth = r.Header.Get("Authorization")
		gotApp = r.Header.Get("X-App-ID")
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"message":"ok","data":{"user_id":1,"role_id":2,"open_id":10001}}`))
	}))
	t.Cleanup(srv.Close)

	corec, err := core.NewClient(
		core.WithURL(srv.URL),
		core.WithAppID("shim_test"),
		core.WithOrganizationID("10000000000"),
		core.WithServerSignSecret("test-auth-secret"),
	)
	if err != nil {
		t.Fatal(err)
	}
	cli, err := New(corec)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := cli.VerifyAt(context.Background(), "tok", 1735689600)
	if err != nil {
		t.Fatal(err)
	}
	if rep.OpenID != 10001 {
		t.Fatalf("open_id=%d", rep.OpenID)
	}
	wantSign := core.UserVerifySign("shim_test", "tok", "test-auth-secret", 1735689600)
	if gotBody.Sign != wantSign || gotBody.Token != "tok" || gotBody.Ts != 1735689600 {
		t.Fatalf("body=%+v want sign %s", gotBody, wantSign)
	}
	if wantSign != "bI5G%2BNdANAgwcvqvBGDTUHZe990%3D" {
		t.Fatalf("golden drifted: %s", wantSign)
	}
	if gotApp != "shim_test" || gotAuth == "" {
		t.Fatalf("app=%s auth=%q", gotApp, gotAuth)
	}
}

func TestNewNilClient(t *testing.T) {
	t.Parallel()
	if _, err := New(nil); err == nil {
		t.Fatal("空 Client 应失败")
	}
}
