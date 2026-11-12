package mail

import (
	"context"
	"net/http"
)

const (
	pathSendPlayerV2 = "/mail/v2/mail/server/send/player"
	pathSendGroupV2  = "/mail/v2/mail/server/send/group"
	pathSendServerV2 = "/mail/v2/mail/server/send/server"
)

// SendPlayerV2 通过独立 V2 入口发送个人邮；参数校验由服务端执行，失败不降级。
// 网络结果未知时使用相同 Serial、收件人及内容重试；同 Serial 不更新已发邮件。
func (c *Client) SendPlayerV2(ctx context.Context, req SendPlayerV2Request) (*SendPlayerV2Reply, error) {
	var out SendPlayerV2Reply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendPlayerV2, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendGroupV2 发送多区服群组邮，失败不降级。任一收件人失败即返回 nil 和错误，已投递部分不回滚。
// 修复失败原因后使用原 Serial、收件人及内容重试，只补未投递部分；新业务使用新 Serial。
func (c *Client) SendGroupV2(ctx context.Context, req SendGroupV2Request) (*SendGroupV2Reply, error) {
	var out SendGroupV2Reply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendGroupV2, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendServerV2 写入各区服邮箱，失败不降级；任一区服失败即返回 nil 和错误，已写入部分不回滚。
// 结果未知或部分失败时使用原 Serial、区服范围及内容重试；玩家通过 Sync 收取，发送成功不等于已进箱。
func (c *Client) SendServerV2(ctx context.Context, req SendServerV2Request) (*SendServerV2Reply, error) {
	var out SendServerV2Reply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendServerV2, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
