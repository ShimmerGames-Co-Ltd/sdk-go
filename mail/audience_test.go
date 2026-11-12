package mail

import (
	"encoding/json"
	"testing"
)

// 校验固定 AST 形状及动态值的线格式，不复制服务端的 Catalog 业务校验。
func TestAudienceJSONPreservesNestedGroupsAndValues(t *testing.T) {
	audience := Audience{
		Op: AudienceOpAnd,
		Children: []AudienceNode{
			{
				Key: "register_at",
				Op:  AudienceOpRange,
				Value: AudienceRange{
					Min: int64(0),
					Max: int64(1700000000),
				},
			},
			{
				Op: AudienceOpOr,
				Children: []AudienceNode{
					{
						Key:   "server_id",
						Op:    AudienceOpIn,
						Value: []string{"1", "2"},
					},
					{
						Key: "client_version",
						Op:  AudienceOpRange,
						Value: AudienceRange{
							Min: "1.2.3",
						},
					},
				},
			},
			{
				Key:   "custom_enabled",
				Op:    AudienceOpEq,
				Value: false,
			},
			{
				Key:   "custom_name",
				Op:    AudienceOpEq,
				Value: "",
			},
			{
				Key:   "custom_level",
				Op:    AudienceOpGte,
				Value: 0,
			},
		},
	}
	raw, err := json.Marshal(audience)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, raw, `{"op":"and","children":[{"key":"register_at","op":"range","value":{"min":0,"max":1700000000}},{"op":"or","children":[{"key":"server_id","op":"in","value":["1","2"]},{"key":"client_version","op":"range","value":{"min":"1.2.3"}}]},{"key":"custom_enabled","op":"eq","value":false},{"key":"custom_name","op":"eq","value":""},{"key":"custom_level","op":"gte","value":0}]}`)
}
