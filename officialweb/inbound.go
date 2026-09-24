package officialweb

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	PathLookupRole   = "/payment/lookup_role"
	PathListProducts = "/payment/get_goods_list"
	PathPreCheck     = "/payment/pre_check"
	PathCreateOrder  = "/payment/add_order"
	PathNotifyPaid   = "/payment/buypayment"
	PathHealthz      = "/healthz"

	// CodeSuccess 官充入站成功码（HTTP 恒 200，业务看 body.code）。
	CodeSuccess = 0
)

// Mux 宿主路由（如 net/http.ServeMux、Kratos HTTP Server）。
type Mux interface {
	Handle(path string, h http.Handler)
}

// Reply 入站业务结果，由宿主状态机填写。
type Reply struct {
	Code    int
	Message string
	Data    any
}

// PaymentHandler 官充五条业务回调。验签与信封由 Inbound 处理。
type PaymentHandler interface {
	LookupRole(ctx context.Context, r *http.Request, body map[string]any) Reply
	ListProducts(ctx context.Context, r *http.Request, body map[string]any) Reply
	PreCheck(ctx context.Context, r *http.Request, body map[string]any) Reply
	CreateOrder(ctx context.Context, r *http.Request, body map[string]any) Reply
	NotifyPaid(ctx context.Context, r *http.Request, body map[string]any) Reply
}

// Inbound 把官充 /payment/* 注册到宿主 mux。
type Inbound struct {
	Secret  string
	Handler PaymentHandler
}

func NewInbound(secret string, h PaymentHandler) (*Inbound, error) {
	if h == nil {
		return nil, &core.ConfigError{Msg: "officialweb PaymentHandler 不能为空"}
	}
	return &Inbound{Secret: secret, Handler: h}, nil
}

// RegisterPayment 注册五条官充入站。
func (in *Inbound) RegisterPayment(mux Mux) {
	if in == nil || mux == nil || in.Handler == nil {
		return
	}
	h := in.Handler
	mux.Handle(PathLookupRole, in.wrap(h.LookupRole))
	mux.Handle(PathListProducts, in.wrap(h.ListProducts))
	mux.Handle(PathPreCheck, in.wrap(h.PreCheck))
	mux.Handle(PathCreateOrder, in.wrap(h.CreateOrder))
	mux.Handle(PathNotifyPaid, in.wrap(h.NotifyPaid))
}

// Register 注册五条入站 + /healthz（GET/POST 正文 ok，无验签）。
func (in *Inbound) Register(mux Mux) {
	in.RegisterPayment(mux)
	if mux != nil {
		mux.Handle(PathHealthz, http.HandlerFunc(Healthz))
	}
}

// Healthz 探活：GET/POST 返回明文 ok。
func Healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// WriteResponseBody 写 HTTP 200 + {code,message,data}。
func WriteResponseBody(w http.ResponseWriter, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if data == nil {
		data = map[string]any{}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func (in *Inbound) wrap(next func(ctx context.Context, r *http.Request, body map[string]any) Reply) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			WriteResponseBody(w, 400, "method not allowed", nil)
			return
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			WriteResponseBody(w, 500, "read body", nil)
			return
		}
		body := map[string]any{}
		if len(raw) > 0 {
			// UseNumber：保留 JSON number 十进制文本，避免雪花 int64 解成 float64 后验签失败。
			dec := json.NewDecoder(bytes.NewReader(raw))
			dec.UseNumber()
			if err := dec.Decode(&body); err != nil {
				WriteResponseBody(w, 400, "invalid json", nil)
				return
			}
		}
		gotSign, _ := body["sign"].(string)
		delete(body, "sign")
		want := SignSortedQSMD5(body, in.Secret)
		if strings.TrimSpace(gotSign) == "" || gotSign != want {
			WriteResponseBody(w, 401, "invalid sign", nil)
			return
		}
		rep := next(r.Context(), r, body)
		WriteResponseBody(w, rep.Code, rep.Message, rep.Data)
	})
}
