package iap

type OrderExternalState int32

const (
	OrderExternalStatePending       OrderExternalState = 1
	OrderExternalStatePayFailed     OrderExternalState = 2
	OrderExternalStatePaySuccess    OrderExternalState = 3
	OrderExternalStateVerifyFailed  OrderExternalState = 4
	OrderExternalStateVerifyPending OrderExternalState = 5
	OrderExternalStateVerifySuccess OrderExternalState = 6
	OrderExternalStateDone          OrderExternalState = 7
)

type (

	// VerifyPaymentRequest
	// OrderId 平台订单ID
	// Extras 平台透传信息(可选)
	VerifyPaymentRequest struct {
		OrderId string `json:"order_id"`
		Extras  string `json:"extras"`
	}

	// VerifyPaymentResponse
	// OrderId 平台订单ID
	// ProductId 平台商品ID
	// State 订单状态
	// Success 是否验证成功 (发奖依据)
	// Extras 平台透传信息(可选)
	VerifyPaymentResponse struct {
		OrderId   string             `json:"order_id"`
		ProductId string             `json:"product_id"`
		State     OrderExternalState `json:"state"`
		Success   bool               `json:"-"`
		Extras    string             `json:"extras"`
	}
)
