package core

import "time"

const (
	// DefaultHTTPTimeout 默认 HTTP 超时。
	DefaultHTTPTimeout = 30 * time.Second
	// DefaultServerVersion 必须严格大于 v1.0.0，服务端才走新 body（HTTP 200+code）。
	// 与 sdk-go v3.5.0 / 中台版本叙事对齐。
	DefaultServerVersion = "v3.5.0"

	headerAccept          = "Accept"
	headerContentType     = "Content-Type"
	headerContentJSON     = "application/json"
	headerAuthorization   = "Authorization"
	headerAppID           = "X-App-ID"
	headerOrganizationID  = "X-Organization-ID"
	headerUserToken       = "X-User-Token"
	headerServerVersion   = "X-Server-Version"
	headerUserAgent       = "User-Agent"
	headerUserAgentFormat = "ShimmerSDK-Go %s %s"

	// authorizationValue 与旧 sdk-go / 服务端 parseSignature 一致（HMAC 头，不是 RSA）。
	authorizationValue = `organization_id="%s",appid="%s",nonce="%s",timestamp="%d",signature="%s"`

	// SignatureFormat 签算原文：URI + timestamp + nonce + body。
	SignatureFormat = "%s%d%s%s"
)
