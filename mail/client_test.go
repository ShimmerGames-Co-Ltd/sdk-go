package mail

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

func TestNewNilClient(t *testing.T) {
	t.Parallel()
	if _, err := New(nil); err == nil {
		t.Fatal("空 Client 应失败")
	}
}

func TestSendPlayerPathAndSnakeCaseBody(t *testing.T) {
	t.Parallel()
	var gotPath, gotBody string
	c := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.Method != http.MethodPost {
			t.Errorf("method=%s", r.Method)
		}
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"data":{"mail_ref_id":"m1","badge":{"unread_mail_count":1}}}`))
	})
	cli, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	rep, err := cli.SendPlayer(context.Background(), SendPlayerRequest{
		ServerID: "s1", To: "p1", Serial: "ser", Title: "t", Content: "c", From: "sys",
		ContentMode: ContentModePlatformTemplate, TemplateID: "welcome", TemplateArgs: `{"name":"a"}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != pathSendPlayer {
		t.Fatalf("path=%s", gotPath)
	}
	if !strings.Contains(gotBody, `"server_id"`) || strings.Contains(gotBody, `"serverID"`) {
		t.Fatalf("请求应为 snake_case: %s", gotBody)
	}
	if !strings.Contains(gotBody, `"content_mode":"platform_template"`) || !strings.Contains(gotBody, `"template_id":"welcome"`) {
		t.Fatalf("应含模板字段: %s", gotBody)
	}
	if rep.MailRefID != "m1" || rep.Badge.UnreadMailCount != 1 {
		t.Fatalf("reply=%+v", rep)
	}
}

func TestSyncAPIErrorVsHTTPError(t *testing.T) {
	t.Parallel()
	c := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":1003,"message":"bad"}`))
	})
	cli, err := New(c)
	if err != nil {
		t.Fatal(err)
	}
	_, err = cli.Sync(context.Background(), SyncRequest{PlayerID: "p", ServerID: "1", PlayerAttrs: map[string]any{"register_at": 1}})
	api, ok := core.AsAPIError(err)
	if !ok || api.Code != 1003 {
		t.Fatalf("want APIError 1003 got %v", err)
	}

	c2 := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`oops`))
	})
	cli2, _ := New(c2)
	_, err = cli2.Login(context.Background(), LoginRequest{PlayerID: "p", ServerID: "1", Lang: "zh"})
	if _, ok := core.AsHTTPError(err); !ok {
		t.Fatalf("want HTTPError got %v", err)
	}
}

func TestListDecodesSnakeCase(t *testing.T) {
	t.Parallel()
	c := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pathList {
			t.Errorf("path=%s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{
				"total": 1,
				"headers": []map[string]any{{
					"id": 9, "have_attachment": true, "template_id": "tpl",
				}},
			},
		})
	})
	cli, _ := New(c)
	rep, err := cli.List(context.Background(), ListRequest{PlayerID: "p", ServerID: "1", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if rep.Total != 1 || len(rep.Headers) != 1 || !rep.Headers[0].HaveAttachment {
		t.Fatalf("%+v", rep)
	}
}
