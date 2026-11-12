# 邮件 V2 接入

[安装与 Client 创建](integration.md)。V2 指邮件发送接口版本，仍使用 `sdk-go/v3/mail` 包；服务端需支持 V2 接口。签名由 SDK 自动处理。

**必填 From、Serial、收件人；内容选择普通文案或模板。** 参数表使用 SDK 的 Go 字段名。已接入 V1 的项目先看[迁移对照](#7-邮件发送-v1-到-v2-迁移)。

## 1. 接口

| SDK 方法 | HTTP 路径（POST） | 收件人 |
|---|---|---|
| SendPlayerV2 | /mail/v2/mail/server/send/player | 一个区服的一个玩家 |
| SendGroupV2 | /mail/v2/mail/server/send/group | 按区服分组的指定玩家 |
| SendServerV2 | /mail/v2/mail/server/send/server | 指定区服或 App 全服；玩家通过 Sync 收取 |

## 2. 公共参数

以下字段在三种请求中相同。Required＝必填；Optional＝可选；Conditional＝满足条件时必填。

| 字段 | Go 类型 | 必填性 | 默认值 | 说明 |
|---|---|---|---|---|
| AppID | string | Optional | 创建 Client 时的 AppID | 通常省略；填写时须与 Client 一致 |
| From | string | Required | — | 发件人，最多128字节 |
| Serial | string | Required | — | 业务唯一编号，最多128字节；重试不变 |
| ContentMode | string | Optional | plain | 内容模式，见下表 |
| I18n | []I18nBlock | Conditional | 无 | 普通文案必填；模板模式省略 |
| ContentCategory | string | Optional | 空 | 项目自定义分类，最多64字节；无需后台登记，不继承模板分类 |
| ContentFormat | string | Optional | 见下表 | plain_text / rich_text；声明展示格式，不转换正文 |
| Attachments | []MailAttachment | Optional | 无附件 | 本次奖励，不从模板继承；JSON 编码后最多4096字节 |
| ExpireAt | int64 | Optional | 发送后14天 | 过期时刻的 Unix 秒；0用默认值，非零须晚于服务端处理时刻 |
| TemplateID | string | Conditional | 空 | 模板模式必填，最多128字节 |
| TemplateArgs | map[string]any | Conditional | 无 | 平台模板有占位符时须完整提供；客户端模板按项目约定 |

`ContentMode` 使用下列 `mail.ContentMode*` 常量：

| ContentMode | I18n | TemplateID | TemplateArgs | ContentFormat 默认值 |
|---|---|---|---|---|
| ContentModePlain | 必填 | 省略 | 省略 | rich_text |
| ContentModePlatformTemplate | 省略 | 后台已启用模板的 code | 与模板占位符完全一致 | 模板格式；显式值须与模板一致 |
| ContentModeClientTemplate | 省略 | 项目客户端模板 ID，不与后台 code 冲突 | 按项目约定传给客户端 | rich_text |

`ContentFormat` 可使用 `mail.ContentFormatPlainText` 或 `mail.ContentFormatRichText`。

## 3. 收件人参数

| 请求类型 | 字段 | Go 类型 | 必填性 | 说明 |
|---|---|---|---|---|
| SendPlayerV2Request | ServerID | string | Required | 实际区服 ID，最多64字节 |
| SendPlayerV2Request | PlayerID | string | Required | 邮箱玩家标识，最多64字节；来源见下表 |
| SendGroupV2Request | RecipientTargets | []RecipientTarget | Required | 至少一个区服，每个区服只列一次 |
| SendServerV2Request | ServerIDs | []string | Required | 非空、不重复，各最多64字节；全服填 `[]string{"*"}`，不能与具体区服混用 |
| SendServerV2Request | Audience | *Audience | Optional | nil 表示不加筛选；填写时只筛选 ServerIDs 范围内的玩家；不能传空结构体 |
| SendServerV2Request | IncludeLateRegistrants | bool | Optional | 默认 false；true 允许发送后注册的玩家收取，仍须满足 Audience |

PlayerID 与读取邮箱时的标识一致，按后台邮箱配置选择：

| 邮箱配置 | PlayerID 取值 |
|---|---|
| client_direct + user_token（默认） | Auth 返回的 RoleID，转为十进制字符串 |
| client_direct + client_body | 项目约定的自定义标识 |
| via_game_server | 项目约定的邮箱标识，与 Login/Sync 一致 |

个人和群组请求填写实际区服 ID，不使用 `"*"`。

**区服邮件默认只允许发送前已注册的玩家收取。** `IncludeLateRegistrants=false` 时，Sync 的 `PlayerAttrs["register_at"]` 须严格早于发送时刻，即使未填 Audience 也一样。属性缺失或不满足条件时，Sync 可能成功但收不到此邮件。个人和群组不使用这项策略。

## 4. 子对象

### RecipientTarget

| 字段 | Go 类型 | 必填性 | 说明 |
|---|---|---|---|
| ServerID | string | Required | 实际区服 ID，最多64字节，不使用 `"*"` |
| PlayerIDs | []string | Required | 至少一个玩家；标识非空、同服不重复，各最多64字节 |

### I18nBlock

| 字段 | Go 类型 | 必填性 | 说明 |
|---|---|---|---|
| Lang | string | Required | App 配置的语言标识，如 English；有文案的语言须已启用，不能重复 |
| Title | string | Conditional | 默认语言必填；其他语言与 Content 同时填写；默认语言标题最多256字节 |
| Content | string | Conditional | 默认语言必填；其他语言与 Title 同时填写 |
| SubTitle | string | Optional | 副标题，最多80个 Unicode 字符 |

未使用的语言省略。缺少玩家语言的完整文案时，标题、正文和副标题一起使用默认语言。

### MailAttachment

| 字段 | Go 类型 | 必填性 | 默认值与说明 |
|---|---|---|---|
| Type | int32 | Optional | 0；项目约定的奖励类型，须≥0 |
| ID | string | Required | 非空业务道具 ID；发送时不检查后台道具登记 |
| Count | int32 | Required | 奖励数量，须>0 |
| Param | string | Optional | 空；项目扩展字符串，原样传递 |

### TemplateArgs 的键和值

平台模板的 map 键与模板中的 `{{key}}` 一致，不能缺传或多传，JSON 编码后最多8192字节：

- `plain.*`：字符串、数字或布尔值，例如 `"plain.player_name": "Ada"`，各语言共用。
- `i18n.*`：“语言标识 → 值”的 map，例如 `"i18n.item_name": map[string]any{"English": "Sword"}`，必须含 App 默认语言；每种语言的值也是字符串、数字或布尔值。

发送时校验模板要求的参数名和类型；参数目录由后台维护。客户端模板参数按项目约定原样传递。

### Audience

Audience 是筛选条件树：最外层为条件组，至少包含一条规则。`and` 为全部满足，`or` 为满足任一。

| 类型 | 字段 | Go 类型 | 必填性 | 说明 |
|---|---|---|---|---|
| Audience | Op | string | Required | AudienceOpAnd / AudienceOpOr |
| Audience | Children | []AudienceNode | Required | 规则或嵌套条件组 |
| AudienceNode | Op | string | Required | 组使用 and/or；规则使用下表的比较操作符 |
| AudienceNode | Key | string | Conditional | 规则必填，填写可筛选字段名；组省略 |
| AudienceNode | Value | any | Conditional | 规则必填，类型见下表；组省略；0、false、空字符串会保留 |
| AudienceNode | Children | []AudienceNode | Conditional | 嵌套组必填；规则省略 |
| AudienceRange | Min | any | Required | 包含的下界；非负数字或版本字符串，0有效 |
| AudienceRange | Max | any | Optional | 不包含的上界；省略表示无上界，填写时须大于 Min |

Key、值类型和可用操作符以后台配置的 **Audience Catalog（可筛选字段列表）** 为准。Key 对应 `SyncRequest.PlayerAttrs` 的同名属性。

| 操作符（使用 mail.AudienceOp* 常量） | Value 的写法 |
|---|---|
| eq | 单个字符串、数字或布尔值，按字段类型填写 |
| in | 非空的同类型切片，如 `[]string{"CN", "US"}` |
| range | AudienceRange，包含 Min、不包含 Max；register_at 默认使用此操作 |
| gte / lte / gt / lt | 可比较的单个值；版本字段填写版本字符串 |

`Value` 按字段类型传 Go 数值、字符串或 bool；区间边界只用数字或版本字符串。`any` 用于容纳这些类型，不能一律转成字符串。完整构造见[迁移示例](mail-v1-to-v2.md#2-迁移后v2)。

### SyncRequest.PlayerAttrs 怎么填

以下键填入 `PlayerAttrs` 这个 map，按规则需要提供：

- `register_at`：玩家真实注册 Unix 秒，游戏服上传，不从 Login 自动补齐。
- `server_id`：可省略，服务端使用 Sync 顶层的 ServerID；已填值优先。
- `package_name`、`client_version`、`country`：由游戏服提供，不从请求头或用户 token 补齐。
- `distribution_channel`、`platform_type`：缺省时尝试从 package_name 对应的 Hub 包配置补齐，已有值不覆盖；补齐失败不阻断 Sync。
- 自定义属性：按 Catalog 约定提供。

规则依赖的属性仍缺失时，该规则不匹配。完整读信流程见第 6 节的 mail-receive 示例。

## 5. 返回值与重试

| 返回类型 | 字段 | 含义 |
|---|---|---|
| SendPlayerV2Reply | MailRefID、Badge | 邮件引用、玩家总角标；玩家删信后，同 Serial 重试成功时 MailRefID 可为空 |
| SendGroupV2Reply | Refs | 每项包含 ServerID、PlayerID、MailRefID、Badge；同服玩家可共享邮件引用 |
| SendServerV2Reply | Refs | 每项包含 ServerID、MailRefID；全服的 ServerID 为 `"*"`；成功表示已写入区服邮件，玩家还需 Sync 收取 |

- `MailRefID` 不能用于 Get/Read/Claim；这些方法使用 List 或通知中的 `MailHeader.ID`。
- 各发送类型分别去重。同一 App 内，个人按“区服＋玩家＋Serial”，群组和区服按“区服＋Serial”去重。同 Serial 不更新内容、过期时间或后注册玩家策略，删信后也不重新投递。新业务用新 Serial。
- 群组和多区服请求可能部分成功；失败时 SDK 返回 `nil, err`，已投递部分不回滚。重试复用原请求，包括 Serial、收件人、内容和 ExpireAt。SDK 不自动重试或降级到 V1。
- HTTP 200 也可能业务失败，SDK 会检查响应的 code；调用方检查 err 即可。详见[错误处理](integration.md#3-返回值与错误)。

读取和领取接口沿用原用法：附件及客户端模板参数仍是 JSON 字符串，需调用方解码；平台模板返回渲染后的文案，模板字段为空。旧邮件的 ContentFormat 为空时沿用项目原展示方式。Badge 表示邮箱总角标，SyncReply.NewMailCount 表示本次新增数。

## 6. 调用示例

以下代码放在返回 `error` 的业务函数中；`c` 已按[接入指南](integration.md#2-创建-client)创建，App 默认语言为 English。本例使用纯文本，未填 ExpireAt，有效期为14天。

```go
cli, err := mail.New(c)
if err != nil {
	return err
}
reply, err := cli.SendPlayerV2(ctx, mail.SendPlayerV2Request{
	From:          "system",
	Serial:        "welcome-role-10001",
	ServerID:      "1",
	PlayerID:      "10001",
	ContentFormat: mail.ContentFormatPlainText,
	I18n: []mail.I18nBlock{
		{
			Lang:    "English",
			Title:   "Welcome",
			Content: "Welcome to the game.",
		},
	},
})
if err != nil {
	return err
}
fmt.Println(reply.MailRefID)
return nil
```

更多场景见[示例源码](../examples/gameserver/mail_v2.go)，通过[示例入口](../examples/gameserver/main.go)运行。先设置 `SHIMO_URL`、`SHIMO_APP_ID`、`SHIMO_AUTH_SECRET`；`SHIMO_ORG_ID` 可选。发送还须设置 `SHIMO_MAIL_SERIAL`，重试不变，新业务换新值。

| SHIMO_DEMO 的值 | SHIMO_MAIL_SERIAL 示例值 | 额外环境变量或前置配置 |
|---|---|---|
| mail-send-v2 | welcome-role-10001 | 普通个人邮；SHIMO_PLAYER_ID |
| mail-group-v2 | group-reward-1001 | 跨服群组、附件；SHIMO_PLAYER_ID、SHIMO_SECOND_SERVER_ID、SHIMO_SECOND_PLAYER_ID；两个区服须不同 |
| mail-server-v2 | server-announcement-1001 | 全服，允许后注册玩家 |
| mail-audience-v2 | audience-gift-1001 | Catalog 的 register_at 支持 range；接收方 Sync 上传真实注册时间 |
| mail-platform-template-v2 | platform-reward-1001 | SHIMO_TEMPLATE_CODE；模板已启用，恰好使用 plain.player_name / i18n.item_name 两个占位符 |
| mail-client-template-v2 | client-reward-1001 | SHIMO_PLAYER_ID、SHIMO_CLIENT_TEMPLATE_ID；客户端按源码中的参数渲染 |
| mail-receive | 无需设置 | runtime_access=via_game_server；SHIMO_PLAYER_ID、SHIMO_REGISTER_AT；解码客户端模板时使用对应 SHIMO_CLIENT_TEMPLATE_ID，默认 local_reward |

上表为示例程序的环境变量。`SHIMO_SERVER_ID` 默认1，`SHIMO_DEFAULT_LANG` 默认 English，按项目配置设置。在 SDK 根目录运行：

```bash
SHIMO_DEMO=mail-send-v2 SHIMO_MAIL_SERIAL=welcome-role-10001 \
  SHIMO_PLAYER_ID=10001 go run ./examples/gameserver
```

不同场景使用不同 Serial，避免把普通文案和模板邮件误判为同一业务重试。新一轮测试也使用新的业务编号。

示例中的玩家、奖励、模板和筛选值按项目实际配置替换。个人邮和条件邮先构造 req，再调用 SDK；业务重试复用原 req。命令行示例每次启动都会重新计算时间，重新启动不能代替原请求重试。收信示例将已知客户端模板参数直接解码为项目结构体，数字可直接使用。

## 7. 邮件发送 V1 到 V2 迁移

项目按需将 `SendPlayer` / `SendGroup` / `SendServer` 换为对应的 `*V2` 方法，请求和返回类型同步替换。Client、签名和读取方法沿用，无需更换 module 路径；原 V1 方法仍可用。

### 字段对照

| V1 SDK 字段 | V2 SDK 字段 | 迁移方法 |
|---|---|---|
| AppID、From、Serial | AppID、From、Serial | 沿用原值 |
| SendPlayerRequest.To | SendPlayerV2Request.PlayerID | 沿用原玩家标识；ServerID 不变 |
| SendGroupRequest.ServerID、To | SendGroupV2Request.RecipientTargets | 每个区服对应一个 RecipientTarget，将 To 放入 PlayerIDs |
| Title、Content | I18n | 按语言合并标题和正文；个人邮件的原字符串放入 App 默认语言块 |
| Attachment | Attachments | JSON 字符串改为 []MailAttachment，保留原奖励值 |
| TemplateID | TemplateID | 平台模板仍填原 code，客户端模板仍填项目标识 |
| TemplateArgs | TemplateArgs | JSON 对象字符串改为 map[string]any，保留原键和值 |
| ContentMode | ContentMode | 普通文案用 ContentModePlain；模板显式选择对应模式 |
| TTL | ExpireAt | “有效期多少秒”改为“在哪一秒过期”，见下文 |
| SendServerRequest.ServerIDs | SendServerV2Request.ServerIDs | 沿用原区服范围 |
| SendServerRequest.MatchRule | SendServerV2Request.Audience | JSON 字符串改为条件结构；不筛选时省略 |
| GroupMailRef.To | GroupMailV2Ref.PlayerID | 按返回的 ServerID 和 PlayerID 关联收件人 |

### 默认行为差异

- **内容模式**：V2 未填 ContentMode 时始终是普通文案，不再根据 TemplateID 推断客户端模板。模板模式必须显式设置 ContentMode，并省略 I18n。
- **后注册玩家**：V1 区服邮件默认允许后注册玩家；V2 默认不允许。保持原行为时设置 `IncludeLateRegistrants: true`，已有 Audience 条件仍然生效。个人和群组无需增加此字段。
- **过期时间**：V1 的 TTL 填有效期秒数，0或负数使用应用默认值；V2 的 ExpireAt 填过期时刻的 Unix 秒，0固定表示发送后14天。保持原有效期的换算见[迁移示例](mail-v1-to-v2.md)。重试复用原请求，已过期的请求会报错。
- **展示格式**：V2 普通文案和客户端模板默认 rich_text。保持原展示效果时，按项目原约定显式填写 ContentFormat；平台模板省略此字段即可使用模板格式。

[完整迁移示例](mail-v1-to-v2.md)直接对照 V1/V2 请求。切换前确认服务端支持 V2，核对项目用到的字段及收取结果；无需重发历史邮件。
