package leaderboard

import (
	"context"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/services"
)

const (
	ApiPath = "/leaderboard/"
)

type Service services.Service

// Register 注册排行榜
func (svc *Service) Register(ctx context.Context, req RegisterRequest) (resp *RegisterReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/register"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(RegisterReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// Delete 删除排行榜
func (svc *Service) Delete(ctx context.Context, req DeleteRequest) (resp *DeleteReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/delete"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(DeleteReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// DeleteApp 清空 APP 全部排行榜
func (svc *Service) DeleteApp(ctx context.Context, req DeleteAppRequest) (resp *DeleteAppReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/app/delete"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(DeleteAppReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// AllOfAPP 获取 APP 全部排行榜
func (svc *Service) AllOfAPP(ctx context.Context, req AllOfAPPRequest) (resp *AllOfAPPReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/app/all"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(AllOfAPPReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// Ask 查询排行榜状态
func (svc *Service) Ask(ctx context.Context, req AskRequest) (resp *AskReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/ask"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(AskReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// SetScore 玩家设定绝对分
func (svc *Service) SetScore(ctx context.Context, req SetScoreRequest) (resp *SetScoreReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/score/set"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(SetScoreReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// IncrScore 玩家增量分
func (svc *Service) IncrScore(ctx context.Context, req IncrScoreRequest) (resp *IncrScoreReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/score/incr"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(IncrScoreReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// GetScore 获取玩家当前分
func (svc *Service) GetScore(ctx context.Context, req GetScoreRequest) (resp *GetScoreReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/score/get"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(GetScoreReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// List 分页拉取排行榜
func (svc *Service) List(ctx context.Context, req ListRequest) (resp *ListReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/list"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(ListReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// BanAdd 封禁玩家
func (svc *Service) BanAdd(ctx context.Context, req BanAddRequest) (resp *BanAddReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/ban/add"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(BanAddReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// BanRemove 解除封禁
func (svc *Service) BanRemove(ctx context.Context, req BanRemoveRequest) (resp *BanRemoveReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/ban/remove"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(BanRemoveReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// BanClean 清空 APP 全部封禁
func (svc *Service) BanClean(ctx context.Context, req BanCleanRequest) (resp *BanCleanReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/ban/clean"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(BanCleanReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// BanGet 获取封禁详情
func (svc *Service) BanGet(ctx context.Context, req BanGetRequest) (resp *BanGetReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/ban/get"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(BanGetReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// BanList 分页拉取封禁列表
func (svc *Service) BanList(ctx context.Context, req BanListRequest) (resp *BanListReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/ban/list"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(BanListReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}

// MemberInfo 查询玩家当日上报次数
func (svc *Service) MemberInfo(ctx context.Context, req MemberInfoRequest) (resp *MemberInfoReply, result *core.APIResult, err error) {
	path := ApiPath + "leaderboard/member/info"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(MemberInfoReply)
	err = core.UnmarshalResponse(result.Response, resp)
	return
}
