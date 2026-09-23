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

func (stubHandler) LookupRole(_ context.Context, _ *http.Request, body map[string]any) Reply {
	return Reply{Code: CodeSuccess, Message: "ok", Data: map[string]any{"project_role_id": body["project_role_id"]}}
}
func (stubHandler) ListProducts(context.Context, *http.Request, map[string]any) Reply {
	return Reply{Code: CodeSuccess, Message: "ok", Data: map[string]any{"goods_list": []any{}}}
}
func (stubHandler) PreCheck(context.Context, *http.Request, map[string]any) Reply {
	return Reply{Code: CodeSuccess, Message: "ok"}
}
func (stubHandler) CreateOrder(context.Context, *http.Request, map[string]any) Reply {
	return Reply{Code: CodeSuccess, Message: "ok", Data: map[string]any{"cp_oid": "cp-1"}}
}
func (stubHandler) NotifyPaid(context.Context, *http.Request, map[string]any) Reply {
	return Reply{Code: CodeSuccess, Message: "ok"}
}

func TestInboundLookupAndBadSign(t *testing.T) {
	in, err := NewInbound("test-sign-secret", stubHandler{})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	in.Register(mux)

	body := map[string]any{"project_role_id": "player-web-prod"}
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

	bad, _ := json.Marshal(map[string]any{"project_role_id": "player-web-prod", "sign": "bad"})
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

// 雪花 role_id 超过 float64 安全整数（2^53）。JSON number 必须用 UseNumber 验签，不能解成 float64。
func TestInboundListProductsSnowflakeInt64JSONNumber(t *testing.T) {
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
	if rec.Code != http.StatusOK {
		t.Fatalf("http %d", rec.Code)
	}
	var env map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if int(env["code"].(float64)) != 0 {
		t.Fatalf("want code=0 got %+v", env)
	}
}
