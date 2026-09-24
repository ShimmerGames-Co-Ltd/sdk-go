package redeemcode

import (
	"context"
	"net/http"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	pathRedeem = "/redeemcode/v1/server/codes:redeem"
	pathCheck  = "/redeemcode/v1/server/codes:check"
)

// Client 封装 RedeemServer 代玩家兑换/检查。不封装用户 JWT 兑换与 Admin。
type Client struct {
	core *core.Client
}

func New(c *core.Client) (*Client, error) {
	if c == nil {
		return nil, &core.ConfigError{Msg: "core.Client 不能为空"}
	}
	if c.AuthMode() != core.AuthServerSignature {
		return nil, &core.ConfigError{Msg: "RedeemServer 需要 WithAuthSecret"}
	}
	return &Client{core: c}, nil
}

func (c *Client) Redeem(ctx context.Context, req RedeemRequest) (*RedeemReply, error) {
	var out RedeemReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathRedeem, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Check(ctx context.Context, req CheckRequest) (*CheckReply, error) {
	var out CheckReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathCheck, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
