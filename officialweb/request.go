package officialweb

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// LookupRoleRequest POST /payment/lookup_role。
type LookupRoleRequest struct {
	CpRoleID   string `json:"cp_role_id"`
	CpServerID string `json:"cp_server_id"`
}

// ListProductsRequest POST /payment/get_goods_list。
type ListProductsRequest struct {
	PlatformRoleID string `json:"platform_role_id"`
}

// PreCheckRequest POST /payment/pre_check。
type PreCheckRequest struct {
	PlatformRoleID string  `json:"platform_role_id"`
	ProductID      string  `json:"product_id"`
	Param          []int64 `json:"param"` // 可选
}

// AddOrderRequest POST /payment/add_order。
type AddOrderRequest struct {
	PlatformRoleID  string  `json:"platform_role_id"`
	ProductID       string  `json:"product_id"`
	Param           []int64 `json:"param"`             // 可选
	ProductNameI18n string  `json:"product_name_i18n"` // 可选；线格式 JSON 字符串
}

// BuyPaymentRequest POST /payment/buypayment。
type BuyPaymentRequest struct {
	PlatformRoleID  string `json:"platform_role_id"`
	PlatformOrderID string `json:"platform_order_id"` // 可为空串
	CpOrderID       string `json:"cp_order_id"`
	CpExt           string `json:"cp_ext"` // 可为空串
}

// unmarshalBody：JSON → struct（sign 等未知字段忽略）。不做必填校验（V0）。
func unmarshalBody(raw []byte, dst any) error {
	if len(bytes.TrimSpace(raw)) == 0 {
		raw = []byte("{}")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("officialweb: invalid body: %w", err)
	}
	return nil
}
