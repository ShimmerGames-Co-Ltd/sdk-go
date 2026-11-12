package mail

import (
	"context"
	"net/http"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	pathLogin  = "/mail/v1/mail/server/login"
	pathSync   = "/mail/v1/mail/server/sync"
	pathList   = "/mail/v1/mail/server/list"
	pathGet    = "/mail/v1/mail/server/get"
	pathRead   = "/mail/v1/mail/server/read"
	pathClaim  = "/mail/v1/mail/server/claim"
	pathRemove = "/mail/v1/mail/server/remove"
)

// Client 封装 MailServer（SERVER_SIGNATURE）。不封装 MailUser / MailLegacy。
// 发信不受 runtime_access 限制；Login/Sync/List/Get 等邮箱读写要求 App 配置为 via_game_server。
type Client struct {
	core *core.Client
}

func New(c *core.Client) (*Client, error) {
	if c == nil {
		return nil, &core.ConfigError{Msg: "core.Client 不能为空"}
	}
	return &Client{core: c}, nil
}

// Login 初始化邮箱并设置读取语言；切换语言会刷新邮件的本地化摘要，注册时间在 Sync 中提供。
func (c *Client) Login(ctx context.Context, req LoginRequest) (*LoginReply, error) {
	var out LoginReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathLogin, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Sync 在 Login 后收取符合本次玩家属性的区服邮件；缺少属性可能成功但不匹配对应邮件。
func (c *Client) Sync(ctx context.Context, req SyncRequest) (*SyncReply, error) {
	var out SyncReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSync, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List 返回邮箱内 ID，并推进列表已看游标；后续 Get/Read/Claim/Remove 使用这些 ID。
func (c *Client) List(ctx context.Context, req ListRequest) (*ListReply, error) {
	var out ListReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathList, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get 按 List/通知中的邮箱内 ID 获取详情，不接受 Send 返回的正文引用。
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
