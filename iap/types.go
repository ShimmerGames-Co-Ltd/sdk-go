package iap

// 外部订单状态，对齐 iap.v1.OrderExternalState。
const (
	StatePending       int32 = 1
	StatePayFailed     int32 = 2
	StatePaySuccess    int32 = 3
	StateVerifyFailed  int32 = 4
	StateVerifyPending int32 = 5
	StateVerifySuccess int32 = 6
	StateDone          int32 = 7
)

type AdjustInfo struct {
	SourceADID  string `json:"source_adid,omitempty"`
	ADID        string `json:"adid,omitempty"`
	Environment string `json:"environment,omitempty"`
}

type VerifyOrderRequest struct {
	OrderID string `json:"order_id"`
	Extras  string `json:"extras,omitempty"`
}

type VerifyOrderReply struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	State     int32  `json:"state"`
	Extras    string `json:"extras"`
}

// VerifySuccess 是否可作为发奖依据（对齐旧 sdk-go Success 派生）。
func VerifySuccess(state int32) bool {
	return state == StateVerifySuccess
}

type RecordPurchaseRequest struct {
	PackageName  string      `json:"package_name"`
	PayChannel   string      `json:"pay_channel"`
	ProductID    string      `json:"product_id"`
	PurchaseToken string     `json:"purchase_token"`
	Extras       string      `json:"extras,omitempty"`
	Additional   string      `json:"additional,omitempty"`
	AdjustInfo   *AdjustInfo `json:"adjust_info,omitempty"`
	UserID       int64       `json:"user_id,omitempty"`
	AuthOpenID   string      `json:"auth_openid,omitempty"`
	AuthChannel  string      `json:"auth_channel,omitempty"`
	RoleID       string      `json:"role_id,omitempty"`
}
