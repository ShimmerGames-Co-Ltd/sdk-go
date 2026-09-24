package iap

import (
	"context"
	"net/http"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	pathVerifyOrder    = "/iap/v1/server/verify_order"
	pathRecordPurchase = "/iap/v1/server/record_purchase"
)

// Client 封装游戏服 IAP：验单与代录购买。不封装用户创单/pending/回调。
type Client struct {
	core *core.Client
}

func New(c *core.Client) (*Client, error) {
	if c == nil {
		return nil, &core.ConfigError{Msg: "core.Client 不能为空"}
	}
	return &Client{core: c}, nil
}

func (c *Client) VerifyOrder(ctx context.Context, req VerifyOrderRequest) (*VerifyOrderReply, error) {
	var out VerifyOrderReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathVerifyOrder, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) RecordPurchase(ctx context.Context, req RecordPurchaseRequest) error {
	return c.core.DoJSON(ctx, http.MethodPost, pathRecordPurchase, req, nil)
}
