package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/auth"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/iap"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/leaderboard"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/mail"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/redeemcode"
)

func main() {
	baseURL := os.Getenv("SHIMO_URL")
	appID := os.Getenv("SHIMO_APP_ID")
	orgID := os.Getenv("SHIMO_ORG_ID")
	secret := os.Getenv("SHIMO_AUTH_SECRET")
	demo := os.Getenv("SHIMO_DEMO")
	if baseURL == "" || appID == "" || secret == "" || demo == "" {
		fmt.Fprintln(os.Stderr, "用法: SHIMO_URL SHIMO_APP_ID SHIMO_AUTH_SECRET SHIMO_DEMO=mail|leaderboard|auth|iap|redeem [SHIMO_ORG_ID] [SHIMO_ACCOUNT_ID|SHIMO_PLAYER_ID|SHIMO_TOKEN|SHIMO_ORDER_ID|SHIMO_CODE]")
		os.Exit(2)
	}

	c, err := core.NewClient(
		core.WithURL(baseURL),
		core.WithAppID(appID),
		core.WithOrganizationID(orgID),
		core.WithAuthSecret(secret),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	switch demo {
	case "mail":
		cli, err := mail.New(c)
		must(err)
		rep, err := cli.Login(ctx, mail.LoginRequest{
			PlayerID: os.Getenv("SHIMO_PLAYER_ID"),
			ServerID: envOr("SHIMO_SERVER_ID", "1"),
			Lang:     envOr("SHIMO_LANG", "zh"),
		})
		must(err)
		fmt.Printf("mail unread=%d\n", rep.Badge.UnreadMailCount)
	case "leaderboard":
		cli, err := leaderboard.New(c)
		must(err)
		rep, err := cli.Ask(ctx, leaderboard.AskRequest{ID: envOr("SHIMO_LB_ID", "demo")})
		must(err)
		fmt.Printf("leaderboard id=%s name=%s\n", rep.ID, rep.Name)
	case "auth":
		cli, err := auth.New(c)
		must(err)
		rep, err := cli.Verify(ctx, os.Getenv("SHIMO_TOKEN"))
		must(err)
		fmt.Printf("auth user_id=%d role_id=%d\n", rep.UserID, rep.RoleID)
	case "iap":
		cli, err := iap.New(c)
		must(err)
		rep, err := cli.VerifyOrder(ctx, iap.VerifyOrderRequest{OrderID: os.Getenv("SHIMO_ORDER_ID")})
		must(err)
		fmt.Printf("iap order=%s state=%d success=%v\n", rep.OrderID, rep.State, iap.VerifySuccess(rep.State))
	case "redeem":
		cli, err := redeemcode.New(c)
		must(err)
		rep, err := cli.Check(ctx, redeemcode.CheckRequest{
			Code:      os.Getenv("SHIMO_CODE"),
			UserIDStr: os.Getenv("SHIMO_PLAYER_ID"),
		})
		must(err)
		fmt.Printf("redeem result=%d\n", rep.Result)
	default:
		fmt.Fprintf(os.Stderr, "未知 SHIMO_DEMO=%s\n", demo)
		os.Exit(2)
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
