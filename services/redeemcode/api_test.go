package redeemcode

import (
	"context"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core/option"
	"testing"
)

func TestRedeemCode(t *testing.T) {
	sdkOpts := []core.ClientOption{
		// 服务侧验证信息
		option.WithAuthCipher("10000000000", "shim01", "H34RWX1Gdn4giFa7vAoffCEzWJ2RC1F0"),
	}
	sdkOpts = append(sdkOpts, option.WithUrl("http://platform-api.shimmer.lab"))
	ctx := context.Background()
	// sdk client
	client, err := core.NewClient(ctx, sdkOpts...)
	if err != nil {
		t.Logf("create sdk client failed. %s", err)
		return
	}
	redeemService := RedeemCodeService{Client: client}
	rsp1, _, err := redeemService.RedeemCode(ctx, RedeemCodeRequest{Code: "NSOIDROY47BID"})
	if err != nil {
		t.Errorf("redeem code failed. %s", err)
	}
	t.Log(rsp1)
}
