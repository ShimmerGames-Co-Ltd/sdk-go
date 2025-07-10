package authorization

type (
	// UserAuthorizeRequest
	// Token 登录凭据
	// Timestamp 时间戳
	// Sign 验签数据
	UserAuthorizeRequest struct {
		Token     string `json:"token"`
		Timestamp int64  `json:"ts"`
		Sign      string `json:"sign"`
	}

	// UserAuthorizeResponse
	// UserId 平台用户ID, 多应用下同一个ID
	// RoleId 平台角色ID, 当前应用下唯一ID
	// ShortId 平台角色短ID, 当前应用下唯一ID
	// DispatchServer 用于服务器分发(未接入可忽略)
	// RemoteAddr 用户访问IP地址
	UserAuthorizeResponse struct {
		UserId         int64  `json:"user_id"`
		RoleId         int64  `json:"role_id"`
		ShortId        int32  `json:"open_id"` // Alias: 平台OpenID 等同于 角色短ID
		DispatchServer string `json:"dispatch_server"`
		RemoteAddr     string `json:"remote_addr"`
	}
)
