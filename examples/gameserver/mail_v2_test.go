package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

// 执行用户会运行的示例代码，预期 JSON 独立书写，避免示例与测试各维护一套请求构造。
func TestMailV2Examples(t *testing.T) {
	cases := []struct {
		demo     string
		path     string
		body     string
		response string
	}{
		{
			demo:     "mail-send-v2",
			path:     "/mail/v2/mail/server/send/player",
			body:     `{"from":"system","serial":"business-1","server_id":"1","player_id":"10001","content_mode":"plain","content_format":"plain_text","content_category":"welcome","i18n":[{"lang":"English","title":"Welcome","content":"Welcome to the game.","sub_title":"A new adventure"}]}`,
			response: `{"mail_ref_id":"player-ref","badge":{"pending_new_mail_count":1,"unread_mail_count":2}}`,
		},
		{
			demo:     "mail-group-v2",
			path:     "/mail/v2/mail/server/send/group",
			body:     `{"from":"system","serial":"business-1","recipient_targets":[{"server_id":"1","player_ids":["10001"]},{"server_id":"2","player_ids":["20001"]}],"i18n":[{"lang":"English","title":"Reward","content":"Your reward is ready."}],"attachments":[{"type":0,"id":"gold","count":100,"param":"event-reward"}]}`,
			response: `{"refs":[{"server_id":"1","player_id":"10001","mail_ref_id":"group-1","badge":{"unread_mail_count":1}},{"server_id":"2","player_id":"20001","mail_ref_id":"group-2","badge":{"unread_mail_count":2}}]}`,
		},
		{
			demo:     "mail-server-v2",
			path:     "/mail/v2/mail/server/send/server",
			body:     `{"from":"system","serial":"business-1","server_ids":["*"],"include_late_registrants":true,"i18n":[{"lang":"English","title":"Announcement","content":"Welcome, all players."}]}`,
			response: `{"refs":[{"server_id":"*","mail_ref_id":"all-ref"}]}`,
		},
		{
			demo:     "mail-audience-v2",
			path:     "/mail/v2/mail/server/send/server",
			body:     `{"from":"system","serial":"business-1","server_ids":["1"],"audience":{"op":"and","children":[{"key":"register_at","op":"range","value":{}}]},"i18n":[{"lang":"English","title":"New player gift","content":"A gift for players registered in the last seven days."}]}`,
			response: `{"refs":[{"server_id":"1","mail_ref_id":"conditional-ref"}]}`,
		},
		{
			demo:     "mail-platform-template-v2",
			path:     "/mail/v2/mail/server/send/server",
			body:     `{"from":"system","serial":"business-1","server_ids":["1"],"content_mode":"platform_template","template_id":"welcome_reward","include_late_registrants":true,"template_args":{"plain.player_name":"Ada","i18n.item_name":{"English":"Legendary Sword"}}}`,
			response: `{"refs":[{"server_id":"1","mail_ref_id":"platform-ref"}]}`,
		},
		{
			demo:     "mail-client-template-v2",
			path:     "/mail/v2/mail/server/send/player",
			body:     `{"from":"system","serial":"business-1","server_id":"1","player_id":"10001","content_mode":"client_template","content_format":"rich_text","template_id":"local_reward","template_args":{"amount":7,"item_id":"gold"}}`,
			response: `{"mail_ref_id":"client-ref","badge":{"unread_mail_count":1}}`,
		},
	}
	for _, tt := range cases {
		t.Run(tt.demo, func(t *testing.T) {
			setMailExampleEnv(t)
			started := time.Now().Unix()
			calls := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.URL.Path != tt.path || r.Method != http.MethodPost {
					t.Errorf("route=%s %s", r.Method, r.URL.Path)
				}
				var got map[string]any
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// 两处时间来自示例的当前时间，先独立断言，再比较其余完整请求。
				if tt.demo == "mail-send-v2" {
					assertExampleTime(t, got["expire_at"], started+86400, time.Now().Unix()+86400)
					delete(got, "expire_at")
				}
				if tt.demo == "mail-audience-v2" {
					audience, ok := got["audience"].(map[string]any)
					if !ok {
						t.Error("audience must be an object")
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					children, ok := audience["children"].([]any)
					if !ok || len(children) != 1 {
						t.Error("audience requires a rule")
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					rule := children[0].(map[string]any)
					bounds := rule["value"].(map[string]any)
					assertExampleTime(t, bounds["min"], started-7*86400, time.Now().Unix()-7*86400)
					delete(bounds, "min")
				}
				var want map[string]any
				if err := json.Unmarshal([]byte(tt.body), &want); err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("body=%v\nwant=%v", got, want)
				}
				_, _ = w.Write([]byte(`{"code":0,"data":` + tt.response + `}`))
			}))
			defer srv.Close()
			client, err := core.NewClient(
				core.WithURL(srv.URL),
				core.WithAppID("app1"),
				core.WithServerSignSecret("test-secret"),
			)
			if err != nil {
				t.Fatal(err)
			}
			if err := mailV2Demo(context.Background(), client, tt.demo); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatalf("calls=%d", calls)
			}
		})
	}
}

func TestMailReceiveExample(t *testing.T) {
	cases := []struct {
		name      string
		mails     string
		wantError bool
	}{
		{
			name:  "decode attachment and client parameters",
			mails: `[{"id":"7","content_format":"rich_text","template_id":"local_reward","template_args":"{\"amount\":7}","attachment":"[{\"type\":0,\"id\":\"gold\",\"count\":7}]"},{"id":8,"content_format":"","template_args":"","attachment":""}]`,
		},
		{
			name:      "attachment decoding failure",
			mails:     `[{"id":"7","attachment":"invalid json"}]`,
			wantError: true,
		},
		{
			name:      "template parameter decoding failure",
			mails:     `[{"id":"7","template_id":"local_reward","template_args":"invalid json"}]`,
			wantError: true,
		},
		{
			name:      "template amount must decode as an integer",
			mails:     `[{"id":"7","template_id":"local_reward","template_args":"{\"amount\":\"7\"}"}]`,
			wantError: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			setMailExampleEnv(t)
			var paths []string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				paths = append(paths, r.URL.Path)
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if body["server_id"] != "1" || body["player_id"] != "10001" {
					t.Errorf("mailbox identity=%v", body)
				}
				switch r.URL.Path {
				case "/mail/v1/mail/server/login":
					if body["lang"] != "English" {
						t.Errorf("lang=%v", body["lang"])
					}
					_, _ = w.Write([]byte(`{"code":0,"data":{"badge":{}}}`))
				case "/mail/v1/mail/server/sync":
					wantAttrs := map[string]any{"register_at": float64(1700000000)}
					if !reflect.DeepEqual(body["player_attrs"], wantAttrs) {
						t.Errorf("attrs=%v", body["player_attrs"])
					}
					_, _ = w.Write([]byte(`{"code":0,"data":{"new_mail_count":1,"badge":{"unread_mail_count":2}}}`))
				case "/mail/v1/mail/server/list":
					_, _ = w.Write([]byte(`{"code":0,"data":{"headers":[{"id":"7"},{"id":8}],"total":2,"badge":{}}}`))
				case "/mail/v1/mail/server/get":
					if !reflect.DeepEqual(body["mail_ids"], []any{float64(7), float64(8)}) {
						t.Errorf("Get must use mailbox IDs: %v", body["mail_ids"])
					}
					_, _ = w.Write([]byte(`{"code":0,"data":{"mails":` + tt.mails + `}}`))
				default:
					t.Errorf("unexpected path %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer srv.Close()
			client, err := core.NewClient(
				core.WithURL(srv.URL),
				core.WithAppID("app1"),
				core.WithServerSignSecret("test-secret"),
			)
			if err != nil {
				t.Fatal(err)
			}
			err = mailV2Demo(context.Background(), client, "mail-receive")
			if (err != nil) != tt.wantError {
				t.Fatalf("err=%v wantError=%v", err, tt.wantError)
			}
			wantPaths := []string{
				"/mail/v1/mail/server/login",
				"/mail/v1/mail/server/sync",
				"/mail/v1/mail/server/list",
				"/mail/v1/mail/server/get",
			}
			if !reflect.DeepEqual(paths, wantPaths) {
				t.Fatalf("paths=%v", paths)
			}
		})
	}
}

func TestMailExamplesRejectMissingInputsBeforeSending(t *testing.T) {
	cases := []struct {
		demo string
		key  string
	}{
		{
			demo: "mail-send-v2",
			key:  "SHIMO_MAIL_SERIAL",
		},
		{
			demo: "mail-send-v2",
			key:  "SHIMO_PLAYER_ID",
		},
		{
			demo: "mail-group-v2",
			key:  "SHIMO_SECOND_SERVER_ID",
		},
		{
			demo: "mail-platform-template-v2",
			key:  "SHIMO_TEMPLATE_CODE",
		},
		{
			demo: "mail-client-template-v2",
			key:  "SHIMO_CLIENT_TEMPLATE_ID",
		},
		{
			demo: "mail-receive",
			key:  "SHIMO_REGISTER_AT",
		},
	}
	for _, tt := range cases {
		t.Run(tt.demo+"/"+tt.key, func(t *testing.T) {
			setMailExampleEnv(t)
			t.Setenv(tt.key, "")
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("missing input must not send a request")
				w.WriteHeader(http.StatusBadRequest)
			}))
			defer srv.Close()
			client, err := core.NewClient(
				core.WithURL(srv.URL),
				core.WithAppID("app1"),
				core.WithServerSignSecret("test-secret"),
			)
			if err != nil {
				t.Fatal(err)
			}
			if _, ok := mailV2Demo(context.Background(), client, tt.demo).(*core.ConfigError); !ok {
				t.Fatal("missing input must return ConfigError")
			}
		})
	}
}

func setMailExampleEnv(t *testing.T) {
	t.Helper()
	t.Setenv("SHIMO_SERVER_ID", "1")
	t.Setenv("SHIMO_PLAYER_ID", "10001")
	t.Setenv("SHIMO_DEFAULT_LANG", "English")
	t.Setenv("SHIMO_MAIL_SERIAL", "business-1")
	t.Setenv("SHIMO_SECOND_SERVER_ID", "2")
	t.Setenv("SHIMO_SECOND_PLAYER_ID", "20001")
	t.Setenv("SHIMO_TEMPLATE_CODE", "welcome_reward")
	t.Setenv("SHIMO_CLIENT_TEMPLATE_ID", "local_reward")
	t.Setenv("SHIMO_REGISTER_AT", "1700000000")
}

func assertExampleTime(t *testing.T, value any, min, max int64) {
	t.Helper()
	got, ok := value.(float64)
	if !ok || got < float64(min) || got > float64(max) {
		t.Errorf("time=%v want [%d,%d]", value, min, max)
	}
}
