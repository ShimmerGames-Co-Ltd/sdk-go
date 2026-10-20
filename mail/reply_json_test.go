package mail

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestReplyIDLists(t *testing.T) {
	for _, field := range []string{"read_mail_ids", "removed_mail_ids"} {
		for _, input := range []string{`["9007199254740993",102]`, `[]`, `null`} {
			t.Run(field+input, func(t *testing.T) {
				data := []byte(`{"` + field + `":` + input + `,"badge":{"unread_mail_count":3},"unknown":1}`)
				var out any
				var ids *[]int64
				var badge *MailboxBadge
				if field == "read_mail_ids" {
					r := &ReadReply{}
					out, ids, badge = r, &r.ReadMailIDs, &r.Badge
				} else {
					r := &RemoveReply{}
					out, ids, badge = r, &r.RemovedMailIDs, &r.Badge
				}
				if err := json.Unmarshal(data, out); err != nil {
					t.Fatal(err)
				}
				var want []int64
				if input == "[]" {
					want = []int64{}
				} else if input != "null" {
					want = []int64{9007199254740993, 102}
				}
				if !reflect.DeepEqual(*ids, want) || badge.UnreadMailCount != 3 {
					t.Fatalf("ids=%#v badge=%+v", *ids, badge)
				}
				raw, err := json.Marshal(out)
				if err != nil {
					t.Fatal(err)
				}
				var encoded map[string]json.RawMessage
				if err := json.Unmarshal(raw, &encoded); err != nil {
					t.Fatal(err)
				}
				wantJSON, _ := json.Marshal(want)
				if string(encoded[field]) != string(wantJSON) {
					t.Fatalf("encode=%s want %s", encoded[field], wantJSON)
				}
				// 缺失字段和根 null 沿用 encoding/json 的保留已有值语义。
				for _, update := range []string{`{}`, `null`} {
					if err := json.Unmarshal([]byte(update), out); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(*ids, want) || badge.UnreadMailCount != 3 {
						t.Fatalf("%s modified reply: %+v", update, out)
					}
				}
				if err := json.Unmarshal([]byte(`{"`+field+`":["bad"]}`), out); err == nil {
					t.Fatal("accepted invalid ID")
				}
			})
		}
	}
}
