package mail

import (
	"context"
	"net/http"
)

const (
	pathSendPlayer = "/mail/v1/mail/server/send/player"
	pathSendGroup  = "/mail/v1/mail/server/send/group"
	pathSendServer = "/mail/v1/mail/server/send/server"
)

func (c *Client) SendPlayer(ctx context.Context, req SendPlayerRequest) (*SendPlayerReply, error) {
	var out SendPlayerReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendPlayer, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendGroup(ctx context.Context, req SendGroupRequest) (*SendGroupReply, error) {
	var out SendGroupReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendGroup, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SendServer(ctx context.Context, req SendServerRequest) (*SendServerReply, error) {
	var out SendServerReply
	if err := c.core.DoJSON(ctx, http.MethodPost, pathSendServer, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
