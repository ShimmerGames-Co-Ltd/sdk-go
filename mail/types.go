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

// 内容模式：见主仓 lib/consts.MailContentMode*（如 platform_template / client_template）。
const (
	ContentModeUnspecified      = ""
	ContentModePlain            = "plain"
	ContentModeClientTemplate   = "client_template"
	ContentModePlatformTemplate = "platform_template"
)

type MailboxBadge struct {
	PendingNewMailCount int32 `json:"pending_new_mail_count"`
	UnreadMailCount     int32 `json:"unread_mail_count"`
}

type MailStatus struct {
	ReadAt  core.JSONInt64 `json:"read_at"`
	ClaimAt core.JSONInt64 `json:"claim_at"`
}

type MailHeader struct {
	ID              core.JSONInt64 `json:"id"`
	Type            MailType       `json:"type"`
	From            string         `json:"from"`
	Title           string         `json:"title"`
	TemplateID      string         `json:"template_id"`
	HaveAttachment  bool           `json:"have_attachment"`
	Status          MailStatus     `json:"status"`
	CreateAt        core.JSONInt64 `json:"create_at"`
	ExpireAt        core.JSONInt64 `json:"expire_at"`
	ExpireCD        core.JSONInt64 `json:"expire_cd"`
	ContentCategory string         `json:"content_category"`
	SubTitle        string         `json:"sub_title"`
}

type MailDetail struct {
	ID              core.JSONInt64 `json:"id"`
	Type            MailType       `json:"type"`
	From            string         `json:"from"`
	Title           string         `json:"title"`
	Content         string         `json:"content"`
	TemplateID      string         `json:"template_id"`
	TemplateArgs    string         `json:"template_args"`
	Attachment      string         `json:"attachment"`
	CreateAt        core.JSONInt64 `json:"create_at"`
	ExpireAt        core.JSONInt64 `json:"expire_at"`
	ExpireCD        core.JSONInt64 `json:"expire_cd"`
	ContentCategory string         `json:"content_category"`
	SubTitle        string         `json:"sub_title"`
}

type Attachment struct {
	MailID     core.JSONInt64 `json:"mail_id"`
	TemplateID string         `json:"template_id"`
	Attachment string         `json:"attachment"`
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
	AppID       string         `json:"app_id,omitempty"`
	PlayerID    string         `json:"player_id"`
	ServerID    string         `json:"server_id"`
	PlayerAttrs map[string]any `json:"player_attrs,omitempty"`
}

type SyncReply struct {
	NewMailCount int32        `json:"new_mail_count"`
	Badge        MailboxBadge `json:"badge"`
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

type SendPlayerRequest struct {
	AppID        string         `json:"app_id,omitempty"`
	From         string         `json:"from"`
	ServerID     string         `json:"server_id"`
	To           string         `json:"to"`
	Serial       string         `json:"serial"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	TemplateID   string         `json:"template_id,omitempty"`
	TemplateArgs string         `json:"template_args,omitempty"`
	Attachment   string         `json:"attachment,omitempty"`
	TTL          core.JSONInt64 `json:"ttl,omitempty"`
	ContentMode  string         `json:"content_mode,omitempty"`
}

type SendPlayerReply struct {
	MailRefID string       `json:"mail_ref_id"`
	Badge     MailboxBadge `json:"badge"`
}

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
