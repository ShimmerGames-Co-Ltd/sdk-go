# sdk-go v3 接入指南

读者：游戏服 Go 服务端。版本 v3.5.0。

## 1. 安装

```bash
go get github.com/ShimmerGames-Co-Ltd/sdk-go/v3@v3.5.0
```

| | v1 | v3 |
|--|--|--|
| module | `github.com/ShimmerGames-Co-Ltd/sdk-go` | `github.com/ShimmerGames-Co-Ltd/sdk-go/v3` |
| import | `.../sdk-go/core` | `.../sdk-go/v3/core` |
| 分支 | `version-1.1.0` | `main` |

本地未发布 tag 时：

```
replace github.com/ShimmerGames-Co-Ltd/sdk-go/v3 => /path/to/sdk-go
```

工作树必须在 **main**，且 `go.mod` 第一行是 `module github.com/ShimmerGames-Co-Ltd/sdk-go/v3`。

## 2. 创建 Client

`WithURL` 填**网关根**，不要自行再拼 `/auth`。

```go
c, err := core.NewClient(
    core.WithURL("http://platform-api.shimmer.lab"),
    core.WithAppID(appID),
    core.WithOrganizationID(orgID),
    core.WithAuthSecret(authSecret),
)
```

默认带 `X-Server-Version: v3.5.0`。不要设置 `X-SDK-Version: v1.2.0`。密钥只从环境变量或密钥系统读取。

同一 `*core.Client` 可交给 `auth.New`、`mail.New`、`iap.New`、`leaderboard.New`、`redeemcode.New`。这些包要求 `WithAuthSecret`。

## 3. 新 body 与错误

服务端成功和业务失败都是 **HTTP 200**：

| 情况 | 判断 |
|------|------|
| 成功 | `code==0`，数据在 `data` |
| 业务失败 | `*core.APIError`（`code!=0`） |
| 网关/传输非 200 | `*core.HTTPError` |
| 网络失败 | `*core.TransportError` |
| 200 但不是 `{code,data,message}` | `*core.ResponseBodyError` |

不要把 HTTP 状态码当作唯一成败标准。

## 4. 路径与 Kong

SDK 路径已经包含网关服务前缀。Kong `strip_path` 之后，服务用 `X-Original-URI` 验签，因此签的是带前缀的 URI。

| 包 | 示例 path |
|----|-----------|
| auth | `/auth/v1/server/verify` |
| mail | `/mail/v1/mail/server/login` |
| iap | `/iap/v1/server/verify_order` |
| leaderboard | `/leaderboard/v1/server/register` |
| redeemcode | `/redeemcode/v1/server/codes:check` |

`/mail/v1/mail/...` 是预期形态：前缀 `/mail`，服务自己的路径仍是 `/v1/mail/...`。

直连没有前缀的 upstream 时，改用该服务的 base URL，或用 `core.WithPathPrefix` 调整签算 URI。不要由调用方设置 `X-Original-URI`。

`officialweb` 挂在**游戏服**上，路径是 `/payment/*`，不走上表前缀。

## 5. 按包

### auth

```go
cli, _ := auth.New(c)
rep, err := cli.Verify(ctx, userToken)
```

`Verify` 的 body 另有 HMAC-SHA1（`appid`/`token`/`ts`），与上面的 Authorization 不是同一套算法。

### mail

请求 JSON 为 snake_case。平台模板发信：

```go
cli, _ := mail.New(c)
_, err = cli.SendPlayer(ctx, mail.SendPlayerRequest{
    ServerID: "1", To: "player-1", Serial: "s1", From: "sys",
    ContentMode: mail.ContentModePlatformTemplate,
    TemplateID: "welcome", TemplateArgs: `{"name":"a"}`,
})
```

路径：`/mail/v1/mail/server/login|sync|list|get|read|claim|remove` 与 `send/player|group|server`。不要调用旧的 `/mail/send/*`。

### iap

`VerifyOrder` → `/iap/v1/server/verify_order`。`RecordPurchase` → `/iap/v1/server/record_purchase`。不含玩家创单和渠道回调。

### leaderboard

JSON 键与 proto 字段名一致（如 `subId`、`startAt`），不是 snake_case。可 `Register` 后 `SetScore` / `GetScore` 自举一条榜。

### redeemcode

`Check` 不消耗次数；`Redeem` 会核销。

### officialweb

游戏服实现 `PaymentHandler`，`Inbound.Register` 挂上 `/payment/lookup_role`、`get_goods_list`、`pre_check`、`add_order`、`buypayment` 与 `/healthz`。验签算法是 `SignSortedQSMD5`。业务状态机留在宿主。

## 6. JSON

| 范围 | 键名 |
|------|------|
| Mail / Auth / IAP / Redeem 请求与成功 data | snake_case |
| Leaderboard | proto 原名（camelCase） |
| 错误 | `code`、`message` |

Go 字段仍是 PascalCase。

## 7. 从 v1 迁移

- import 换成 `sdk-go/v3/...`。
- 不再使用 `APIResult` / `PlatError`；用 `errors.As` 区分 `APIError` 与 `HTTPError`。
- 版本头改为默认 `v3.5.0`（不要再发 `v1.0.0`）。
- 路径改为第 4 节，不要继续拼 `/mail/send/player`。
- 成功体从扁平 JSON 改为只读 `data`。

## 8. 密钥

示例与仓库只放环境变量名。不要把 `auth_secret`、用户 token 写进代码或提交到 git。

## 9. 联调

中台夹具 mockgame 提供 UAT smoke（auth、mail、排行榜自举、固定订单验单、兑换码 Check）。环境变量模板见夹具 `configs/sdk_v3_smoke.env.example`。真实值放本地 gitignore 文件。

夹具可以同时链 v1 与 v3，出站用配置切换；默认 v3。官充入站只有 v3。

## 10. 非目标

无 chat、无 realtime、无 Admin、无玩家邮件 JWT 面、无 IAP 用户创单。
