package mail

import (
	"context"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/services"
)

const (
	ApiPath = "/mail/"
)

type Service services.Service

// SendToUser 发送邮件到用户
func (svc *Service) SendToUser(ctx context.Context, req SendPlayerMailRequest) (resp *SendPlayerMailReply, result *core.APIResult, err error) {

	path := ApiPath + "send/player"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(SendPlayerMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// SendToServer 发送邮件到全服务
func (svc *Service) SendToServer(ctx context.Context, req SendServerMailRequest) (resp *SendServerMailReply, result *core.APIResult, err error) {

	path := ApiPath + "send/server"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(SendServerMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// SendToGroup 群发邮件
func (svc *Service) SendToGroup(ctx context.Context, req SendGroupMailRequest) (resp *SendGroupMailReply, result *core.APIResult, err error) {

	path := ApiPath + "send/group"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(SendGroupMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// UserLogin 用户登录邮件
func (svc *Service) UserLogin(ctx context.Context, req LoginRequest) (resp *LoginReply, result *core.APIResult, err error) {

	path := ApiPath + "player/login"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(LoginReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// UserSync 用户同步邮件
func (svc *Service) UserSync(ctx context.Context, req SyncMailRequest) (resp *SyncMailReply, result *core.APIResult, err error) {

	path := ApiPath + "player/sync"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(SyncMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// UserList 用户拉取邮件列表
func (svc *Service) UserList(ctx context.Context, req ListMailRequest) (resp *ListMailReply, result *core.APIResult, err error) {

	path := ApiPath + "player/list"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(ListMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// UserGet 用户获取邮件详情
func (svc *Service) UserGet(ctx context.Context, req GetMailRequest) (resp *GetMailReply, result *core.APIResult, err error) {

	path := ApiPath + "player/get"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(GetMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// UserSetRead 用户设置已读
func (svc *Service) UserSetRead(ctx context.Context, req ReadMailRequest) (resp *ReadMailReply, result *core.APIResult, err error) {

	path := ApiPath + "player/read"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(ReadMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// UserClaimAttachment 用户领取附件
func (svc *Service) UserClaimAttachment(ctx context.Context, req ClaimMailRequest) (resp *ClaimMailReply, result *core.APIResult, err error) {

	path := ApiPath + "player/claim"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(ClaimMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}

// UserRemove 用户删除邮件
func (svc *Service) UserRemove(ctx context.Context, req RemoveMailRequest) (resp *RemoveMailReply, result *core.APIResult, err error) {

	path := ApiPath + "player/remove"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(RemoveMailReply)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}
