package iap

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
	// State verifying: 创建完成; purchased: 已支付; consumed: 已消费; done: 订单完结; canceled: 订单取消; invalid: 订单无效
	// Extras 平台透传信息(可选)
	VerifyPaymentResponse struct {
		OrderId   string `json:"order_id"`
		ProductId string `json:"product_id"`
		State     string `json:"state"`
		Extras    string `json:"extras"`
	}
)
