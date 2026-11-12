# sdk-go v3 接入指南

读者：游戏服 Go 服务端。module：`github.com/ShimmerGames-Co-Ltd/sdk-go/v3`。

## 1. 安装

```bash
go get github.com/ShimmerGames-Co-Ltd/sdk-go/v3@latest
```

| | 旧 SDK module（v1） | 当前 SDK module（v3） |
|--|--|--|
| module | `github.com/ShimmerGames-Co-Ltd/sdk-go` | `github.com/ShimmerGames-Co-Ltd/sdk-go/v3` |
| import | `.../sdk-go/core` | `.../sdk-go/v3/core` |

依赖版本记录在游戏服项目的 `go.mod` / `go.sum` 中。需要使用本地 SDK 源码时，在游戏服项目的 `go.mod` 中添加：

```
replace github.com/ShimmerGames-Co-Ltd/sdk-go/v3 => /path/to/sdk-go
```

将 `/path/to/sdk-go` 替换为源码目录；恢复使用发布版本时移除这行。

## 2. 创建 Client

`WithURL` 填**网关根地址**，例如 `http://platform-api.shimmer.lab`，不含 `/auth`、`/mail` 等服务路径。以下代码放在返回 `error` 的业务函数中；`ctx` 和配置值由项目提供。

```go
c, err := core.NewClient(
	core.WithURL(os.Getenv("SHIMO_URL")),
	core.WithAppID(appID),
	core.WithOrganizationID(orgID),
	core.WithServerSignSecret(serverSignSecret), // Hub app.auth_secret
)
if err != nil {
	return err
}
```

`WithURL`、`WithAppID`、`WithServerSignSecret` 必填；`WithOrganizationID` 与 `WithHTTPClient` 可选。签名密钥使用 Hub 的 `app.auth_secret`，不要使用客户端的 `sdk_secret`。

SDK 自动处理签名和兼容请求头。密钥从环境变量或密钥系统读取。

同一 `*core.Client` 可用于创建各业务 Client，例如 `auth.New(c)`、`mail.New(c)`。

## 3. 返回值与错误

先检查调用返回的 `err`；成功时直接使用返回的 Go 对象。SDK 已解析服务端的 `{code,data,message}` 响应：

| 情况 | 判断 |
|------|------|
| 成功 | `err == nil`；服务端 `code==0`，SDK 将 `data` 解析为返回对象 |
| 业务失败 | `*core.APIError`（`code!=0`） |
| HTTP 非 200 | `*core.HTTPError` |
| 网络失败 | `*core.TransportError` |
| HTTP 200，但响应不符合 `{code,data,message}` 格式 | `*core.ResponseBodyError` |

业务失败也可能返回 HTTP 200，不能只凭 HTTP 状态码判断成功。需要区分错误类型时使用 `errors.As`。

## 4. 请求路径

SDK 自动拼接以下完整路径，调用方只需配置根地址。

| 包 | 示例 path |
|----|-----------|
| auth | `/auth/v1/server/verify` |
| mail | `/mail/v1/mail/server/login`、`/mail/v2/mail/server/send/player` |
| iap | `/iap/v1/server/verify_order` |
| leaderboard | `/leaderboard/v1/server/register` |
| redeemcode | `/redeemcode/v1/server/codes:check` |

`/mail/v1/mail/...` 中的两个 `mail` 都需要保留。路径与验签由 SDK 处理，无需自行设置 `X-Original-URI`。

`officialweb` 挂在**游戏服**上，路径是 `/payment/*`，不走上表前缀。

## 5. 各包用法

### auth

```go
cli, err := auth.New(c)
if err != nil {
	return err
}
reply, err := cli.Verify(ctx, userToken)
if err != nil {
	return err
}
fmt.Println(reply.RoleID)
return nil
```

`c` 为第 2 节创建的 Client；`reply` 为解析后的验证结果。

### mail

#### V2：新接入使用

使用 `SendPlayerV2`、`SendGroupV2`、`SendServerV2`。V2 指邮件发送接口版本，仍使用 `sdk-go/v3` module；服务端须支持 V2 接口。

[参数表、必填性与默认值](mail-v2.md) · [V1 迁移对照](mail-v2.md#7-邮件发送-v1-到-v2-迁移) · [完整迁移示例](mail-v1-to-v2.md)。

读取邮箱继续使用 `Login`、`Sync`、`List`、`Get` 等方法。游戏服读写玩家邮箱时，后台须将 App 的 `runtime_access` 配置为 `via_game_server`；发送邮件不受此配置限制。使用 `client_direct` 的项目由客户端 SDK 读取邮箱。

#### V1：已有接入仍可使用

`SendPlayer`、`SendGroup`、`SendServer` 的原用法保留。平台模板模式下，`TemplateID` 填后台模板的 `code`，`Title`、`Content` 留空。

V1 的 `TemplateArgs` 是 **包含 JSON 对象的字符串**，键名与模板中的 `{{key}}` 一致：

- `plain.*` 的值为字符串、数字或布尔值，各语言共用。
- `i18n.*` 的值为“语言标识 → 值”的对象，必须含 App 默认语言；每种语言的值也是字符串、数字或布尔值。

示例：后台已启用 `welcome` 模板，恰好使用 `{{plain.player_name}}` 和 `{{i18n.item_name}}`，App 默认语言为 `English`。

```go
cli, err := mail.New(c)
if err != nil {
	return err
}
_, err = cli.SendPlayer(ctx, mail.SendPlayerRequest{
	From:         "system",
	Serial:       "welcome-role-10001",
	ServerID:     "1",
	To:           "10001",
	ContentMode:  mail.ContentModePlatformTemplate,
	TemplateID:   "welcome",
	TemplateArgs: `{"plain.player_name":"Ada","i18n.item_name":{"English":"Legendary Sword"}}`,
})
return err
```

参数须与模板要求的键完全一致，不能缺传或多传，JSON 字符串最多 8192 字节。`i18n.*` 按玩家登录语言取值，缺少该语言时使用 App 默认语言。

### iap

`VerifyOrder` → `/iap/v1/server/verify_order`。`RecordPurchase` → `/iap/v1/server/record_purchase`。不含玩家创单和渠道回调。

### leaderboard

JSON 键使用 camelCase（如 `subId`、`startAt`）。先用 `Register` 创建榜单，再用 `SetScore` / `GetScore` 写入和查询分数。

### redeemcode

`Check` 不消耗次数；`Redeem` 会核销。

### officialweb

游戏服实现 `PaymentHandler` 的五个回调，再用 `Inbound.Register` 注册 `/payment/*` 和 `/healthz`。

SDK 负责验签、解析 JSON，再调用对应回调。业务接口返回 HTTP 200，错误看响应的 `code`：方法错误、非法 JSON 或字段类型不符为 `400`；验签失败为 `401`。字段是否必填、是否为空由游戏服业务校验。

- `PlatformRoleID` 是字符串；从 Auth 的 `VerifyReply.RoleID` 转换时使用 `strconv.FormatInt(int64(reply.RoleID), 10)`。
- `AddOrderRequest.ProductNameI18n` 是 JSON 字符串，例如 `{"English":"Pack"}`；有值时参与验签。需要用于付后邮件时保存此值，发送 V2 邮件时先解码为对象，再填入 `TemplateArgs["i18n.product_name"]`。支付通知 `BuyPaymentRequest` 不含此字段。

完整字段见[回调请求类型](../officialweb/request.go)，回调定义见[入站接口](../officialweb/inbound.go)。

## 6. JSON

| 范围 | 键名 |
|------|------|
| Mail / Auth / IAP / Redeem 请求与成功 data | snake_case |
| Leaderboard | proto 原名（camelCase） |
| 错误 | `code`、`message` |

文档参数表和 Go 示例使用 SDK 的 Go 字段名，JSON 键名由 SDK 转换。

邮件响应的 `ReadReply.ReadMailIDs`、`RemoveReply.RemovedMailIDs` 均为 `[]int64`。SDK 自动兼容服务端返回的数字和数字字符串，无需调用方转换；空数组与 `null` 保持区分。请求中的 `MailIDs` 仍编码为 JSON 数字数组。

## 7. 旧 SDK module 迁移

本节指旧 SDK module 到 `sdk-go/v3` 的迁移。同一个 `sdk-go/v3/mail` 包内的邮件发送 V1 到 V2，请看[邮件迁移说明](mail-v2.md#7-邮件发送-v1-到-v2-迁移)。

- import 换成 `github.com/ShimmerGames-Co-Ltd/sdk-go/v3/...`。
- 不再使用 `APIResult` / `PlatError`；用 `errors.As` 区分 `APIError` 与 `HTTPError`。
- 签名及兼容请求头交由 SDK 处理，移除接入方自行设置或覆盖版本头的逻辑。
- 使用 SDK 方法发请求，无需手工拼路径。
- 使用方法返回的 Go 对象，SDK 自动从响应的 `data` 中解析。

## 8. 密钥

示例与仓库只放环境变量名。不要把 `auth_secret`、用户 token 写进代码或提交到 git。

## 9. 联调

SDK 自带[调用示例](../examples/gameserver/main.go)。先设置 `SHIMO_URL`、`SHIMO_APP_ID`、`SHIMO_AUTH_SECRET`，再选择场景运行，例如：

```bash
SHIMO_DEMO=auth SHIMO_TOKEN='<user-token>' go run ./examples/gameserver
```

邮件示例的场景和所需环境变量见[邮件 V2 调用示例](mail-v2.md#6-调用示例)。

## 10. 支持范围

SDK 面向游戏服，不包含 chat、realtime、Admin、客户端邮件接口或玩家侧 IAP 创单接口。
