package mail

// I18nBlock 是一组完整的语言内容；副标题可空，标题和正文按块回退。
type I18nBlock struct {
	// Lang 使用 App 配置的语言标识（如 English），有效文案须属于已启用语言，不能重复。
	Lang string `json:"lang"`
	// Title 与 Content 须成对填写；必须包含 App 默认语言的完整块，默认语言标题最多 256 bytes。
	Title string `json:"title"`
	// Content 是正文原文，ContentFormat 不会自动转换、转义或改写正文。
	Content string `json:"content"`
	// SubTitle 可空，非空最多 80 个 Unicode 字符；随完整语言块一起回退。
	SubTitle string `json:"sub_title,omitempty"`
}

// MailAttachment 是发送附件；与领取结果 Attachment 是不同的协议对象。
type MailAttachment struct {
	// Type 是项目约定的道具/奖励类型（对应 Hub game_items.item_type），须 >= 0。
	Type int32 `json:"type"`
	// ID 是项目业务道具 ID（对应 game_items.item_id），必填；发送不检查道具是否登记。
	ID string `json:"id"`
	// Count 是本项奖励数量，须 > 0；附件整体 JSON 编码后最多 4096 bytes。
	Count int32 `json:"count"`
	// Param 是项目扩展字符串，可空；中台原样透传，不解析其内容。
	Param string `json:"param,omitempty"`
}

// RecipientTarget 是一个区服的群组收件人，同次请求中 ServerID 不可重复。
type RecipientTarget struct {
	// ServerID 是实际区服 ID，必填且最多 64 bytes，不使用全服通配符。
	ServerID string `json:"server_id"`
	// PlayerIDs 非空，元素须非空、同服不重复且各最多 64 bytes；身份语义同 SendPlayerV2Request.PlayerID。
	PlayerIDs []string `json:"player_ids"`
}

// 正文格式元数据，不改变正文内容；旧邮件读取到的空字符串表示格式未指定。
const (
	ContentFormatPlainText = "plain_text"
	ContentFormatRichText  = "rich_text"
)

// SendPlayerV2Request 发送个人邮。普通文案填写 I18n；两种模板模式只填模板字段。
// V2 空模式始终是 plain，不通过 TemplateID 自动推断；模板模式必须显式选择。
type SendPlayerV2Request struct {
	// AppID 可省略，以 core.Client 的签名 App 为准；非空时必须一致。
	AppID string `json:"app_id,omitempty"`
	// From 是展示的发件人，必填且最多 128 bytes。
	From string `json:"from"`
	// Serial 是稳定业务编号，必填且最多 128 bytes；按 App+区服+玩家+Serial 幂等，重试不改内容。
	Serial string `json:"serial"`
	// ContentMode 使用 ContentMode*；空为 plain，plain 禁止模板 ID 和非空模板参数。
	ContentMode string `json:"content_mode,omitempty"`
	// I18n 在 plain 模式必填，须包含默认语言完整块；两种模板模式必须留空。
	I18n []I18nBlock `json:"i18n,omitempty"`
	// ContentCategory 是游戏自定义分类，最多 64 bytes；空表示无分类，不校验登记、不继承模板分类。
	ContentCategory string `json:"content_category,omitempty"`
	// ContentFormat 使用 ContentFormat*；普通/客户端模板空为 rich_text；平台模板空用固定版本，非空须一致。
	ContentFormat string `json:"content_format,omitempty"`
	// Attachments 是此次发送的奖励；可省略，不从模板继承。
	Attachments []MailAttachment `json:"attachments,omitempty"`
	// ExpireAt 是 Unix 秒；0 为发送后 14 天，非零须晚于发送时刻；重试不更新原邮件过期时间。
	ExpireAt int64 `json:"expire_at,omitempty"`
	// TemplateID 在模板模式必填，最多 128 bytes；平台模板填后台 code，客户端模板填项目标识。
	// 已登记后台的模板 code 不能作为客户端模板 ID。
	TemplateID string `json:"template_id,omitempty"`
	// TemplateArgs 直接填写对象；平台模板 plain.* 用标量，i18n.* 用语言到标量的 map，须含默认语言。
	// 平台参数须与固定版本占位符完全一致，编码后最多 8192 bytes；客户端参数按项目约定透传。
	// 参数目录的登记/启用要求在后台维护模板时检查，发送时按固定版本的占位符集合校验。
	TemplateArgs map[string]any `json:"template_args,omitempty"`
	// ServerID 是实际区服 ID，必填且最多 64 bytes，不使用全服通配符。
	ServerID string `json:"server_id"`
	// PlayerID 是该区服邮箱的玩家标识，必填且最多 64 bytes，须与读取邮箱时使用的身份一致。
	// 默认客户端直连模式使用 Auth 角色 role_id 的十进制字符串；自定义/游戏服模式使用约定的项目标识。
	PlayerID string `json:"player_id"`
}

// SendGroupV2Request 按区服分组发送给明确玩家，不使用 Audience 或后注册玩家过滤。
// 内容字段规则同 SendPlayerV2Request；幂等范围为 App+区服+Serial，同一业务重试保留原收件人和内容。
type SendGroupV2Request struct {
	// AppID 可省略；非空须与 core.Client 签名 App 一致。
	AppID string `json:"app_id,omitempty"`
	// From 是发件人，必填且最多 128 bytes。
	From string `json:"from"`
	// Serial 是稳定业务编号，必填且最多 128 bytes；按 App+区服+Serial 幂等。
	Serial string `json:"serial"`
	// ContentMode 使用 ContentMode*；空为 plain，模板模式必须显式选择。
	ContentMode string `json:"content_mode,omitempty"`
	// I18n 在 plain 模式必填且含默认语言完整块；两种模板模式留空。
	I18n []I18nBlock `json:"i18n,omitempty"`
	// ContentCategory 是游戏自定义分类，最多 64 bytes；空表示无分类，不校验登记、不继承模板分类。
	ContentCategory string `json:"content_category,omitempty"`
	// ContentFormat 普通/客户端模板默认 rich_text；平台模板空用固定版本，非空须一致。
	ContentFormat string `json:"content_format,omitempty"`
	// Attachments 是此次发送的奖励，可省略；规则见 MailAttachment。
	Attachments []MailAttachment `json:"attachments,omitempty"`
	// ExpireAt 是 Unix 秒；0 为发送后 14 天，非零须晚于发送时刻。
	ExpireAt int64 `json:"expire_at,omitempty"`
	// TemplateID 模板模式必填；平台填后台 code，客户端填未与后台 code 冲突的项目标识，最多 128 bytes。
	TemplateID string `json:"template_id,omitempty"`
	// TemplateArgs 直接传对象；参数规则与大小边界见 SendPlayerV2Request.TemplateArgs。
	TemplateArgs map[string]any `json:"template_args,omitempty"`
	// RecipientTargets 非空，每服只出现一次；各服玩家列表非空且不重复。
	RecipientTargets []RecipientTarget `json:"recipient_targets"`
}

// SendServerV2Request 写入区服邮箱，玩家在 Sync 时按受众和注册截止条件收取。
// 内容字段规则同 SendPlayerV2Request；幂等范围为 App+区服+Serial。
type SendServerV2Request struct {
	// AppID 可省略；非空须与 core.Client 签名 App 一致。
	AppID string `json:"app_id,omitempty"`
	// From 是发件人，必填且最多 128 bytes。
	From string `json:"from"`
	// Serial 是稳定业务编号，必填且最多 128 bytes；按 App+区服+Serial 幂等。
	Serial string `json:"serial"`
	// ContentMode 使用 ContentMode*；空为 plain，模板模式必须显式选择。
	ContentMode string `json:"content_mode,omitempty"`
	// I18n 在 plain 模式必填且含默认语言完整块；两种模板模式留空。
	I18n []I18nBlock `json:"i18n,omitempty"`
	// ContentCategory 是游戏自定义分类，最多 64 bytes；空表示无分类，不校验登记、不继承模板分类。
	ContentCategory string `json:"content_category,omitempty"`
	// ContentFormat 普通/客户端模板默认 rich_text；平台模板空用固定版本，非空须一致。
	ContentFormat string `json:"content_format,omitempty"`
	// Attachments 是此次发送的奖励，可省略；规则见 MailAttachment。
	Attachments []MailAttachment `json:"attachments,omitempty"`
	// ExpireAt 是 Unix 秒；0 为发送后 14 天，非零须晚于发送时刻。
	ExpireAt int64 `json:"expire_at,omitempty"`
	// TemplateID 模板模式必填；平台填后台 code，客户端填未与后台 code 冲突的项目标识，最多 128 bytes。
	TemplateID string `json:"template_id,omitempty"`
	// TemplateArgs 直接传对象；参数规则与大小边界见 SendPlayerV2Request.TemplateArgs。
	TemplateArgs map[string]any `json:"template_args,omitempty"`
	// ServerIDs 非空、不重复，各最多 64 bytes；[]string{"*"} 表示 App 全服，不与具体区服混用。
	ServerIDs []string `json:"server_ids"`
	// Audience 为 nil 时不附加筛选；非空规则与 ServerIDs 取交集，使用接收方 Sync.PlayerAttrs 匹配。
	Audience *Audience `json:"audience,omitempty"`
	// IncludeLateRegistrants 默认 false：即使 Audience 为 nil，Sync 也须提供合法且严格早于发送时刻的
	// PlayerAttrs["register_at"]（Unix 秒）；true 允许后注册玩家，但不绕过 Audience 规则。
	IncludeLateRegistrants bool `json:"include_late_registrants,omitempty"`
}

// SendPlayerV2Reply 表示发信成功或命中既有幂等记录。
type SendPlayerV2Reply struct {
	// MailRefID 是正文引用，不是 Get/Read/Claim 的邮箱内 ID；删信后同 Serial 重试成功时可为空。
	MailRefID string `json:"mail_ref_id"`
	// Badge 是该玩家邮箱的角标快照，不表示本次新增投递数量。
	Badge MailboxBadge `json:"badge"`
}

// GroupMailV2Ref 是一个明确收件人的发送结果；同服玩家可能共享同一正文引用。
type GroupMailV2Ref struct {
	// ServerID 是该收件人所属的区服。
	ServerID string `json:"server_id"`
	// PlayerID 是该区服的邮箱玩家标识，与请求一致。
	PlayerID string `json:"player_id"`
	// MailRefID 是正文引用，不是 List 返回的邮箱内 ID。
	MailRefID string `json:"mail_ref_id"`
	// Badge 是该收件人的邮箱角标快照。
	Badge MailboxBadge `json:"badge"`
}

// SendGroupV2Reply 仅在整次请求成功时返回；失败时可能已有部分收件人收到邮件。
type SendGroupV2Reply struct {
	// Refs 按区服和玩家标识关联结果，不以正文引用判断收件人数量。
	Refs []GroupMailV2Ref `json:"refs"`
}

// ServerMailV2Ref 表示已写入一个区服范围，玩家是否进箱仍由 Sync 决定。
type ServerMailV2Ref struct {
	// ServerID 是请求的区服范围；全服请求返回 "*"，不会展开成玩家或实际区服列表。
	ServerID string `json:"server_id"`
	// MailRefID 是区服邮件正文引用，不是玩家邮箱内 ID。
	MailRefID string `json:"mail_ref_id"`
}

// SendServerV2Reply 仅在整次请求成功时返回，结果不包含玩家角标。
type SendServerV2Reply struct {
	// Refs 是各区服范围的写入结果，不代表玩家已经收取。
	Refs []ServerMailV2Ref `json:"refs"`
}

// 条件组使用 And/Or；规则使用其余操作符，且须属于该 App Catalog 字段的 allowed_ops。
const (
	AudienceOpAnd   = "and"
	AudienceOpOr    = "or"
	AudienceOpEq    = "eq"
	AudienceOpIn    = "in"
	AudienceOpRange = "range"
	AudienceOpGte   = "gte"
	AudienceOpLte   = "lte"
	AudienceOpGt    = "gt"
	AudienceOpLt    = "lt"
)

// Audience 是区服邮件的根条件组，与 ServerIDs 取交集；不筛选时将请求的 Audience 留 nil。
// 根组及嵌套组使用 and/or，整棵树至少包含一条规则；字段和值域由 App 的 Audience Catalog 决定。
type Audience struct {
	// Op 使用 AudienceOpAnd 或 AudienceOpOr。
	Op string `json:"op"`
	// Children 包含规则或嵌套条件组；根组必须包含至少一条有效规则。
	Children []AudienceNode `json:"children"`
}

// AudienceNode 表示一条规则或嵌套条件组，直接使用结构体字面量组合即可。
// 规则填写 Key/Op/Value，Children 留空；条件组填写 Op/Children，Key/Value 留空。
type AudienceNode struct {
	// Op：组用 and/or；规则用 eq/in/range/gte/lte/gt/lt，受 Catalog allowed_ops 限制。
	Op string `json:"op"`
	// Key 是 Catalog 中的字段名，与接收方 Sync.PlayerAttrs 的同名属性匹配。
	Key string `json:"key,omitempty"`
	// Value 按 Catalog 类型填写标量；in 用同类型切片，range 用 AudienceRange。
	// 数值使用 Go 数字，版本使用字符串，布尔使用 bool；0、false 和空字符串均会保留在请求中。
	Value any `json:"value,omitempty"`
	// Children 仅供嵌套条件组使用；规则留空，服务端据 Key/Op 识别规则节点。
	Children []AudienceNode `json:"children,omitempty"`
}

// AudienceRange 是区间 [Min, Max)，Min 必填，Max 可省略；同时提供时 Max 须大于 Min。
// 边界按 Catalog 类型使用非负数字或版本字符串；any 用于保留整数/小数/版本及显式零值。
type AudienceRange struct {
	// Min 是包含的下界，必填；数值 0 会保留，不能用 nil 省略。
	Min any `json:"min"`
	// Max 是不包含的上界；nil 表示不限制上界。
	Max any `json:"max,omitempty"`
}
