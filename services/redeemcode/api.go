package redeemcode

import (
	"context"
	"fmt"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/services"
)

const (
	ApiPath = "/redeemcode/"
)

type RedeemCodeService services.Service

// RedeemCode
// 兑换码
func (svc *RedeemCodeService) RedeemCode(ctx context.Context, req RedeemCodeRequest) (resp *RedeemCodeResponse, result *core.APIResult, err error) {
	path := ApiPath + "user/code/redeem"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(RedeemCodeResponse)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}
	return resp, result, nil
}

// CheckCode
// 代玩家只读检查兑换码是否可兑换, 不会产生兑换记录
func (svc *RedeemCodeService) CheckCode(ctx context.Context, req CheckCodeRequest) (resp *CheckCodeResponse, result *core.APIResult, err error) {
	if req.Code == "" {
		err = fmt.Errorf("redeem code is empty")
		return
	}
	// 服务端要求玩家身份非空
	if req.RoleId == 0 {
		err = fmt.Errorf("role id is empty")
		return
	}
	path := ApiPath + "v1/server/codes:check"
	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}
	resp = new(CheckCodeResponse)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}
	return resp, result, nil
}
