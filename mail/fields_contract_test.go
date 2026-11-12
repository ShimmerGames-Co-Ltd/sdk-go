package mail

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

// 完整线格式夹具保护 List/Get/通知的全部公开字段，包括 number/string 两种 int64。
func TestReadAndNotifyCompleteFields(t *testing.T) {
	t.Parallel()
	headerJSON := `{"id":"7","type":3,"from":"GM","title":"Title","sub_title":"Subtitle","content_category":"event","content_format":"plain_text","template_id":"client","have_attachment":true,"status":{"read_at":"11","claim_at":12},"create_at":100,"expire_at":"200","expire_cd":100}`
	detailJSON := `{"id":7,"type":3,"from":"GM","title":"Title","sub_title":"Subtitle","content":"Body","content_category":"event","content_format":"plain_text","template_id":"client","template_args":"{\"amount\":7}","attachment":"[{\"type\":2,\"id\":\"gold\",\"count\":7,\"param\":\"opaque\"}]","create_at":"100","expire_at":200,"expire_cd":"100"}`
	wantHeader := MailHeader{
		ID:              7,
		Type:            MailType(3),
		From:            "GM",
		Title:           "Title",
		SubTitle:        "Subtitle",
		ContentCategory: "event",
		ContentFormat:   "plain_text",
		TemplateID:      "client",
		HaveAttachment:  true,
		Status: MailStatus{
			ReadAt:  11,
			ClaimAt: 12,
		},
		CreateAt: 100,
		ExpireAt: 200,
		ExpireCD: 100,
	}
	wantDetail := MailDetail{
		ID:              7,
		Type:            MailType(3),
		From:            "GM",
		Title:           "Title",
		SubTitle:        "Subtitle",
		Content:         "Body",
		ContentCategory: "event",
		ContentFormat:   "plain_text",
		TemplateID:      "client",
		TemplateArgs:    `{"amount":7}`,
		Attachment:      `[{"type":2,"id":"gold","count":7,"param":"opaque"}]`,
		CreateAt:        100,
		ExpireAt:        200,
		ExpireCD:        100,
	}
	core := testCore(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case pathList:
			_, _ = w.Write([]byte(`{"code":0,"data":{"headers":[` + headerJSON + `],"total":1,"badge":{"pending_new_mail_count":2,"unread_mail_count":3}}}`))
		case pathGet:
			_, _ = w.Write([]byte(`{"code":0,"data":{"mails":[` + detailJSON + `]}}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	client, err := New(core)
	if err != nil {
		t.Fatal(err)
	}
	list, err := client.List(context.Background(), ListRequest{
		ServerID: "s1",
		PlayerID: "p1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if list.Total != 1 || !reflect.DeepEqual(list.Headers, []MailHeader{wantHeader}) || list.Badge.PendingNewMailCount != 2 || list.Badge.UnreadMailCount != 3 {
		t.Fatalf("incomplete list: %+v", list)
	}
	get, err := client.Get(context.Background(), GetRequest{
		ServerID: "s1",
		PlayerID: "p1",
		MailIDs:  []int64{7},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(get.Mails, []MailDetail{wantDetail}) {
		t.Fatalf("incomplete get: %+v", get)
	}
	var notify NewMailNotify
	if err := json.Unmarshal([]byte(`{"player_mail":{"server_id":"s1","player_id":"p1","headers":[`+headerJSON+`]}}`), &notify); err != nil {
		t.Fatal(err)
	}
	if notify.PlayerMail == nil || notify.PlayerMail.ServerID != "s1" || notify.PlayerMail.PlayerID != "p1" || !reflect.DeepEqual(notify.PlayerMail.Headers, []MailHeader{wantHeader}) {
		t.Fatalf("incomplete notify: %+v", notify)
	}
	var legacy MailHeader
	if err := json.Unmarshal([]byte(`{"id":7,"title":"Legacy"}`), &legacy); err != nil {
		t.Fatal(err)
	}
	if legacy.Title != "Legacy" || legacy.ContentFormat != "" {
		t.Fatalf("legacy defaults changed: %+v", legacy)
	}
}
