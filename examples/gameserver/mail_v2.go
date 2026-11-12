package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/mail"
)

// mailV2Demo 的各分支都是完整调用；模板、语言和 Catalog 前置见 docs/mail-v2.md。
func mailV2Demo(ctx context.Context, c *core.Client, demo string) error {
	cli, err := mail.New(c)
	if err != nil {
		return err
	}
	serverID := envOr("SHIMO_SERVER_ID", "1")
	playerID := os.Getenv("SHIMO_PLAYER_ID")
	lang := envOr("SHIMO_DEFAULT_LANG", "English")
	// 每种发送场景使用文档中各自的 Serial；同一业务重试保留原值。
	serial := os.Getenv("SHIMO_MAIL_SERIAL")
	if demo != "mail-receive" {
		if err := requireEnv("SHIMO_MAIL_SERIAL"); err != nil {
			return err
		}
	}

	switch demo {
	case "mail-send-v2":
		if err := requireEnv("SHIMO_PLAYER_ID"); err != nil {
			return err
		}
		// 首次发送前构造请求；若要重试，复用 req，过期时刻不重新计算。
		req := mail.SendPlayerV2Request{
			From:            "system",
			Serial:          serial,
			ServerID:        serverID,
			PlayerID:        playerID,
			ContentMode:     mail.ContentModePlain,
			ContentFormat:   mail.ContentFormatPlainText,
			ContentCategory: "welcome",
			ExpireAt:        time.Now().Add(24 * time.Hour).Unix(),
			I18n: []mail.I18nBlock{
				{
					Lang:     lang,
					Title:    "Welcome",
					Content:  "Welcome to the game.",
					SubTitle: "A new adventure",
				},
			},
		}
		reply, err := cli.SendPlayerV2(ctx, req)
		if err != nil {
			return err
		}
		fmt.Printf("mail reference=%s unread=%d\n", reply.MailRefID, reply.Badge.UnreadMailCount)
	case "mail-group-v2":
		if err := requireEnv("SHIMO_PLAYER_ID", "SHIMO_SECOND_SERVER_ID", "SHIMO_SECOND_PLAYER_ID"); err != nil {
			return err
		}
		// 每个区服只列一次，Type=0 仅为本示例的项目奖励类型约定。
		reply, err := cli.SendGroupV2(ctx, mail.SendGroupV2Request{
			From:   "system",
			Serial: serial,
			RecipientTargets: []mail.RecipientTarget{
				{
					ServerID:  serverID,
					PlayerIDs: []string{playerID},
				},
				{
					ServerID:  os.Getenv("SHIMO_SECOND_SERVER_ID"),
					PlayerIDs: []string{os.Getenv("SHIMO_SECOND_PLAYER_ID")},
				},
			},
			I18n: []mail.I18nBlock{
				{
					Lang:    lang,
					Title:   "Reward",
					Content: "Your reward is ready.",
				},
			},
			Attachments: []mail.MailAttachment{
				{
					Type:  0,
					ID:    "gold",
					Count: 100,
					Param: "event-reward",
				},
			},
		})
		if err != nil {
			return err
		}
		for _, ref := range reply.Refs {
			fmt.Printf("server=%s player=%s reference=%s\n", ref.ServerID, ref.PlayerID, ref.MailRefID)
		}
	case "mail-server-v2":
		// App 全服：不与具体区服混用。这里显式允许发送后注册的玩家收取。
		reply, err := cli.SendServerV2(ctx, mail.SendServerV2Request{
			From:                   "system",
			Serial:                 serial,
			ServerIDs:              []string{"*"},
			IncludeLateRegistrants: true,
			I18n: []mail.I18nBlock{
				{
					Lang:    lang,
					Title:   "Announcement",
					Content: "Welcome, all players.",
				},
			},
		})
		if err != nil {
			return err
		}
		for _, ref := range reply.Refs {
			fmt.Printf("server=%s reference=%s\n", ref.ServerID, ref.MailRefID)
		}
	case "mail-audience-v2":
		// Catalog 的 register_at 须支持 range；Sync 上传真实注册时间。
		// 只构造一次筛选时间范围，重试时复用 req。
		req := mail.SendServerV2Request{
			From:      "system",
			Serial:    serial,
			ServerIDs: []string{serverID},
			Audience: &mail.Audience{
				Op: mail.AudienceOpAnd,
				Children: []mail.AudienceNode{
					{
						Key: "register_at",
						Op:  mail.AudienceOpRange,
						Value: mail.AudienceRange{
							Min: time.Now().Add(-7 * 24 * time.Hour).Unix(),
						},
					},
				},
			},
			I18n: []mail.I18nBlock{
				{
					Lang:    lang,
					Title:   "New player gift",
					Content: "A gift for players registered in the last seven days.",
				},
			},
		}
		reply, err := cli.SendServerV2(ctx, req)
		if err != nil {
			return err
		}
		for _, ref := range reply.Refs {
			fmt.Printf("server=%s reference=%s\n", ref.ServerID, ref.MailRefID)
		}
	case "mail-platform-template-v2":
		if err := requireEnv("SHIMO_TEMPLATE_CODE"); err != nil {
			return err
		}
		// 后台模板须启用且恰好使用这两个占位符；不传 I18n，格式使用固定模板版本。
		reply, err := cli.SendServerV2(ctx, mail.SendServerV2Request{
			From:                   "system",
			Serial:                 serial,
			ServerIDs:              []string{serverID},
			ContentMode:            mail.ContentModePlatformTemplate,
			TemplateID:             os.Getenv("SHIMO_TEMPLATE_CODE"),
			IncludeLateRegistrants: true,
			TemplateArgs: map[string]any{
				"plain.player_name": "Ada",
				"i18n.item_name": map[string]any{
					lang: "Legendary Sword",
				},
			},
		})
		if err != nil {
			return err
		}
		for _, ref := range reply.Refs {
			fmt.Printf("server=%s reference=%s\n", ref.ServerID, ref.MailRefID)
		}
	case "mail-client-template-v2":
		if err := requireEnv("SHIMO_PLAYER_ID", "SHIMO_CLIENT_TEMPLATE_ID"); err != nil {
			return err
		}
		// 客户端本地渲染；模板 ID 不能与后台已登记 code 冲突，参数不使用平台目录。
		reply, err := cli.SendPlayerV2(ctx, mail.SendPlayerV2Request{
			From:          "system",
			Serial:        serial,
			ServerID:      serverID,
			PlayerID:      playerID,
			ContentMode:   mail.ContentModeClientTemplate,
			ContentFormat: mail.ContentFormatRichText,
			TemplateID:    os.Getenv("SHIMO_CLIENT_TEMPLATE_ID"),
			TemplateArgs: map[string]any{
				"amount":  7,
				"item_id": "gold",
			},
		})
		if err != nil {
			return err
		}
		fmt.Printf("mail reference=%s\n", reply.MailRefID)
	case "mail-receive":
		if err := requireEnv("SHIMO_PLAYER_ID", "SHIMO_REGISTER_AT"); err != nil {
			return err
		}
		registeredAt, err := strconv.ParseInt(os.Getenv("SHIMO_REGISTER_AT"), 10, 64)
		if err != nil || registeredAt <= 0 {
			return &core.ConfigError{
				Msg: "SHIMO_REGISTER_AT 须为玩家真实注册时间的正整数 Unix 秒",
			}
		}
		// 只有 via_game_server 模式在这里读邮箱；client_direct 由客户端 SDK 走用户接口。
		_, err = cli.Login(ctx, mail.LoginRequest{
			ServerID: serverID,
			PlayerID: playerID,
			Lang:     lang,
		})
		if err != nil {
			return err
		}
		synced, err := cli.Sync(ctx, mail.SyncRequest{
			ServerID: serverID,
			PlayerID: playerID,
			PlayerAttrs: map[string]any{
				"register_at": registeredAt,
			},
		})
		if err != nil {
			return err
		}
		listed, err := cli.List(ctx, mail.ListRequest{
			ServerID: serverID,
			PlayerID: playerID,
			Limit:    20,
		})
		if err != nil {
			return err
		}
		fmt.Printf("new=%d total=%d\n", synced.NewMailCount, listed.Total)
		if len(listed.Headers) == 0 {
			return nil
		}
		// 使用 List 的邮箱内 ID；发送返回的 MailRefID 不能用于 Get。
		ids := make([]int64, 0, len(listed.Headers))
		for _, header := range listed.Headers {
			ids = append(ids, int64(header.ID))
		}
		detail, err := cli.Get(ctx, mail.GetRequest{
			ServerID: serverID,
			PlayerID: playerID,
			MailIDs:  ids,
		})
		if err != nil {
			return err
		}
		for _, item := range detail.Mails {
			var attachments []mail.MailAttachment
			if item.Attachment != "" {
				if err := json.Unmarshal([]byte(item.Attachment), &attachments); err != nil {
					return err
				}
			}
			// 按项目已知模板解码成业务结构，不需要对 map 值做类型断言或数值转换。
			if item.TemplateID == envOr("SHIMO_CLIENT_TEMPLATE_ID", "local_reward") && item.TemplateArgs != "" {
				var args struct {
					Amount int    `json:"amount"`
					ItemID string `json:"item_id"`
				}
				if err := json.Unmarshal([]byte(item.TemplateArgs), &args); err != nil {
					return err
				}
				fmt.Printf("reward amount=%d item=%s\n", args.Amount, args.ItemID)
			}
			fmt.Printf("id=%d format=%q attachments=%d template=%q\n",
				item.ID, item.ContentFormat, len(attachments), item.TemplateID)
		}
	default:
		return &core.ConfigError{
			Msg: "未知邮件示例: " + demo,
		}
	}
	return nil
}

func requireEnv(names ...string) error {
	for _, name := range names {
		if os.Getenv(name) == "" {
			return &core.ConfigError{
				Msg: "缺少必填环境变量 " + name,
			}
		}
	}
	return nil
}
