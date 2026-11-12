# 邮件 V1 到 V2 迁移示例

[字段与默认行为对照](mail-v2.md#7-邮件发送-v1-到-v2-迁移)。项目用到哪些能力，就迁移哪些字段。

本例将一封带附件和注册时间筛选的区服邮件从 V1 改为 V2，保持**有效期1天、允许后注册玩家、纯文本展示**。假定 App 默认语言为 English，Catalog 的 register_at 支持 range；活动开始时间为 `1700000000`（Unix 秒）。区服、道具和条件值按项目替换。

`ctx` 为业务请求上下文，`cli` 为 `mail.New(c)` 创建的 Client，初始化见[接入指南](integration.md#2-创建-client)。下面两段分别是迁移前后的业务代码，替换使用即可。

## 1. 迁移前：V1

TTL 填有效期的秒数：`86400` 表示1天。附件和筛选条件使用 JSON 字符串。

```go
_, err := cli.SendServer(ctx, mail.SendServerRequest{
	From:        "system",
	Serial:      "activity-1001-reward",
	ServerIDs:   []string{"1", "2"},
	ContentMode: mail.ContentModePlain,
	Title: []mail.I18nText{
		{
			Lang: "English",
			Text: "Reward",
		},
	},
	Content: []mail.I18nText{
		{
			Lang: "English",
			Text: "Your reward is ready.",
		},
	},
	Attachment: `[{"type":0,"id":"gold","count":100,"param":""}]`,
	MatchRule:  `{"op":"and","children":[{"key":"register_at","op":"range","value":{"min":1700000000}}]}`,
	TTL:        core.JSONInt64(86400),
})
return err
```

## 2. 迁移后：V2

ExpireAt 填**过期时刻的 Unix 秒**。例如今天上午10点发送、有效期1天，就填明天上午10点的时间戳。下面在构造请求时计算一次，重试复用同一份 `req`。

```go
req := mail.SendServerV2Request{
	From:          "system",
	Serial:        "activity-1001-reward",
	ServerIDs:     []string{"1", "2"},
	ContentMode:   mail.ContentModePlain,
	ContentFormat: mail.ContentFormatPlainText,
	I18n: []mail.I18nBlock{
		{
			Lang:    "English",
			Title:   "Reward",
			Content: "Your reward is ready.",
		},
	},
	Attachments: []mail.MailAttachment{
		{
			Type:  0,
			ID:    "gold",
			Count: 100,
		},
	},
	Audience: &mail.Audience{
		Op: mail.AudienceOpAnd,
		Children: []mail.AudienceNode{
			{
				Key: "register_at",
				Op:  mail.AudienceOpRange,
				Value: mail.AudienceRange{
					Min: int64(1700000000),
				},
			},
		},
	},
	IncludeLateRegistrants: true,
	ExpireAt:               time.Now().Add(24 * time.Hour).Unix(),
}
_, err := cli.SendServerV2(ctx, req)
return err
```

- `IncludeLateRegistrants: true` 保留 V1 允许后注册玩家的行为；玩家仍需满足 Audience 条件。本例只有活动开始时及之后注册的玩家能收取。个人和群组请求无需增加此字段。
- `ContentFormatPlainText` 对应本例项目原有的纯文本展示方式；项目使用富文本时选择 `ContentFormatRichText`。
- 使用模板的业务将 ContentMode 显式设为对应模板模式，填写 TemplateID、TemplateArgs，省略 I18n；参数规则见[公共参数表](mail-v2.md#2-公共参数)。

## 3. 原来没有填写 TTL

V1 的 TTL 为0或负数时，服务端使用应用默认有效期；V2 的 ExpireAt=0 固定表示发送后14天。因此要先确认旧应用的默认有效期：

- 旧默认是2天：使用 `time.Now().Add(48 * time.Hour).Unix()`，不能直接填0，否则会变成14天。
- 旧默认是14天：可以省略 ExpireAt 或填0。

V1 从服务端处理时开始计算有效期；V2 显式填写的 ExpireAt 是调用方给出的过期时刻。为保持原行为，首次发送前计算，重试时保留原值。已经过期的请求会报错，同 Serial 重试也一样。
