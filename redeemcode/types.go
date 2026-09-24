package redeemcode

import "github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"

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
	ID          core.JSONInt64  `json:"id"`
	Code        string `json:"code"`
	Channel     string `json:"channel"`
	Reward      string `json:"reward"`
	GroupID     core.JSONInt64  `json:"group_id"`
	AppID       string `json:"app_id"`
	Type        int32  `json:"type"`
	CreateTime  core.JSONInt64  `json:"create_time"`
	ActiveTime  core.JSONInt64  `json:"active_time"`
	EffectTime  core.JSONInt64  `json:"effect_time"`
	ExpireTime  core.JSONInt64  `json:"expire_time"`
	MaxUseCount core.JSONInt64  `json:"max_use_count"`
	UseCount    core.JSONInt64  `json:"use_count"`
	Status      int32  `json:"status"`
}

type RedeemedRecord struct {
	UserID     string `json:"user_id"`
	CodeID     core.JSONInt64  `json:"code_id"`
	AppID      string `json:"app_id"`
	Code       string `json:"code"`
	Channel    string `json:"channel"`
	Rewarded   bool   `json:"rewarded"`
	GroupID    core.JSONInt64  `json:"group_id"`
	RedeemTime core.JSONInt64  `json:"redeem_time"`
}

type RedeemRequest struct {
	Code      string `json:"code"`
	UserID    core.JSONInt64  `json:"user_id,omitempty"`
	UserIDStr string `json:"user_id_str,omitempty"`
	RoleID    core.JSONInt64  `json:"role_id,omitempty"`
}

type RedeemReply struct {
	Result int32           `json:"result"`
	Code   *RedeemedCode   `json:"code"`
	Record *RedeemedRecord `json:"record"`
}

type CheckRequest = RedeemRequest
type CheckReply = RedeemReply
