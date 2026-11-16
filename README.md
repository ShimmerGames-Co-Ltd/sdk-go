# sdk-go v3

游戏服接入中台（Shimo）的 Go SDK。module：`github.com/ShimmerGames-Co-Ltd/sdk-go/v3`。

完整说明：[docs/integration.md](docs/integration.md)。

## 快速开始

```bash
go get github.com/ShimmerGames-Co-Ltd/sdk-go/v3@latest
```

以下代码放在返回 `error` 的业务函数中；`ctx` 和配置值由项目提供。

```go
c, err := core.NewClient(
	core.WithURL(os.Getenv("SHIMO_URL")), // 网关根地址，不含 /auth 等服务路径
	core.WithAppID(appID),
	core.WithOrganizationID(orgID),
	core.WithServerSignSecret(secret), // Hub app.auth_secret
)
if err != nil {
	return err
}
cli, err := auth.New(c)
if err != nil {
	return err
}
_, err = cli.Verify(ctx, userToken)
return err
```

- 签名和兼容请求头由 SDK 自动处理。
- 调用成功时 `err == nil`，SDK 已从响应的 `data` 中解析返回值。
- 业务失败返回 `*core.APIError`；HTTP 非 200 返回 `*core.HTTPError`。

## SDK 版本

这里的 v1/v3 指 **Go module 版本**。邮件发送 V1/V2 是另一套接口版本，两者都在 `sdk-go/v3/mail` 包中。

| | v1 `github.com/ShimmerGames-Co-Ltd/sdk-go` | v3 `…/sdk-go/v3` |
|--|--|--|
| 响应 | 扁平；失败看非 2xx | HTTP 200 + `{code,data,message}` |
| 路径 | `/mail/send/*` 等 | `/auth/v1/...`、`/mail/v1/mail/...`、`/mail/v2/mail/...` 等 |

新项目使用 v3；旧 module 的迁移见[接入指南](docs/integration.md#7-旧-sdk-module-迁移)。

## 包

`core` `auth` `mail` `iap` `leaderboard` `redeemcode` `officialweb`。

不含 chat、realtime、Admin。

邮件新接入使用 `SendPlayerV2`、`SendGroupV2`、`SendServerV2`；现有 V1 发送方法继续可用。服务端需支持 V2 接口。字段、默认值和示例见[邮件 V2 接入](docs/mail-v2.md)。

接收新邮 HTTP 回调使用 `mail.NewHTTPNotifyHandler`；已有路由可调用 `mail.ParseHTTPNotify`。SDK 负责验签和单封事件解析，项目负责处理或入队。见[新邮通知接入](docs/mail-notify.md)。
