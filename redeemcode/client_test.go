package redeemcode

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

func TestRedeemAndCheckPaths(t *testing.T) {
	t.Parallel()
	var paths []string
	var lastBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		b, _ := io.ReadAll(r.Body)
		lastBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{"result":0,"code":{"code":"ABC","reward":"gold"},"record":{"user_id":"u1","rewarded":true}}}`))
	}))
	t.Cleanup(srv.Close)
	corec, _ := core.NewClient(core.WithURL(srv.URL), core.WithAppID("app1"), core.WithAuthSecret("s"))
	cli, err := New(corec)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := cli.Redeem(context.Background(), RedeemRequest{Code: "ABC", UserIDStr: "u1", RoleID: 9})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(lastBody, `"user_id_str"`) || !strings.Contains(lastBody, `"role_id"`) {
		t.Fatalf("body=%s", lastBody)
	}
	if rep.Result != ResultOK || rep.Code == nil || rep.Code.Reward != "gold" {
		t.Fatalf("%+v", rep)
	}
	_, err = cli.Check(context.Background(), CheckRequest{Code: "ABC", UserID: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 || paths[0] != pathRedeem || paths[1] != pathCheck {
		t.Fatalf("paths=%v", paths)
	}
}
