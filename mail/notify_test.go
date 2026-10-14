package mail

import (
	"encoding/json"
	"testing"
)

func TestNewMailNotify_unmarshalPlayerMailSnakeCase(t *testing.T) {
	t.Parallel()
	raw := `{
		"player_mail": {
			"server_id": "1",
			"player_id": "10001",
			"headers": [{
				"id": 7,
				"type": 3,
				"title": "hello",
				"have_attachment": true,
				"create_at": 1,
				"expire_at": 3601,
				"expire_cd": 3600,
				"sub_title": "sub",
				"content_category": "system"
			}]
		}
	}`
	var n NewMailNotify
	if err := json.Unmarshal([]byte(raw), &n); err != nil {
		t.Fatal(err)
	}
	if n.PlayerMail == nil || n.PlayerMail.ServerID != "1" {
		t.Fatalf("unexpected payload: %+v", n)
	}
	if len(n.PlayerMail.Headers) != 1 || !n.PlayerMail.Headers[0].HaveAttachment {
		t.Fatalf("unexpected headers: %+v", n.PlayerMail.Headers)
	}
}
