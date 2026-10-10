package officialweb

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubHandler struct{}

func (stubHandler) LookupRole(_ context.Context, _ *http.Request, req LookupRoleRequest) Reply {
	return Reply{Code: CodeSuccess, Message: "ok", Data: map[string]any{"cp_role_id": req.CpRoleID}}
}
func (stubHandler) ListProducts(context.Context, *http.Request, ListProductsRequest) Reply {
	return Reply{Code: CodeSuccess, Message: "ok", Data: map[string]any{"goods_list": []any{}}}
}
func (stubHandler) PreCheck(context.Context, *http.Request, PreCheckRequest) Reply {
	return Reply{Code: CodeSuccess, Message: "ok"}
}
func (stubHandler) CreateOrder(context.Context, *http.Request, AddOrderRequest) Reply {
	return Reply{Code: CodeSuccess, Message: "ok", Data: map[string]any{"cp_order_id": "cp-1"}}
}
func (stubHandler) NotifyPaid(context.Context, *http.Request, BuyPaymentRequest) Reply {
	return Reply{Code: CodeSuccess, Message: "ok"}
}

func TestInboundLookupAndBadSign(t *testing.T) {
	in, err := NewInbound("test-sign-secret", stubHandler{})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	in.Register(mux)

	body := map[string]any{"cp_role_id": "player-web-prod", "cp_server_id": "1"}
	body["sign"] = SignSortedQSMD5(body, "test-sign-secret")
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, PathLookupRole, bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("http %d", rec.Code)
	}
	var env map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if int(env["code"].(float64)) != 0 {
		t.Fatalf("env %+v", env)
	}

	bad, _ := json.Marshal(map[string]any{"cp_role_id": "player-web-prod", "sign": "bad"})
	req = httptest.NewRequest(http.MethodPost, PathLookupRole, bytes.NewReader(bad))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if int(env["code"].(float64)) != 401 {
		t.Fatalf("want 401 got %+v", env)
	}

	req = httptest.NewRequest(http.MethodGet, PathHealthz, nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	b, _ := io.ReadAll(rec.Body)
	if rec.Code != 200 || string(b) != "ok" {
		t.Fatalf("healthz %d %q", rec.Code, b)
	}
}

func TestNewInboundNilHandler(t *testing.T) {
	_, err := NewInbound("x", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

// 验签：雪花 JSON number 必须用 UseNumber；struct 为 string 时 number 进线 → 400。
func TestInboundSignAcceptsSnowflakeJSONNumberThenTypeError(t *testing.T) {
	in, err := NewInbound("test-sign-secret", stubHandler{})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	in.Register(mux)

	const snowflake int64 = 3747523271598286848
	body := map[string]any{"platform_role_id": snowflake}
	body["sign"] = SignSortedQSMD5(body, "test-sign-secret")
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, PathListProducts, bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var env map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if int(env["code"].(float64)) != 400 {
		t.Fatalf("want 400 type error got %+v", env)
	}
}

func TestInboundListProductsSnowflakeString(t *testing.T) {
	in, err := NewInbound("test-sign-secret", stubHandler{})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	in.Register(mux)

	body := map[string]any{"platform_role_id": "3747523271598286848"}
	body["sign"] = SignSortedQSMD5(body, "test-sign-secret")
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, PathListProducts, bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var env map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if int(env["code"].(float64)) != 0 {
		t.Fatalf("want code=0 got %+v", env)
	}
}

// V0：缺 product_id 仍进 Handler（stub 成功）。
func TestInboundCreateOrderMissingProductIDReachesHandler(t *testing.T) {
	in, err := NewInbound("test-sign-secret", stubHandler{})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	in.Register(mux)
	body := map[string]any{"platform_role_id": "1"}
	body["sign"] = SignSortedQSMD5(body, "test-sign-secret")
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, PathCreateOrder, bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	var env map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &env)
	if int(env["code"].(float64)) != 0 {
		t.Fatalf("want 0 (handler) got %+v", env)
	}
}
