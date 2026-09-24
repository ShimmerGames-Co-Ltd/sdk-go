package leaderboard

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

func testCore(t *testing.T, h http.HandlerFunc) *core.Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c, err := core.NewClient(core.WithURL(srv.URL), core.WithAppID("app1"), core.WithServerSignSecret("s"), core.WithOrganizationID("1"))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSetScorePathAndAuthHeader(t *testing.T) {
	t.Parallel()
	var gotPath, gotAuth, gotApp, gotBody string
	c := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotApp = r.Header.Get("X-App-ID")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{"member":{"uid":"u1","score":10,"rank":1}}}`))
	})
	cli, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := cli.SetScore(context.Background(), SetScoreRequest{ID: "lb1", UID: "u1", Score: 10, SubID: 2})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != pathSetScore || gotApp != "app1" || !strings.Contains(gotAuth, "signature=") {
		t.Fatalf("path=%s app=%s auth=%q", gotPath, gotApp, gotAuth)
	}
	if !strings.Contains(gotBody, `"subId":2`) {
		t.Fatalf("body=%s", gotBody)
	}
	if rep.Member.UID != "u1" || rep.Member.Score != 10 {
		t.Fatalf("%+v", rep)
	}
}

func TestHistoryMemberInfoAndAPIError(t *testing.T) {
	t.Parallel()
	c := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pathHistoryMemberInfo {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{"exists": true, "member": map[string]any{"uid": "u", "score": 3, "rank": 2}},
		})
	})
	cli, _ := New(c)
	rep, err := cli.HistoryMemberInfo(context.Background(), HistoryMemberInfoRequest{ID: "lb", SubID: 1, UID: "u"})
	if err != nil || !rep.Exists || rep.Member.Rank != 2 {
		t.Fatalf("rep=%+v err=%v", rep, err)
	}

	c2 := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":40401,"message":"not found"}`))
	})
	cli2, _ := New(c2)
	_, err = cli2.Ask(context.Background(), AskRequest{ID: "missing"})
	api, ok := core.AsAPIError(err)
	if !ok || api.Code != 40401 {
		t.Fatalf("want APIError got %v", err)
	}
}

func TestNewNilClient(t *testing.T) {
	t.Parallel()
	if _, err := New(nil); err == nil {
		t.Fatal("空 Client 应失败")
	}
}
