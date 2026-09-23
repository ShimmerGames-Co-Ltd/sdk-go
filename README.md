# sdk-go v3（Shimo 游戏服 Go SDK）

> **v3.5.0**：module `github.com/ShimmerGames-Co-Ltd/sdk-go/v3`  
> 契约真源：Shimo 主仓 proto、签名规范与各服务 api-docs。游戏服 **不要** import 主仓内部包。

## 与 v1（同仓旧 module）的关系

| | v1 `github.com/ShimmerGames-Co-Ltd/sdk-go` | v3 `…/sdk-go/v3` |
|--|---------------------------------------------|------------------|
| 分支 | `version-1.1.0` / `main` | **`v3`** |
| 响应 | 扁平体；失败看非 2xx | **新 body**：HTTP 200 + `{code,data,message}` |
| 默认头 | `X-Server-Version: v1.0.0` | **`X-Server-Version: v3.5.0`** |
| 路径 | 旧 `/mail/send/*` 等 | 仅 **`/v1/...` 主路径** |
| 功能 | **冻结**，仅必要 fix | 本线演进 |

存量游戏服继续用 v1；**新项目只用 v3**。

原独立仓 `shimo-sdk-go` 已迁入本 module，请改依赖到本坐标。

## 模块（v3.5.0）

| 包 | 说明 |
|----|------|
| `core` | HMAC-SHA256 `Authorization`、新 body、GET/POST |
| `auth` | `/v1/server/verify`、`/v1/server/query4solo` |
| `mail` | MailServer `/v1/mail/server/*`（含 `content_mode` / 模板发信） |
| `leaderboard` | `/v1/server/*`（含 HistoryMemberInfo） |
| `iap` | 验单 / 代录购买 |
| `redeemcode` | 代玩家兑换/检查 |
| `officialweb` | 官充游戏入站 Server Kit |

**不包含**：`chat` / `realtime`（封存，不由本 SDK 发布）；Admin；MailLegacy。

## JSON

- **Mail / Auth / IAP / Redeem**：请求与成功 `data` 均为 **snake_case**。
- **Leaderboard**：与 proto 字段名一致（proto 本身为 camelCase，如 `subId`）。
- Go 结构体字段仍为 PascalCase。

## 安装

```bash
go get github.com/ShimmerGames-Co-Ltd/sdk-go/v3@v3.5.0
```

本地尚未 push tag 时：

```
replace github.com/ShimmerGames-Co-Ltd/sdk-go/v3 => /path/to/sdk-go
```

（工作树须在 **`v3` 分支**。）

## 同一 Client

```go
c, err := core.NewClient(
    core.WithURL(os.Getenv("SHIMO_URL")),
    core.WithAppID(appID),
    core.WithOrganizationID(orgID),
    core.WithAuthSecret(authSecret),
)
mailCli, _ := mail.New(c)
```

示例：`examples/gameserver`（`SHIMO_DEMO=mail|leaderboard|auth|iap|redeem`）。

## 错误

| 类型 | 含义 |
|------|------|
| `*core.TransportError` | 网络层 |
| `*core.HTTPError` | HTTP ≠ 200 |
| `*core.APIError` | HTTP 200 且 `code != 0` |

业务成败看 **`code`**，不要只看 HTTP 状态。

## Kong / 签算

签算原文：`URI + timestamp + nonce + body`。网关前缀用 `core.WithPathPrefix` / `WithSignRequestURI`。不要客户端设 `X-Original-URI`。

## 验收

```bash
go test ./... -count=1
go build ./...
```
