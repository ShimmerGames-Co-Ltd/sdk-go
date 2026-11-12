package mail

import "github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"

// 请求与响应 data 均为 snake_case（sdk-go v3）；服务端入站仍可接受旧 camelCase。

type MailType int32

const (
	MailTypeNone   MailType = 0
	MailTypePlayer MailType = 1
	MailTypeServer MailType = 2
	MailTypeGroup  MailType = 3
)

// 内容模式：普通文案、客户端本地模板、平台固定版本模板。
// V2 空值等价 plain；模板模式必须显式选择，不能通过 TemplateID 推断。
const (
	ContentModeUnspecified      = ""
	ContentModePlain            = "plain"
	ContentModeClientTemplate   = "client_template"
	ContentModePlatformTemplate = "platform_template"
)

type MailboxBadge struct {
	// PendingNewMailCount 是邮箱中尚未被 List 确认看过的邮件总数，不等于本次新增数。
	PendingNewMailCount int32 `json:"pending_new_mail_count"`
	// UnreadMailCount 是当前未读邮件总数。
	UnreadMailCount int32 `json:"unread_mail_count"`
}

type MailStatus struct {
	// ReadAt 是已读 Unix 秒；0 表示未读。
	ReadAt core.JSONInt64 `json:"read_at"`
	// ClaimAt 是已领取 Unix 秒；0 表示未领取。
	ClaimAt core.JSONInt64 `json:"claim_at"`
}

// MailHeader 是 List 及个人/群组通知中的邮箱摘要；ID 属于一个玩家的邮箱。
type MailHeader struct {
	// ID 是邮箱内数值 ID，用于 Get/Read/Claim/Remove；不是发送返回的 MailRefID。
	ID core.JSONInt64 `json:"id"`
	// Type 使用 MailType*，表示个人、区服或群组来源。
	Type MailType `json:"type"`
	// From 是发送时的发件人。
	From string `json:"from"`
	// Title 是按邮箱当前语言及默认语言回退后得到的标题。
	Title string `json:"title"`
	// TemplateID 仅客户端模板透传；平台模板对接收方表现为普通文案，此字段为空。
	TemplateID string `json:"template_id"`
	// HaveAttachment 表示邮件声明有附件，是否已领取见 Status.ClaimAt。
	HaveAttachment bool `json:"have_attachment"`
	// Status 是该玩家的阅读/领取状态。
	Status MailStatus `json:"status"`
	// CreateAt 是创建 Unix 秒。
	CreateAt core.JSONInt64 `json:"create_at"`
	// ExpireAt 是过期 Unix 秒。
	ExpireAt core.JSONInt64 `json:"expire_at"`
	// ExpireCD 是服务端生成的剩余有效秒数。
	ExpireCD core.JSONInt64 `json:"expire_cd"`
	// ContentCategory 是发送时的分类标识；未设置时为空。
	ContentCategory string `json:"content_category"`
	// ContentFormat 使用 ContentFormat*；旧无格式邮件为空，调用方按原有展示策略处理。
	ContentFormat string `json:"content_format"`
	// SubTitle 与 Title/Content 使用同一个完整语言块回退，可为空。
	SubTitle string `json:"sub_title"`
}

// MailDetail 是 Get 返回的详情；附件和客户端模板参数仍是 JSON 字符串，区别于 V2 发送对象。
type MailDetail struct {
	// ID 是邮箱内数值 ID，与 List/通知中的 MailHeader.ID 对应。
	ID core.JSONInt64 `json:"id"`
	// Type 使用 MailType*。
	Type MailType `json:"type"`
	// From 是发送时的发件人。
	From string `json:"from"`
	// Title 是按邮箱当前语言回退后的标题。
	Title string `json:"title"`
	// Content 是正文原文或平台模板渲染结果，按 ContentFormat 展示。
	Content string `json:"content"`
	// TemplateID 仅客户端模板透传；普通文案和平台模板为空。
	TemplateID string `json:"template_id"`
	// TemplateArgs 透传客户端模板参数字符串；平台模板返回渲染后文案，此字段为空。
	TemplateArgs string `json:"template_args"`
	// Attachment 是 JSON 数组字符串，可解码为 []MailAttachment；空值表示无附件。
	Attachment string `json:"attachment"`
	// CreateAt 是创建 Unix 秒。
	CreateAt core.JSONInt64 `json:"create_at"`
	// ExpireAt 是过期 Unix 秒。
	ExpireAt core.JSONInt64 `json:"expire_at"`
	// ExpireCD 是服务端生成的剩余有效秒数。
	ExpireCD core.JSONInt64 `json:"expire_cd"`
	// ContentCategory 是发送时的分类标识，未设置时为空。
	ContentCategory string `json:"content_category"`
	// ContentFormat 使用 ContentFormat*；旧无格式邮件为空，调用方按原有展示策略处理。
	ContentFormat string `json:"content_format"`
	// SubTitle 可为空，与标题、正文按完整语言块回退。
	SubTitle string `json:"sub_title"`
}

// Attachment 是 Claim 按邮件返回的领取结果包装，不是一项奖励。
type Attachment struct {
	// MailID 是本次领取的邮箱内 ID。
	MailID core.JSONInt64 `json:"mail_id"`
	// TemplateID 仅客户端模板透传，平台模板为空。
	TemplateID string `json:"template_id"`
	// Attachment 是本封邮件的奖励 JSON 数组字符串，可解码为 []MailAttachment。
	Attachment string `json:"attachment"`
}

type I18nText struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

type LoginRequest struct {
	AppID    string `json:"app_id,omitempty"`
	PlayerID string `json:"player_id"`
	ServerID string `json:"server_id"`
	Lang     string `json:"lang"`
}

type LoginReply struct {
	Badge MailboxBadge `json:"badge"`
}

type SyncRequest struct {
	// AppID 可省略；非空须与 core.Client 签名 App 一致。
	AppID string `json:"app_id,omitempty"`
	// PlayerID 必须与发送及 Login 使用的该区服邮箱玩家标识一致。
	PlayerID string `json:"player_id"`
	// ServerID 是读取的实际区服 ID，不传 "*"。
	ServerID string `json:"server_id"`
	// PlayerAttrs 是本次玩家真实属性，键和类型与 App Audience Catalog 对齐，整体可选。
	// 收取默认 IncludeLateRegistrants=false 的 V2 区服邮时须传 register_at（Unix 秒整数）。
	// 缺省 server_id 由顶层 ServerID 补齐；缺省渠道/平台可由 package_name 对应的包配置补齐。
	// 其余规则所需属性由游戏服提供；游戏服模式不从请求头补包名、版本或国家。
	PlayerAttrs map[string]any `json:"player_attrs,omitempty"`
}

type SyncReply struct {
	// NewMailCount 仅为本次 Sync 新合并进邮箱的邮件数。
	NewMailCount int32 `json:"new_mail_count"`
	// Badge 是当前总角标，Sync 不推进 List 的已看游标。
	Badge MailboxBadge `json:"badge"`
}

type ListRequest struct {
	AppID        string         `json:"app_id,omitempty"`
	PlayerID     string         `json:"player_id"`
	ServerID     string         `json:"server_id"`
	LatestMailID core.JSONInt64 `json:"latest_mail_id,omitempty"`
	Limit        int32          `json:"limit,omitempty"`
}

type ListReply struct {
	Headers []MailHeader `json:"headers"`
	Total   int32        `json:"total"`
	Badge   MailboxBadge `json:"badge"`
}

type GetRequest struct {
	AppID    string  `json:"app_id,omitempty"`
	PlayerID string  `json:"player_id"`
	ServerID string  `json:"server_id"`
	MailIDs  []int64 `json:"mail_ids"`
}

type GetReply struct {
	Mails []MailDetail `json:"mails"`
}

type ReadRequest struct {
	AppID    string  `json:"app_id,omitempty"`
	PlayerID string  `json:"player_id"`
	ServerID string  `json:"server_id"`
	All      bool    `json:"all,omitempty"`
	MailIDs  []int64 `json:"mail_ids,omitempty"`
}

type ReadReply struct {
	ReadMailIDs []int64      `json:"read_mail_ids"`
	Badge       MailboxBadge `json:"badge"`
}

type ClaimRequest struct {
	AppID    string  `json:"app_id,omitempty"`
	PlayerID string  `json:"player_id"`
	ServerID string  `json:"server_id"`
	All      bool    `json:"all,omitempty"`
	MailIDs  []int64 `json:"mail_ids,omitempty"`
}

type ClaimReply struct {
	Attachments []Attachment `json:"attachments"`
	Badge       MailboxBadge `json:"badge"`
}

type RemoveRequest struct {
	AppID    string  `json:"app_id,omitempty"`
	PlayerID string  `json:"player_id"`
	ServerID string  `json:"server_id"`
	All      bool    `json:"all,omitempty"`
	MailIDs  []int64 `json:"mail_ids,omitempty"`
}

type RemoveReply struct {
	RemovedMailIDs []int64      `json:"removed_mail_ids"`
	Badge          MailboxBadge `json:"badge"`
}

// SendPlayerRequest 个人邮。platform_template 时 TemplateArgs 为 JSON object 字符串：plain.* 标量，i18n.* 为 lang→标量 且须含 App 默认语言。
type SendPlayerRequest struct {
	AppID        string         `json:"app_id,omitempty"`
	From         string         `json:"from"`
	ServerID     string         `json:"server_id"`
	To           string         `json:"to"`
	Serial       string         `json:"serial"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	TemplateID   string         `json:"template_id,omitempty"`
	TemplateArgs string         `json:"template_args,omitempty"` // 见 SendPlayerRequest 类型注释
	Attachment   string         `json:"attachment,omitempty"`
	TTL          core.JSONInt64 `json:"ttl,omitempty"`
	ContentMode  string         `json:"content_mode,omitempty"`
}

type SendPlayerReply struct {
	MailRefID string       `json:"mail_ref_id"`
	Badge     MailboxBadge `json:"badge"`
}

// SendGroupRequest 群组邮。platform_template 时 TemplateArgs 规则同 SendPlayerRequest。
type SendGroupRequest struct {
	AppID        string         `json:"app_id,omitempty"`
	From         string         `json:"from"`
	ServerID     string         `json:"server_id"`
	To           []string       `json:"to"`
	Serial       string         `json:"serial"`
	Title        []I18nText     `json:"title"`
	Content      []I18nText     `json:"content"`
	TemplateID   string         `json:"template_id,omitempty"`
	TemplateArgs string         `json:"template_args,omitempty"`
	Attachment   string         `json:"attachment,omitempty"`
	TTL          core.JSONInt64 `json:"ttl,omitempty"`
	ContentMode  string         `json:"content_mode,omitempty"`
}

type GroupMailRef struct {
	To        string       `json:"to"`
	MailRefID string       `json:"mail_ref_id"`
	Badge     MailboxBadge `json:"badge"`
}

type SendGroupReply struct {
	Refs []GroupMailRef `json:"refs"`
}

// SendServerRequest 全服/条件邮。platform_template 时 TemplateArgs 规则同 SendPlayerRequest。
type SendServerRequest struct {
	AppID        string         `json:"app_id,omitempty"`
	From         string         `json:"from"`
	ServerIDs    []string       `json:"server_ids"`
	Serial       string         `json:"serial"`
	Title        []I18nText     `json:"title"`
	Content      []I18nText     `json:"content"`
	TemplateID   string         `json:"template_id,omitempty"`
	TemplateArgs string         `json:"template_args,omitempty"`
	Attachment   string         `json:"attachment,omitempty"`
	TTL          core.JSONInt64 `json:"ttl,omitempty"`
	MatchRule    string         `json:"match_rule,omitempty"`
	ContentMode  string         `json:"content_mode,omitempty"`
}

type ServerMailRef struct {
	ServerID  string `json:"server_id"`
	MailRefID string `json:"mail_ref_id"`
}

type SendServerReply struct {
	Refs []ServerMailRef `json:"refs"`
}

// NATS 新邮推送（snake_case；与 ListMail 的 MailHeader 字段一致）。

// ServerMailNotifyHeader 全服/条件邮 NATS 通知头。
type ServerMailNotifyHeader struct {
	ID        core.JSONInt64 `json:"id"`
	MailRefID string         `json:"mail_ref_id"`
	ExpireAt  core.JSONInt64 `json:"expire_at"`
	ExpireCD  core.JSONInt64 `json:"expire_cd"`
}

// NewPlayerMailNotify 个人/群组邮 NATS 载荷。
type NewPlayerMailNotify struct {
	ServerID string       `json:"server_id"`
	PlayerID string       `json:"player_id"`
	Headers  []MailHeader `json:"headers"`
}

// NewServerMailNotify 全服/条件邮 NATS 载荷。
type NewServerMailNotify struct {
	ServerID string                   `json:"server_id"`
	Headers  []ServerMailNotifyHeader `json:"headers"`
}

// NewMailNotify NATS 推送根类型（player_mail 与 server_mail 二选一非空）。
type NewMailNotify struct {
	PlayerMail *NewPlayerMailNotify `json:"player_mail,omitempty"`
	ServerMail *NewServerMailNotify `json:"server_mail,omitempty"`
}
