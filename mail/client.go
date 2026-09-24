package mail

import (
	"context"
	"net/http"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	pathLogin      = "/mail/v1/mail/server/login"
	pathSync       = "/mail/v1/mail/server/sync"
	pathList       = "/mail/v1/mail/server/list"
	pathGet        = "/mail/v1/mail/server/get"
	pathRead       = "/mail/v1/mail/server/read"
	pathClaim      = "/mail/v1/mail/server/claim"
	pathRemove     = "/mail/v1/mail/server/remove"
	pathSendPlayer = "/mail/v1/mail/server/send/player"
	pathSendGroup  = "/mail/v1/mail/server/send/group"
	pathSendServer = "/mail/v1/mail/server/send/server"
)

// Client 封装 MailServer（SERVER_SIGNATURE）。不封装 MailUser / MailLegacy。
type Client struct {
	core *core.Client
}

func New(c *core.Client) (*Client, error) {
	if c == nil {
		return nil, &core.ConfigError{Msg: "core.Client 不能为空"}
	}
	if c.AuthMode() != core.AuthServerSignature {
		return nil, &core.ConfigError{Msg: "MailServer 需要 WithAuthSecret"}
	}
	return &Client{core: c}, nil
}

func (c *Client) Login(ctx context.Context, req LoginRequest) (*LoginReply, error) {
	var out LoginReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathLogin, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Sync(ctx context.Context, req SyncRequest) (*SyncReply, error) {
	var out SyncReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSync, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) List(ctx context.Context, req ListRequest) (*ListReply, error) {
	var out ListReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Get(ctx context.Context, req GetRequest) (*GetReply, error) {
	var out GetReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathGet, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Read(ctx context.Context, req ReadRequest) (*ReadReply, error) {
	var out ReadReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathRead, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Claim(ctx context.Context, req ClaimRequest) (*ClaimReply, error) {
	var out ClaimReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathClaim, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Remove(ctx context.Context, req RemoveRequest) (*RemoveReply, error) {
	var out RemoveReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathRemove, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendPlayer(ctx context.Context, req SendPlayerRequest) (*SendPlayerReply, error) {
	var out SendPlayerReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendPlayer, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendGroup(ctx context.Context, req SendGroupRequest) (*SendGroupReply, error) {
	var out SendGroupReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendGroup, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendServer(ctx context.Context, req SendServerRequest) (*SendServerReply, error) {
	var out SendServerReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendServer, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
