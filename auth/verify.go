package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	pathServerVerify = "/auth/v1/server/verify"
	pathQueryForSolo = "/auth/v1/server/query4solo"
)

// Client 封装游戏服 ServerVerify（HMAC-SHA1 body.sign，外加 SERVER_SIGNATURE Authorization）。
type Client struct {
	core *core.Client
	now  func() time.Time
}

func New(c *core.Client) (*Client, error) {
	if c == nil {
		return nil, &core.ConfigError{Msg: "core.Client 不能为空"}
	}
	if c.AuthMode() != core.AuthServerSignature {
		return nil, &core.ConfigError{Msg: "UserVerify 需要 WithAuthSecret"}
	}
	return &Client{core: c, now: time.Now}, nil
}

type VerifyReply struct {
	UserID         core.JSONInt64  `json:"user_id"`
	RoleID         core.JSONInt64  `json:"role_id"`
	OpenID         int32  `json:"open_id"`
	DispatchServer string `json:"dispatch_server"`
	RemoteAddr     string `json:"remote_addr"`
}

type verifyBody struct {
	Token string `json:"token"`
	Ts    core.JSONInt64  `json:"ts"`
	Sign  string `json:"sign"`
}

// Verify 校验玩家 token。open_id 字段实际是 short_id（服务端历史口径）。
func (c *Client) Verify(ctx context.Context, token string) (*VerifyReply, error) {
	return c.VerifyAt(ctx, token, c.now().Unix())
}

// VerifyAt 使用固定时间戳，便于对照签名向量。
func (c *Client) VerifyAt(ctx context.Context, token string, ts int64) (*VerifyReply, error) {
	sign := core.UserVerifySign(c.core.AppID(), token, c.core.AuthSecret(), ts)
	body := verifyBody{Token: token, Ts: core.JSONInt64(ts), Sign: sign}
	var out VerifyReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathServerVerify, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

type SoloReply struct {
	UID     string `json:"uid"`
	Channel string `json:"channel"`
	UserID  core.JSONInt64  `json:"user_id"`
	ShortID int32  `json:"short_id"`
	AppID   string `json:"appid"`
}

// QueryForSolo 按 token 查询 Solo 用户。仍走已签名 core.Client（proto 上该 RPC 历史未强制验签）。
func (c *Client) QueryForSolo(ctx context.Context, token string) (*SoloReply, error) {
	var out SoloReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathQueryForSolo, map[string]string{"token": token}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
