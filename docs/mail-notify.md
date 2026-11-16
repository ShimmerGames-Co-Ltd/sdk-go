# 新邮 HTTP 通知接入

游戏服接收 Mail 的新邮事件，V1/V2 发信共用这套通知。后台可配置默认通知及按服覆盖；选择 HTTP 时，将通知地址填写为项目实际注册的 URL。

**一次回调是一封邮件事件。** `Headers` 沿用 NATS 的数组格式，HTTP 接收要求恰好一项。群组发送会逐玩家通知，区服发送会逐区服通知；不把多个玩家或多封邮件合并为一次回调。

## 1. 注册接收路由

以下片段放在返回 `error` 的服务初始化函数中，`mux` 为已有的 `*http.ServeMux`。示例只输出事件，实际项目应在回调中完成处理或入队后再返回 `nil`。

```go
receiver, err := mail.NewHTTPNotifyHandler(
	os.Getenv("SHIMO_APP_ID"),
	os.Getenv("SHIMO_AUTH_SECRET"), // Hub App.auth_secret
	func(ctx context.Context, r *http.Request, notification mail.HTTPNotifyRequest) error {
		if notification.PlayerMail != nil {
			message := notification.PlayerMail
			fmt.Printf("server=%s player=%s mailbox_id=%d\n",
				message.ServerID, message.PlayerID, message.Headers[0].ID)
		} else {
			message := notification.ServerMail
			fmt.Printf("server=%s reference=%s\n", message.ServerID, message.Headers[0].MailRefID)
		}
		return nil
	},
)
if err != nil {
	return err
}
mux.Handle("/mail/new_mail_notify", receiver)
return nil
```

SDK 检查 POST、App、原始正文签名、1 MiB 正文上限和单封事件格式（Headers 恰好一项且 ID>0），通过后才调用业务函数。返回 `nil` 时响应 HTTP 204；处理失败应返回错误。

需要动态读取配置或使用已有路由时，直接解析当前请求：

```go
notification, err := mail.ParseHTTPNotify(r, appID, authSecret)
if err != nil {
	mail.WriteHTTPNotifyError(w, err)
	return
}
fmt.Printf("verified notification from app=%s\n", notification.AppID)
w.WriteHeader(http.StatusNoContent)
```

解析会消费请求体。项目可使用 `notification.RawBody` 留存验签原文，无需再次读取或重新编码。

## 2. SDK 字段

所有通知类型均在 [mail/types.go](../mail/types.go)，接收逻辑在 [mail/http_notify.go](../mail/http_notify.go)。

| 类型 | 字段 | Go 类型 | 含义 |
|---|---|---|---|
| HTTPNotifyRequest | NewMailNotify | NewMailNotify | 内嵌的通知载荷，可直接访问 PlayerMail / ServerMail |
| HTTPNotifyRequest | AppID | string | 已验证的 X-App-ID，不写入 JSON |
| HTTPNotifyRequest | RawBody | []byte | 用于验签的原始 JSON，不写入 JSON |
| NewMailNotify | PlayerMail | *NewPlayerMailNotify | 个人/群组邮件事件；与 ServerMail 二选一 |
| NewMailNotify | ServerMail | *NewServerMailNotify | 区服邮件事件；与 PlayerMail 二选一 |
| NewPlayerMailNotify | ServerID | string | 实际区服 ID |
| NewPlayerMailNotify | PlayerID | string | 该区服的邮箱玩家标识 |
| NewPlayerMailNotify | Headers | []MailHeader | 恰好一项；ID 是 Get/Read/Claim 使用的玩家邮箱内 ID |
| NewServerMailNotify | ServerID | string | 实际区服 ID，不使用全服通配符 |
| NewServerMailNotify | Headers | []ServerMailNotifyHeader | 恰好一项；ID 是服邮箱游标，MailRefID 是正文引用，均不能当玩家邮箱内 ID |

收到区服事件后，项目仍需为相应玩家执行 Sync，再用 List/Get 读取邮箱。通知不包含所有收件玩家，也不等于玩家已经收取邮件。

## 3. 错误与投递语义

| 情况 | HTTP 状态 |
|---|---|
| 验签并处理成功 | 204 |
| 非 POST | 405，Allow: POST |
| App 不符、缺签名、错签名或正文被改动 | 401 |
| 无效 JSON、消息种类错误、区服/玩家为空、Headers 不为一项或 ID 非正数 | 400 |
| 正文超过 1 MiB | 413 |
| 配置缺失或业务处理失败 | 503 |

协议错误类型为 `*mail.HTTPNotifyError`，包含 `StatusCode` 和 `Message`；配置错误为 `*core.ConfigError`。普通业务错误响应不包含原始错误详情。已有框架可用 `WriteHTTPNotifyError` 统一响应，不能包装成 HTTP 200 的业务错误码。

签名使用 `X-App-ID`、`X-Mail-Signature` 请求头，以及 `Base64(HMAC-SHA256(App.auth_secret, 原始 JSON 字节))`。这些检查由 SDK 完成，调用方不需要自行转换载荷或计算签名。

通知是异步尽力提示，不提供可靠送达或顺序保证。个人事件可按 AppID＋ServerID＋PlayerID＋邮件 ID 去重，区服事件按 AppID＋ServerID＋服邮箱游标去重。登录、重连等流程仍通过 Sync/List 获取真实邮箱状态；不要按回调次数累加未读数，以返回的 Badge 为准。SDK 不自动 Sync、不重试、不去重，也不替项目提前入队确认。
