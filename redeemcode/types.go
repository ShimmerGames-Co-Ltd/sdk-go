package redeemcode

const (
	ResultOK             int32 = 0
	ResultNotExist       int32 = 1
	ResultReceivedBySelf int32 = 2
	ResultReceivedByOther int32 = 3
	ResultExpired        int32 = 4
	ResultMaxTimes       int32 = 5
	ResultFailed         int32 = 6
)

type RedeemedCode struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Channel     string `json:"channel"`
	Reward      string `json:"reward"`
	GroupID     int64  `json:"group_id"`
	AppID       string `json:"app_id"`
	Type        int32  `json:"type"`
	CreateTime  int64  `json:"create_time"`
	ActiveTime  int64  `json:"active_time"`
	EffectTime  int64  `json:"effect_time"`
	ExpireTime  int64  `json:"expire_time"`
	MaxUseCount int64  `json:"max_use_count"`
	UseCount    int64  `json:"use_count"`
	Status      int32  `json:"status"`
}

type RedeemedRecord struct {
	UserID     string `json:"user_id"`
	CodeID     int64  `json:"code_id"`
	AppID      string `json:"app_id"`
	Code       string `json:"code"`
	Channel    string `json:"channel"`
	Rewarded   bool   `json:"rewarded"`
	GroupID    int64  `json:"group_id"`
	RedeemTime int64  `json:"redeem_time"`
}

type RedeemRequest struct {
	Code      string `json:"code"`
	UserID    int64  `json:"user_id,omitempty"`
	UserIDStr string `json:"user_id_str,omitempty"`
	RoleID    int64  `json:"role_id,omitempty"`
}

type RedeemReply struct {
	Result int32           `json:"result"`
	Code   *RedeemedCode   `json:"code"`
	Record *RedeemedRecord `json:"record"`
}

type CheckRequest = RedeemRequest
type CheckReply = RedeemReply
