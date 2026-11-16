package mail_test

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/mail"
)

func ExampleNewHTTPNotifyHandler() {
	receiver, err := mail.NewHTTPNotifyHandler(
		os.Getenv("SHIMO_APP_ID"),
		os.Getenv("SHIMO_AUTH_SECRET"),
		func(ctx context.Context, r *http.Request, notification mail.HTTPNotifyRequest) error {
			// 示例只输出事件；项目在此完成处理或入队，再返回 nil。
			// SDK 已检查通知种类和单封约束，可以直接取 Headers[0]。
			if notification.PlayerMail != nil {
				message := notification.PlayerMail
				fmt.Printf("server=%s player=%s mailbox_id=%d\n",
					message.ServerID, message.PlayerID, message.Headers[0].ID)
			} else {
				message := notification.ServerMail
				fmt.Printf("server=%s reference=%s\n", message.ServerID, message.Headers[0].MailRefID)
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/mail/new_mail_notify", receiver)
}
