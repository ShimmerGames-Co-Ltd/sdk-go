package officialweb

import (
	"encoding/json"
	"testing"
)

func TestAddOrderUnmarshal(t *testing.T) {
	raw, _ := json.Marshal(map[string]any{
		"platform_role_id":  "123",
		"product_id":        "p1",
		"param":             []any{1, 2},
		"product_name_i18n": `{"English":"Pack"}`,
		"sign":              "ignored-by-unmarshal",
	})
	var got AddOrderRequest
	if err := unmarshalBody(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.PlatformRoleID != "123" || got.ProductID != "p1" || got.ProductNameI18n != `{"English":"Pack"}` {
		t.Fatalf("%+v", got)
	}
	if len(got.Param) != 2 || got.Param[0] != 1 || got.Param[1] != 2 {
		t.Fatalf("param=%v", got.Param)
	}

	// V0：缺 product_id 仍 Unmarshal 成功，零值交给宿主。
	var miss AddOrderRequest
	if err := unmarshalBody([]byte(`{"platform_role_id":"1"}`), &miss); err != nil {
		t.Fatal(err)
	}
	if miss.ProductID != "" {
		t.Fatalf("%+v", miss)
	}

	var snow AddOrderRequest
	if err := unmarshalBody([]byte(`{"platform_role_id":"3747523271598286848","product_id":"p1"}`), &snow); err != nil {
		t.Fatal(err)
	}
	if snow.PlatformRoleID != "3747523271598286848" || snow.ProductNameI18n != "" || snow.Param != nil {
		t.Fatalf("%+v", snow)
	}
}

func TestBuyPaymentUnmarshal(t *testing.T) {
	var got BuyPaymentRequest
	if err := unmarshalBody([]byte(`{"platform_role_id":"9","cp_order_id":"oid","platform_order_id":"ow-1","cp_ext":""}`), &got); err != nil {
		t.Fatal(err)
	}
	if got.CpOrderID != "oid" || got.PlatformOrderID != "ow-1" || got.CpExt != "" {
		t.Fatalf("%+v", got)
	}
}

func TestLookupListPreCheckUnmarshal(t *testing.T) {
	var lu LookupRoleRequest
	if err := unmarshalBody([]byte(`{"cp_role_id":"r1","cp_server_id":"1"}`), &lu); err != nil {
		t.Fatal(err)
	}
	if lu.CpRoleID != "r1" || lu.CpServerID != "1" {
		t.Fatalf("%+v", lu)
	}
	var lp ListProductsRequest
	if err := unmarshalBody([]byte(`{"platform_role_id":"8"}`), &lp); err != nil {
		t.Fatal(err)
	}
	var pc PreCheckRequest
	if err := unmarshalBody([]byte(`{"platform_role_id":"8","product_id":"p"}`), &pc); err != nil {
		t.Fatal(err)
	}
	if pc.Param != nil {
		t.Fatalf("param=%v", pc.Param)
	}
}

func TestAddOrderBadProductNameType(t *testing.T) {
	var got AddOrderRequest
	err := unmarshalBody([]byte(`{"platform_role_id":"1","product_id":"p","product_name_i18n":{"English":"x"}}`), &got)
	if err == nil {
		t.Fatal("want type error for nested object")
	}
}

func TestPlatformRoleIDRejectsJSONNumber(t *testing.T) {
	// 契约为十进制 string；number 留给验签 map 路径，struct 不解。
	var got ListProductsRequest
	err := unmarshalBody([]byte(`{"platform_role_id":3747523271598286848}`), &got)
	if err == nil {
		t.Fatal("want type error for JSON number")
	}
}
