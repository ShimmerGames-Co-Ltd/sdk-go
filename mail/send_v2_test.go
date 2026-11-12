package mail

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

func TestV2MethodsBindIndependentPathsAndObjectBodies(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		call     func(*Client) (any, error)
		body     string
		response string
		want     any
	}{
		{
			name: "player platform template",
			path: "/mail/v2/mail/server/send/player",
			call: func(c *Client) (any, error) {
				return c.SendPlayerV2(context.Background(), SendPlayerV2Request{
					AppID:           "app1",
					From:            "system",
					Serial:          "player-1",
					ServerID:        "s1",
					PlayerID:        "p1",
					ContentCategory: "project_custom",
					ContentMode:     ContentModePlatformTemplate,
					TemplateID:      "welcome",
					ContentFormat:   ContentFormatPlainText,
					ExpireAt:        2100000000,
					TemplateArgs: map[string]any{
						"plain.name": "Ada",
						"i18n.item": map[string]any{
							"English":           "Sword",
							"ChineseSimplified": "剑",
						},
					},
					Attachments: []MailAttachment{
						{
							Type:  0,
							ID:    "gold",
							Count: 7,
							Param: "opaque",
						},
					},
				})
			},
			body:     `{"app_id":"app1","from":"system","serial":"player-1","server_id":"s1","player_id":"p1","content_category":"project_custom","content_mode":"platform_template","template_id":"welcome","content_format":"plain_text","expire_at":2100000000,"template_args":{"plain.name":"Ada","i18n.item":{"English":"Sword","ChineseSimplified":"剑"}},"attachments":[{"type":0,"id":"gold","count":7,"param":"opaque"}]}`,
			response: `{"mail_ref_id":"player-ref","badge":{"pending_new_mail_count":2,"unread_mail_count":3}}`,
			want: &SendPlayerV2Reply{
				MailRefID: "player-ref",
				Badge: MailboxBadge{
					PendingNewMailCount: 2,
					UnreadMailCount:     3,
				},
			},
		},
		{
			name: "group plain",
			path: "/mail/v2/mail/server/send/group",
			call: func(c *Client) (any, error) {
				return c.SendGroupV2(context.Background(), SendGroupV2Request{
					From:   "system",
					Serial: "group-1",
					RecipientTargets: []RecipientTarget{
						{
							ServerID:  "s1",
							PlayerIDs: []string{"p1", "p2"},
						},
						{
							ServerID:  "s2",
							PlayerIDs: []string{"p3"},
						},
					},
					I18n: []I18nBlock{
						{
							Lang:     "English",
							Title:    "Title",
							Content:  "Body",
							SubTitle: "Subtitle",
						},
						{
							Lang:    "ChineseSimplified",
							Title:   "标题",
							Content: "正文",
						},
					},
				})
			},
			body:     `{"from":"system","serial":"group-1","recipient_targets":[{"server_id":"s1","player_ids":["p1","p2"]},{"server_id":"s2","player_ids":["p3"]}],"i18n":[{"lang":"English","title":"Title","content":"Body","sub_title":"Subtitle"},{"lang":"ChineseSimplified","title":"标题","content":"正文"}]}`,
			response: `{"refs":[{"server_id":"s1","player_id":"p1","mail_ref_id":"shared-ref","badge":{"pending_new_mail_count":1,"unread_mail_count":2}},{"server_id":"s1","player_id":"p2","mail_ref_id":"shared-ref","badge":{"pending_new_mail_count":2,"unread_mail_count":3}},{"server_id":"s2","player_id":"p3","mail_ref_id":"s2-ref","badge":{"pending_new_mail_count":3,"unread_mail_count":4}}]}`,
			want: &SendGroupV2Reply{
				Refs: []GroupMailV2Ref{
					{
						ServerID:  "s1",
						PlayerID:  "p1",
						MailRefID: "shared-ref",
						Badge: MailboxBadge{
							PendingNewMailCount: 1,
							UnreadMailCount:     2,
						},
					},
					{
						ServerID:  "s1",
						PlayerID:  "p2",
						MailRefID: "shared-ref",
						Badge: MailboxBadge{
							PendingNewMailCount: 2,
							UnreadMailCount:     3,
						},
					},
					{
						ServerID:  "s2",
						PlayerID:  "p3",
						MailRefID: "s2-ref",
						Badge: MailboxBadge{
							PendingNewMailCount: 3,
							UnreadMailCount:     4,
						},
					},
				},
			},
		},
		{
			name: "server audience",
			path: "/mail/v2/mail/server/send/server",
			call: func(c *Client) (any, error) {
				return c.SendServerV2(context.Background(), SendServerV2Request{
					From:                   "system",
					Serial:                 "server-1",
					ServerIDs:              []string{"s1", "s2"},
					IncludeLateRegistrants: true,
					Audience: &Audience{
						Op: AudienceOpAnd,
						Children: []AudienceNode{
							{
								Key: "register_at",
								Op:  AudienceOpRange,
								Value: AudienceRange{
									Min: 0,
								},
							},
						},
					},
					I18n: []I18nBlock{
						{
							Lang:    "English",
							Title:   "Title",
							Content: "Body",
						},
					},
				})
			},
			body:     `{"from":"system","serial":"server-1","server_ids":["s1","s2"],"include_late_registrants":true,"audience":{"op":"and","children":[{"key":"register_at","op":"range","value":{"min":0}}]},"i18n":[{"lang":"English","title":"Title","content":"Body"}]}`,
			response: `{"refs":[{"server_id":"s1","mail_ref_id":"s1-ref"},{"server_id":"s2","mail_ref_id":"s2-ref"}]}`,
			want: &SendServerV2Reply{
				Refs: []ServerMailV2Ref{
					{
						ServerID:  "s1",
						MailRefID: "s1-ref",
					},
					{
						ServerID:  "s2",
						MailRefID: "s2-ref",
					},
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			clientCore := testCore(t, func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.URL.Path != tt.path || r.Method != http.MethodPost {
					t.Errorf("unexpected route %s %s", r.Method, r.URL.Path)
				}
				if r.Header.Get("Authorization") == "" {
					t.Error("missing server signature")
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				assertJSONEqual(t, body, tt.body)
				_, _ = w.Write([]byte(`{"code":0,"data":` + tt.response + `}`))
			})
			client, err := New(clientCore)
			if err != nil {
				t.Fatal(err)
			}
			got, err := tt.call(client)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("reply=%+v want=%+v", got, tt.want)
			}
			if calls != 1 {
				t.Fatalf("expected one call, got %d", calls)
			}
		})
	}
}

func TestV2ExpireAtUsesInt64AndPreservesJSONNumbers(t *testing.T) {
	t.Parallel()
	for _, expireAt := range []int64{0, 2100000000, 9007199254740993} {
		cases := []struct {
			name    string
			request any
		}{
			{
				name: "player",
				request: SendPlayerV2Request{
					ExpireAt: expireAt,
				},
			},
			{
				name: "group",
				request: SendGroupV2Request{
					ExpireAt: expireAt,
				},
			},
			{
				name: "server",
				request: SendServerV2Request{
					ExpireAt: expireAt,
				},
			},
		}
		for _, tt := range cases {
			t.Run(tt.name+"/"+strconv.FormatInt(expireAt, 10), func(t *testing.T) {
				raw, err := json.Marshal(tt.request)
				if err != nil {
					t.Fatal(err)
				}
				var body map[string]json.RawMessage
				if err := json.Unmarshal(raw, &body); err != nil {
					t.Fatal(err)
				}
				value, present := body["expire_at"]
				if expireAt == 0 {
					if present {
						t.Fatalf("zero ExpireAt must be omitted: %s", raw)
					}
					return
				}
				if string(value) != strconv.FormatInt(expireAt, 10) {
					t.Fatalf("expire_at=%s want JSON number %d", value, expireAt)
				}
			})
		}
	}
}

func TestSendServerV2OmittedAudienceAndDefaults(t *testing.T) {
	clientCore := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		assertJSONEqual(t, body, `{"from":"system","serial":"all-1","server_ids":["*"],"i18n":[{"lang":"English","title":"Title","content":"Body"}]}`)
		_, _ = w.Write([]byte(`{"code":0,"data":{"refs":[{"server_id":"*","mail_ref_id":"all-ref"}]}}`))
	})
	client, err := New(clientCore)
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.SendServerV2(context.Background(), SendServerV2Request{
		From:      "system",
		Serial:    "all-1",
		ServerIDs: []string{"*"},
		I18n: []I18nBlock{
			{
				Lang:    "English",
				Title:   "Title",
				Content: "Body",
			},
		},
	})
	if err != nil || len(got.Refs) != 1 || got.Refs[0].ServerID != "*" {
		t.Fatalf("reply=%+v err=%v", got, err)
	}
}

func TestV2BusinessFailureDoesNotDowngrade(t *testing.T) {
	cases := []struct {
		path string
		call func(*Client) (any, error)
	}{
		{
			path: "/mail/v2/mail/server/send/player",
			call: func(c *Client) (any, error) {
				return c.SendPlayerV2(context.Background(), SendPlayerV2Request{})
			},
		},
		{
			path: "/mail/v2/mail/server/send/group",
			call: func(c *Client) (any, error) {
				return c.SendGroupV2(context.Background(), SendGroupV2Request{})
			},
		},
		{
			path: "/mail/v2/mail/server/send/server",
			call: func(c *Client) (any, error) {
				return c.SendServerV2(context.Background(), SendServerV2Request{})
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.path, func(t *testing.T) {
			calls := 0
			clientCore := testCore(t, func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.URL.Path != tt.path {
					t.Errorf("unexpected fallback %s", r.URL.Path)
				}
				_, _ = w.Write([]byte(`{"code":1003,"message":"invalid request"}`))
			})
			client, err := New(clientCore)
			if err != nil {
				t.Fatal(err)
			}
			got, err := tt.call(client)
			api, ok := core.AsAPIError(err)
			if !ok || api.Code != 1003 || calls != 1 || !reflect.ValueOf(got).IsNil() {
				t.Fatalf("reply=%+v err=%v calls=%d", got, err, calls)
			}
		})
	}
}

func TestSendPlayerV2AllowsEmptyReferenceOnIdempotentSuccess(t *testing.T) {
	clientCore := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":0,"data":{"mail_ref_id":"","badge":{"unread_mail_count":0}}}`))
	})
	client, err := New(clientCore)
	if err != nil {
		t.Fatal(err)
	}
	got, err := client.SendPlayerV2(context.Background(), SendPlayerV2Request{
		From:        "system",
		Serial:      "existing",
		ServerID:    "s1",
		PlayerID:    "p1",
		ContentMode: ContentModeClientTemplate,
		TemplateID:  "local",
	})
	if err != nil || got == nil || got.MailRefID != "" {
		t.Fatalf("reply=%+v err=%v", got, err)
	}
}

func assertJSONEqual(t *testing.T, got []byte, want string) {
	t.Helper()
	var gotValue any
	var wantValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Error(err)
		return
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("body=%s\nwant=%s", got, want)
	}
}
