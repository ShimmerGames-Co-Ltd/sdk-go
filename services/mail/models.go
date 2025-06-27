package mail

// -------- RPC Request / Reply 对 ------------------------------------------------

// LoginRequest / LoginReply 玩家登录
type (
	// LoginRequest 玩家登录请求参数
	LoginRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		PlayerID string `json:"playerID"`
		// 玩家所在区服ID
		ServerID string `json:"serverID"`
		// 玩家目前的语言
		Lang string `json:"lang"`
		// 玩家注册时间
		RegisterAt int64 `json:"registerAt"`
	}

	// LoginReply 玩家登录返回参数
	LoginReply struct{}
)

// SendPlayerMailRequest / SendPlayerMailReply 发送个人邮件
type (
	// SendPlayerMailRequest 发送个人邮件请求参数
	SendPlayerMailRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 发件人
		From string `json:"from"`
		// 收件人
		To string `json:"to"`
		// 邮件编号（编号相同表示同一封邮件，只处理一次）
		Serial string `json:"serial"`
		// 标题
		Title string `json:"title"`
		// 内容
		Content string `json:"content"`
		// 模板ID
		TemplateID string `json:"templateID"`
		// 模板参数
		TemplateArgs string `json:"templateArgs"`
		// 附件
		Attachment string `json:"attachment"`
		// 过期时间
		Ttl int64 `json:"ttl"`
	}

	// SendPlayerMailReply 发送个人邮件的返回值
	SendPlayerMailReply struct {
		// 邮件唯一标识
		MailRefID string `json:"mailRefID"`
	}
)

// SendServerMailRequest / SendServerMailReply 发送全服邮件
type (
	// SendServerMailRequest 发送全服邮件请求参数
	SendServerMailRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 发件人
		From string `json:"from"`
		// 收件区服IDs
		ServerIDs []string `json:"serverIDs"`
		// 邮件编号（编号相同表示同一封邮件，只处理一次）
		Serial string `json:"serial"`
		// 标题
		Title []I18NText `json:"title"`
		// 内容
		Content []I18NText `json:"content"`
		// 模板ID
		TemplateID string `json:"templateID"`
		// 模板参数
		TemplateArgs string `json:"templateArgs"`
		// 附件
		Attachment string `json:"attachment"`
		// 过期时间
		Ttl int64 `json:"ttl"`
		// 邮件接收者规则（参考：https://gitlab.shimmer.lab/platform/gparser/-/blob/main/README.md）
		MatchRule string `json:"matchRule"`
	}

	// SendServerMailReply 发送全服邮件的返回值
	SendServerMailReply struct {
		// 邮件唯一标识
		Refs []ServerMailRef `json:"refs"`
	}
)

// SendGroupMailRequest / SendGroupMailReply 发送群邮件
type (
	// SendGroupMailRequest 发送群邮件请求参数
	SendGroupMailRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 发送方
		From string `json:"from"`
		// 接收方
		To []string `json:"to"`
		// 邮件编号（编号相同表示同一封邮件，只处理一次）
		Serial string `json:"serial"`
		// 标题
		Title []I18NText `json:"title"`
		// 内容
		Content []I18NText `json:"content"`
		// 模板ID
		TemplateID string `json:"templateID"`
		// 模板参数
		TemplateArgs string `json:"templateArgs"`
		// 附件
		Attachment string `json:"attachment"`
		// 过期时间
		Ttl int64 `json:"ttl"`
	}

	// SendGroupMailReply 发送群邮件返回值
	SendGroupMailReply struct {
		// 邮件唯一标识
		Refs []GroupMailRef `json:"refs"`
	}
)

// SyncMailRequest / SyncMailReply 同步新邮件数量
type (
	// SyncMailRequest 玩家同步自己新邮件的数量请求参数
	SyncMailRequest struct {
		AppID    string `json:"appID"`
		PlayerID string `json:"playerID"`
		// 用于匹配全服邮件的条件
		IntMatcher map[string]int64 `json:"intMatcher"`
		// 用于匹配全服邮件的条件
		StrMatcher map[string]string `json:"strMatcher"`
	}

	// SyncMailReply 玩家同步自己新邮件的数量返回值
	SyncMailReply struct {
		NewMailCount int32 `json:"newMailCount"`
	}
)

// ListMailRequest / ListMailReply 分页拉取邮件
type (
	// ListMailRequest 玩家分页拉取邮件请求参数
	ListMailRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		PlayerID string `json:"playerID"`
		// 上次拉取的最后一份邮件的ID
		LatestMailID int64 `json:"latestMailID"`
		// 拉取邮件数量
		Limit int32 `json:"limit"`
	}

	// ListMailReply 玩家分页拉取邮件返回值
	ListMailReply struct {
		// 拉取到的邮件头
		Headers []MailHeader `json:"headers"`
		// 总邮件数
		Total int32 `json:"total"`
	}
)

// GetMailRequest / GetMailReply 获取邮件详情
type (
	// GetMailRequest 获取邮件详情请求参数
	GetMailRequest struct {
		AppID    string  `json:"appID"`
		PlayerID string  `json:"playerID"`
		MailIDs  []int64 `json:"mailIDs"`
	}

	// GetMailReply 获取邮件详情返回值
	GetMailReply struct {
		Mails []Mail `json:"mails"`
	}
)

// ReadMailRequest / ReadMailReply 标记已读
type (
	// ReadMailRequest 邮件标记为已读请求参数
	ReadMailRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		PlayerID string `json:"playerID"`
		// 是否全部已读
		All bool `json:"all"`
		// 标记为已读的邮件ID
		MailIDs []int64 `json:"mailIDs"`
	}

	// ReadMailReply 邮件标记为已读返回值
	ReadMailReply struct {
		// 所有已读ID
		ReadMailIds []int64 `json:"readMailIds"`
	}
)

// ClaimMailRequest / ClaimMailReply 领取附件
type (
	// ClaimMailRequest 领取邮件附件请求参数
	ClaimMailRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		PlayerID string `json:"playerID"`
		// 是否领取所有
		All bool `json:"all"`
		// 需要领取的邮件ID
		MailIDs []int64 `json:"mailIDs"`
	}

	// ClaimMailReply 领取邮件附件返回值
	ClaimMailReply struct {
		Attachments []Attachment `json:"attachments"`
	}
)

// RemoveMailRequest / RemoveMailReply 删除邮件
type (
	// RemoveMailRequest 删除邮件请求参数
	RemoveMailRequest struct {
		// 游戏ID
		AppID string `json:"appID"`
		// 玩家ID
		PlayerID string `json:"playerID"`
		// 是否全部删除
		All bool `json:"all"`
		// 需要删除的邮件ID
		MailIDs []int64 `json:"mailIDs"`
	}

	// RemoveMailReply 删除邮件返回值
	RemoveMailReply struct {
		// 被成功删除的邮件ID
		RemovedMailIDs []int64 `json:"removedMailIDs"`
	}
)

type Type int32
type Status int32

// Attachment 附件信息
type Attachment struct {
	// 邮件ID
	MailID int64 `json:"mailID"`
	// 邮件模板ID
	TemplateID string `json:"templateID"`
	// 邮件附件
	Attachment string `json:"attachment"`
}

// I18NText 国际化文本
type I18NText struct {
	// 语言
	Lang string `json:"lang"`
	// 文本
	Text string `json:"text"`
}

// 引用
type (
	ServerMailRef struct {
		ServerID  string `json:"serverID"`
		MailRefID string `json:"mailRefID"`
	}

	GroupMailRef struct {
		To        string `json:"to"`
		MailRefID string `json:"mailRefID"`
	}
)

// MailHeader 邮件头
type MailHeader struct {
	// 邮件ID
	Id int64 `json:"id"`
	// 邮件类型
	Type Type `json:"type"`
	// 发送者
	From string `json:"from"`
	// 邮件标题
	Title string `json:"title"`
	// 邮件模板ID
	TemplateId string `json:"templateId"`
	// 邮件是否包含附件
	HaveAttachment bool `json:"haveAttachment"`
	// 邮件状态
	Status Status `json:"status"`
	// 邮件创建时间
	CreateAt int64 `json:"createAt"`
	// 邮件过期时间
	ExpireAt int64 `json:"expireAt"`
	// 邮件过期倒计时
	ExpireCD int64 `json:"expireCD"`
}

// Mail 邮件详情
type Mail struct {
	// 邮件ID
	Id int64 `json:"id"`
	// 邮件类型
	Type Type `json:"type"`
	// 发送者
	From string `json:"from"`
	// 邮件标题
	Title string `json:"title"`
	// 邮件内容
	Content string `json:"content"`
	// 邮件模板ID
	TemplateId string `json:"templateId"`
	// 邮件模板参数
	TemplateArgs string `json:"templateArgs"`
	// 邮件附件
	Attachment string `json:"attachment"`
	// 邮件创建时间
	CreateAt int64 `json:"createAt"`
	// 邮件过期时间
	ExpireAT int64 `json:"expireAT"`
	// 邮件过期倒计时
	ExpireCD int64 `json:"expireCD"`
}
