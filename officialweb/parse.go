package officialweb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// LookupRoleRequest POST /payment/lookup_role（验签后解析）。
type LookupRoleRequest struct {
	CpRoleID   string
	CpServerID string // 可选
}

// ListProductsRequest POST /payment/get_goods_list。
type ListProductsRequest struct {
	PlatformRoleID string // 线格式字符串（亦接受 JSON number）
}

// PreCheckRequest POST /payment/pre_check。
type PreCheckRequest struct {
	PlatformRoleID string
	ProductID      string
	Param          []int64 // 可选；缺省或省略为空 slice
}

// AddOrderRequest POST /payment/add_order。
type AddOrderRequest struct {
	PlatformRoleID  string
	ProductID       string
	Param           []int64 // 可选
	ProductNameI18n string  // 可选；线格式 JSON 字符串，不做二次 Unmarshal
}

// BuyPaymentRequest POST /payment/buypayment。
type BuyPaymentRequest struct {
	PlatformRoleID  string
	PlatformOrderID string // 可为空串
	CpOrderID       string
	CpExt           string // 可为空串
}

// ParseLookupRole 从验签后的 body 解析查角请求。
func ParseLookupRole(body map[string]any) (LookupRoleRequest, error) {
	var out LookupRoleRequest
	id, err := requireStringField(body, "cp_role_id")
	if err != nil {
		return out, err
	}
	out.CpRoleID = id
	out.CpServerID, _ = optionalStringField(body, "cp_server_id")
	return out, nil
}

// ParseListProducts 从验签后的 body 解析商品列表请求。
func ParseListProducts(body map[string]any) (ListProductsRequest, error) {
	var out ListProductsRequest
	id, err := requireStringField(body, "platform_role_id")
	if err != nil {
		return out, err
	}
	out.PlatformRoleID = id
	return out, nil
}

// ParsePreCheck 从验签后的 body 解析购前校验请求。
func ParsePreCheck(body map[string]any) (PreCheckRequest, error) {
	var out PreCheckRequest
	role, err := requireStringField(body, "platform_role_id")
	if err != nil {
		return out, err
	}
	pid, err := requireStringField(body, "product_id")
	if err != nil {
		return out, err
	}
	param, err := optionalInt64SliceField(body, "param")
	if err != nil {
		return out, err
	}
	out.PlatformRoleID = role
	out.ProductID = pid
	out.Param = param
	return out, nil
}

// ParseAddOrder 从验签后的 body 解析创单请求。
func ParseAddOrder(body map[string]any) (AddOrderRequest, error) {
	var out AddOrderRequest
	role, err := requireStringField(body, "platform_role_id")
	if err != nil {
		return out, err
	}
	pid, err := requireStringField(body, "product_id")
	if err != nil {
		return out, err
	}
	param, err := optionalInt64SliceField(body, "param")
	if err != nil {
		return out, err
	}
	nameI18n, err := optionalStringField(body, "product_name_i18n")
	if err != nil {
		return out, err
	}
	out.PlatformRoleID = role
	out.ProductID = pid
	out.Param = param
	out.ProductNameI18n = nameI18n
	return out, nil
}

// ParseBuyPayment 从验签后的 body 解析付后通知请求。
func ParseBuyPayment(body map[string]any) (BuyPaymentRequest, error) {
	var out BuyPaymentRequest
	role, err := requireStringField(body, "platform_role_id")
	if err != nil {
		return out, err
	}
	cpOID, err := requireStringField(body, "cp_order_id")
	if err != nil {
		return out, err
	}
	platOID, _ := optionalStringField(body, "platform_order_id")
	cpExt, err := optionalStringField(body, "cp_ext")
	if err != nil {
		return out, err
	}
	out.PlatformRoleID = role
	out.CpOrderID = cpOID
	out.PlatformOrderID = platOID
	out.CpExt = cpExt
	return out, nil
}

func requireStringField(body map[string]any, key string) (string, error) {
	v, ok := body[key]
	if !ok || v == nil {
		return "", fmt.Errorf("officialweb: missing %s", key)
	}
	s, ok := coerceString(v)
	if !ok {
		return "", fmt.Errorf("officialweb: %s must be string or number", key)
	}
	if s == "" {
		return "", fmt.Errorf("officialweb: missing %s", key)
	}
	return s, nil
}

func optionalStringField(body map[string]any, key string) (string, error) {
	v, ok := body[key]
	if !ok || v == nil {
		return "", nil
	}
	s, ok := coerceString(v)
	if !ok {
		return "", fmt.Errorf("officialweb: %s must be string or number", key)
	}
	return s, nil
}

func coerceString(v any) (string, bool) {
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x), true
	case json.Number:
		return strings.TrimSpace(x.String()), true
	case int64:
		return strconv.FormatInt(x, 10), true
	case int:
		return strconv.Itoa(x), true
	case float64:
		return strconv.FormatInt(int64(x), 10), true
	default:
		return "", false
	}
}

func optionalInt64SliceField(body map[string]any, key string) ([]int64, error) {
	v, ok := body[key]
	if !ok || v == nil {
		return nil, nil
	}
	switch x := v.(type) {
	case []int64:
		out := make([]int64, len(x))
		copy(out, x)
		return out, nil
	case []any:
		out := make([]int64, 0, len(x))
		for i, item := range x {
			n, err := coerceInt64(item)
			if err != nil {
				return nil, fmt.Errorf("officialweb: %s[%d]: %w", key, i, err)
			}
			out = append(out, n)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("officialweb: %s must be an array", key)
	}
}

func coerceInt64(v any) (int64, error) {
	switch x := v.(type) {
	case int64:
		return x, nil
	case int:
		return int64(x), nil
	case json.Number:
		return x.Int64()
	case float64:
		return int64(x), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(x), 10, 64)
	default:
		return 0, fmt.Errorf("must be integer")
	}
}
