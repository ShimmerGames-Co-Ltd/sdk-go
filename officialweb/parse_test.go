package officialweb

import (
	"encoding/json"
	"testing"
)

func TestParseAddOrder(t *testing.T) {
	got, err := ParseAddOrder(map[string]any{
		"platform_role_id":  "123",
		"product_id":        "p1",
		"param":             []any{json.Number("1"), json.Number("2")},
		"product_name_i18n": "{\"English\":\"Pack\"}",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.PlatformRoleID != "123" || got.ProductID != "p1" || got.ProductNameI18n != "{\"English\":\"Pack\"}" {
		t.Fatalf("%+v", got)
	}
	if len(got.Param) != 2 || got.Param[0] != 1 || got.Param[1] != 2 {
		t.Fatalf("param=%v", got.Param)
	}

	_, err = ParseAddOrder(map[string]any{"platform_role_id": "1"})
	if err == nil {
		t.Fatal("want product_id required")
	}
	got, err = ParseAddOrder(map[string]any{
		"platform_role_id": json.Number("3747523271598286848"),
		"product_id":       "p1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.PlatformRoleID != "3747523271598286848" || got.ProductNameI18n != "" || got.Param != nil {
		t.Fatalf("%+v", got)
	}
}

func TestParseBuyPayment(t *testing.T) {
	got, err := ParseBuyPayment(map[string]any{
		"platform_role_id":  "9",
		"cp_order_id":       "oid",
		"platform_order_id": "ow-1",
		"cp_ext":            "",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.CpOrderID != "oid" || got.PlatformOrderID != "ow-1" || got.CpExt != "" {
		t.Fatalf("%+v", got)
	}
	_, err = ParseBuyPayment(map[string]any{"platform_role_id": "9"})
	if err == nil {
		t.Fatal("want cp_order_id")
	}
}

func TestParseLookupListPreCheck(t *testing.T) {
	lu, err := ParseLookupRole(map[string]any{"cp_role_id": "r1", "cp_server_id": "1"})
	if err != nil || lu.CpRoleID != "r1" || lu.CpServerID != "1" {
		t.Fatalf("%+v %v", lu, err)
	}
	_, err = ParseLookupRole(map[string]any{})
	if err == nil {
		t.Fatal("want cp_role_id")
	}
	lp, err := ParseListProducts(map[string]any{"platform_role_id": "8"})
	if err != nil || lp.PlatformRoleID != "8" {
		t.Fatalf("%+v %v", lp, err)
	}
	pc, err := ParsePreCheck(map[string]any{"platform_role_id": "8", "product_id": "p"})
	if err != nil || pc.Param != nil {
		t.Fatalf("%+v %v", pc, err)
	}
}

func TestParseAddOrder_BadProductNameType(t *testing.T) {
	_, err := ParseAddOrder(map[string]any{
		"platform_role_id":  "1",
		"product_id":        "p",
		"product_name_i18n": map[string]any{"English": "x"},
	})
	if err == nil {
		t.Fatal("want type error for nested object")
	}
}
