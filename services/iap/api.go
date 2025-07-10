package iap

import (
	"context"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/services"
)

const (
	ApiPath = "/iap/"
)

type PaymentService services.Service

// VerifyPayment
// 向 支付服务 请求订单状态, 返回订单对应商品与状态
func (svc *PaymentService) VerifyPayment(ctx context.Context, req VerifyPaymentRequest) (resp *VerifyPaymentResponse, result *core.APIResult, err error) {

	path := ApiPath + "server/verify"

	result, err = svc.Client.RequestPost(ctx, path, nil, req)
	if err != nil {
		return
	}

	resp = new(VerifyPaymentResponse)
	err = core.UnmarshalResponse(result.Response, resp)
	if err != nil {
		return nil, result, err
	}

	if resp.State == OrderExternalStateVerifySuccess {
		resp.Success = true
	}

	return resp, result, nil
}
