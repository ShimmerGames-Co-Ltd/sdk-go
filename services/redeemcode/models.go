package redeemcode

type (
	RedeemCodeRequest struct {
		Code   string `json:"code"`
		RoleId int    `json:"role_id"`
	}
	RedeemCodeResponse struct {
		Result int               `json:"result"` //CodeResul 见下面状态码
		Code   *Code             `json:"code"`
		Record *RedeemCodeRecord `json:"record"`
	}
)

const (
	CodeResultOK = 0 // 领取成功

	CodeResultNotExist        = 1 // 兑换码不存在
	CodeResultReceivedBySelf  = 2 // 已领取
	CodeResultReceivedByOther = 3 // 已被他人领取
	CodeResultExpired         = 4 // 已过期
	CodeResultMaxTimes        = 5 // 兑换码已达到最大领取次数
	CodeResultFailed          = 6 // 兑换码领取失败
)

type CodeType int

const (
	CodeTypeNone       CodeType = 0
	CodeTypeDisposable CodeType = 1
	CodeTypeReusable   CodeType = 2
	CodeTypeGlobal     CodeType = 3
)

type CodeStatus int

const (
	CodeStatusCreated CodeStatus = 0
	CodeStatusActive  CodeStatus = 1
	CodeStatusUseOut  CodeStatus = 2
	CodeStatusExpire  CodeStatus = 3
)

type Code struct {
	ID          int64      `json:"ID,omitempty"`
	Code        string     `json:"code,omitempty"`
	Channel     string     `json:"channel,omitempty"`
	Reward      string     `json:"reward,omitempty"`
	GroupId     int64      `json:"groupId,omitempty"`
	AppId       string     `json:"appId,omitempty"`
	Type        CodeType   `json:"type,omitempty"`
	CreateTime  int64      `json:"createTime,omitempty"`
	ActiveTime  int64      `json:"activeTime,omitempty"`
	EffectTime  int64      `json:"effectTime,omitempty"`
	ExpireTime  int64      `json:"expireTime,omitempty"`
	MaxUseCount int64      `json:"maxUseCount,omitempty"`
	UseCount    int64      `json:"useCount,omitempty"`
	Status      CodeStatus `json:"status,omitempty"`
}

type RedeemCodeRecord struct {
	UserId     string `json:"userId,omitempty"`
	CodeId     int64  `json:"codeId,omitempty"`
	AppId      string `json:"appId,omitempty"`
	Code       string `json:"code,omitempty"`
	Channel    string `json:"channel,omitempty"`
	Rewarded   bool   `json:"rewarded,omitempty"`
	GroupId    int64  `json:"groupId,omitempty"`
	RedeemTime int64  `json:"redeemTime,omitempty"`
}
