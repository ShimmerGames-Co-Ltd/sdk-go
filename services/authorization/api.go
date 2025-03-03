package authorization

import (
	"context"
	"fmt"
	"gitlab.shimmer.lab/platform/sdk-go/core"
	"gitlab.shimmer.lab/platform/sdk-go/services"
	"time"
)

const (
	ApiPath = "/auth/"
)

type AuthUserService services.Service

// UserAuthorize
// 根据 账号服务 下发的 `Token`, 进行用户合法性验证, 并返回用户的 `OpenId` 与其他相关数据
func (svc *AuthUserService) UserAuthorize(ctx context.Context, req UserAuthorizeRequest) (resp *UserAuthorizeResponse, result *core.APIResult, err error) {

	path := ApiPath + "server/verify_user"

	now := time.Now().Unix()

	signature, err := svc.Client.TokenSign(req.Token, now)
	if err != nil {
		err = fmt.Errorf("generate user token sign failed. `%s`", err)
		return
	}

	req.Timestamp = now
	req.Sign = signature

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(UserAuthorizeResponse)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	return resp, result, nil
}
