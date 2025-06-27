package consts

import "time"

const (
	HttpDefaultTimeout = 30 * time.Second // HTTP请求默认超时时间
)

// 平台地址
const (
	PlatformRegionCN     = "CN"
	PlatformRegionGlobal = "Global"

	PlatformGlobalApiServer = "https://platform.shimmergames.net"
	PlatformCNApiServer     = "https://platform.shimmergames.com"
)

// Http Header 相关
const (
	HeaderAccept = "Accept"

	HeaderUserAgentValue = "ShimmerGames"

	HeaderContentType                = "Content-Type"
	HeaderContentTypeApplicationJSON = "application/json"

	HeaderServiceSDKVersion      = "X-Server-Version"
	HeaderServiceSDKVersionValue = "v1.0.0"

	HeaderUserAgent       = "User-Agent"
	HeaderUserAgentFormat = "ShimmerSDK-Go %s %s"

	HeaderAuthorization      = "Authorization"
	HeaderAuthorizationValue = "organization_id=\"%s\",appid=\"%s\",nonce=\"%s\",timestamp=\"%d\",signature=\"%s\""

	HeaderAppId          = "X-App-Id"
	HeaderOrganizationId = "X-Organization-Id"
)

const (
	SignatureFormat = "%s%d%s%s" // 数字签名原文格式
)
