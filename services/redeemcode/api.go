package redeemcode

import (
	"context"
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
