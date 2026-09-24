package leaderboard

import (
	"context"
	"net/http"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	pathRegister          = "/leaderboard/v1/server/register"
	pathDelete            = "/leaderboard/v1/server/delete"
	pathDeleteApp         = "/leaderboard/v1/server/app/delete"
	pathAllOfAPP          = "/leaderboard/v1/server/app/all"
	pathAsk               = "/leaderboard/v1/server/ask"
	pathSetScore          = "/leaderboard/v1/server/score/set"
	pathIncrScore         = "/leaderboard/v1/server/score/incr"
	pathGetScore          = "/leaderboard/v1/server/score/get"
	pathList              = "/leaderboard/v1/server/list"
	pathBanAdd            = "/leaderboard/v1/server/ban/add"
	pathBanRemove         = "/leaderboard/v1/server/ban/remove"
	pathBanClean          = "/leaderboard/v1/server/ban/clean"
	pathBanGet            = "/leaderboard/v1/server/ban/get"
	pathBanList           = "/leaderboard/v1/server/ban/list"
	pathMemberInfo        = "/leaderboard/v1/server/member/info"
	pathHistoryMemberInfo = "/leaderboard/v1/server/history/member/info"
)

// Client 封装 Leaderboard 游戏服主路径 /leaderboard/v1/server/*。
type Client struct {
	core *core.Client
}

func New(c *core.Client) (*Client, error) {
	if c == nil {
		return nil, &core.ConfigError{Msg: "core.Client 不能为空"}
	}
	return &Client{core: c}, nil
}

func post[T any](c *Client, ctx context.Context, path string, body any) (*T, error) {
	var out T
	if err := c.core.DoJSON(ctx, http.MethodPost, path, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Register(ctx context.Context, req RegisterRequest) (*RegisterReply, error) {
	return post[RegisterReply](c, ctx, pathRegister, req)
}

func (c *Client) Delete(ctx context.Context, req DeleteRequest) error {
	return c.core.DoJSON(ctx, http.MethodPost, pathDelete, req, nil)
}

func (c *Client) DeleteApp(ctx context.Context, req DeleteAppRequest) error {
	return c.core.DoJSON(ctx, http.MethodPost, pathDeleteApp, req, nil)
}

func (c *Client) AllOfAPP(ctx context.Context) (*AllOfAPPReply, error) {
	return post[AllOfAPPReply](c, ctx, pathAllOfAPP, map[string]any{})
}

func (c *Client) Ask(ctx context.Context, req AskRequest) (*AskReply, error) {
	return post[AskReply](c, ctx, pathAsk, req)
}

func (c *Client) SetScore(ctx context.Context, req SetScoreRequest) (*SetScoreReply, error) {
	return post[SetScoreReply](c, ctx, pathSetScore, req)
}

func (c *Client) IncrScore(ctx context.Context, req IncrScoreRequest) (*IncrScoreReply, error) {
	return post[IncrScoreReply](c, ctx, pathIncrScore, req)
}

func (c *Client) GetScore(ctx context.Context, req GetScoreRequest) (*GetScoreReply, error) {
	return post[GetScoreReply](c, ctx, pathGetScore, req)
}

func (c *Client) List(ctx context.Context, req ListRequest) (*ListReply, error) {
	return post[ListReply](c, ctx, pathList, req)
}

func (c *Client) BanAdd(ctx context.Context, req BanAddRequest) error {
	return c.core.DoJSON(ctx, http.MethodPost, pathBanAdd, req, nil)
}

func (c *Client) BanRemove(ctx context.Context, req BanRemoveRequest) error {
	return c.core.DoJSON(ctx, http.MethodPost, pathBanRemove, req, nil)
}

func (c *Client) BanClean(ctx context.Context) error {
	return c.core.DoJSON(ctx, http.MethodPost, pathBanClean, map[string]any{}, nil)
}

func (c *Client) BanGet(ctx context.Context, req BanGetRequest) (*BanGetReply, error) {
	return post[BanGetReply](c, ctx, pathBanGet, req)
}

func (c *Client) BanList(ctx context.Context, req BanListRequest) (*BanListReply, error) {
	return post[BanListReply](c, ctx, pathBanList, req)
}

func (c *Client) MemberInfo(ctx context.Context, req MemberInfoRequest) (*MemberInfoReply, error) {
	return post[MemberInfoReply](c, ctx, pathMemberInfo, req)
}

func (c *Client) HistoryMemberInfo(ctx context.Context, req HistoryMemberInfoRequest) (*HistoryMemberInfoReply, error) {
	return post[HistoryMemberInfoReply](c, ctx, pathHistoryMemberInfo, req)
}
