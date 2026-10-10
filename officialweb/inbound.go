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

// PaymentHandler 官充五条业务回调。入站：验签 → Unmarshal → 本接口（Go struct）。
// 不做必填非空校验；缺字段为零值，由宿主业务处理。add_order 可含可选 ProductNameI18n（JSON 字符串）。
type PaymentHandler interface {
	LookupRole(ctx context.Context, r *http.Request, req LookupRoleRequest) Reply
	ListProducts(ctx context.Context, r *http.Request, req ListProductsRequest) Reply
	PreCheck(ctx context.Context, r *http.Request, req PreCheckRequest) Reply
	CreateOrder(ctx context.Context, r *http.Request, req AddOrderRequest) Reply
	NotifyPaid(ctx context.Context, r *http.Request, req BuyPaymentRequest) Reply
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
	mux.Handle(PathLookupRole, in.withVerified(in.handleLookupRole))
	mux.Handle(PathListProducts, in.withVerified(in.handleListProducts))
	mux.Handle(PathPreCheck, in.withVerified(in.handlePreCheck))
	mux.Handle(PathCreateOrder, in.withVerified(in.handleCreateOrder))
	mux.Handle(PathNotifyPaid, in.withVerified(in.handleNotifyPaid))
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

func (in *Inbound) withVerified(next func(http.ResponseWriter, *http.Request, []byte)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := in.verifyPOST(w, r)
		if !ok {
			return
		}
		next(w, r, raw)
	})
}

// verifyPOST：统一验签。map+UseNumber 仅服务 SignSortedQSMD5；通过后仍返回 raw 供 Unmarshal。
func (in *Inbound) verifyPOST(w http.ResponseWriter, r *http.Request) (raw []byte, ok bool) {
	if r.Method != http.MethodPost {
		WriteResponseBody(w, 400, "method not allowed", nil)
		return nil, false
	}
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		WriteResponseBody(w, 500, "read body", nil)
		return nil, false
	}
	body := map[string]any{}
	if len(raw) > 0 {
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.UseNumber() // 雪花 ID 勿解成 float64，否则验签失败
		if err := dec.Decode(&body); err != nil {
			WriteResponseBody(w, 400, "invalid json", nil)
			return nil, false
		}
	}
	gotSign, _ := body["sign"].(string)
	delete(body, "sign")
	want := SignSortedQSMD5(body, in.Secret)
	if strings.TrimSpace(gotSign) == "" || gotSign != want {
		WriteResponseBody(w, 401, "invalid sign", nil)
		return nil, false
	}
	return raw, true
}

func (in *Inbound) handleLookupRole(w http.ResponseWriter, r *http.Request, raw []byte) {
	dispatch(w, r, raw, in.Handler.LookupRole)
}

func (in *Inbound) handleListProducts(w http.ResponseWriter, r *http.Request, raw []byte) {
	dispatch(w, r, raw, in.Handler.ListProducts)
}

func (in *Inbound) handlePreCheck(w http.ResponseWriter, r *http.Request, raw []byte) {
	dispatch(w, r, raw, in.Handler.PreCheck)
}

func (in *Inbound) handleCreateOrder(w http.ResponseWriter, r *http.Request, raw []byte) {
	dispatch(w, r, raw, in.Handler.CreateOrder)
}

func (in *Inbound) handleNotifyPaid(w http.ResponseWriter, r *http.Request, raw []byte) {
	dispatch(w, r, raw, in.Handler.NotifyPaid)
}

func dispatch[T any](w http.ResponseWriter, r *http.Request, raw []byte, next func(context.Context, *http.Request, T) Reply) {
	var req T
	if err := unmarshalBody(raw, &req); err != nil {
		WriteResponseBody(w, 400, err.Error(), nil)
		return
	}
	rep := next(r.Context(), r, req)
	WriteResponseBody(w, rep.Code, rep.Message, rep.Data)
}
