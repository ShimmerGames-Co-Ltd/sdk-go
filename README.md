# sdk-go v3

游戏服接入中台（Shimo）的 Go SDK。module：`github.com/ShimmerGames-Co-Ltd/sdk-go/v3`。版本 **v3.5.0**。

完整说明：[docs/integration.md](docs/integration.md)。

## 五分钟

```bash
go get github.com/ShimmerGames-Co-Ltd/sdk-go/v3@v3.5.0
```

```go
c, err := core.NewClient(
    core.WithURL(os.Getenv("SHIMO_URL")), // 网关根，不要带 /auth
    core.WithAppID(appID),
    core.WithOrganizationID(orgID),
    core.WithAuthSecret(secret),
)
cli, err := auth.New(c)
reply, err := cli.Verify(ctx, userToken)
```

- 默认请求头 `X-Server-Version: v3.5.0`。
- 成功：HTTP 200 且 `code==0`，业务字段在 `data`。
- 业务失败：HTTP 200 且 `code!=0`，错误类型 `*core.APIError`。
- 非 200：`*core.HTTPError`。

## 与 v1

| | v1 `github.com/ShimmerGames-Co-Ltd/sdk-go` | v3 `…/sdk-go/v3` |
|--|--|--|
| 分支 | `version-1.1.0` | **`main`** |
| 响应 | 扁平；失败看非 2xx | HTTP 200 + `{code,data,message}` |
| 路径 | `/mail/send/*` 等 | `/auth/v1/...`、`/mail/v1/mail/...` 等 |

存量继续用 v1。新项目只用 v3。

## 包

`core` `auth` `mail` `iap` `leaderboard` `redeemcode` `officialweb`。

不含 chat、realtime、Admin。
