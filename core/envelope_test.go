package core

import (
	"errors"
	"testing"
)

func TestDecodeResponseBodyBusinessCode(t *testing.T) {
	t.Parallel()
	env, err := DecodeResponseBody([]byte(`{"code":14003,"message":"账号已被禁言"}`))
	if err != nil {
		t.Fatal(err)
	}
	if env.Code != 14003 {
		t.Fatalf("code=%d", env.Code)
	}
	apiErr := &APIError{Code: env.Code, Message: env.Message, Body: "x"}
	if !errors.As(apiErr, new(*APIError)) {
		t.Fatal("APIError 应可 errors.As")
	}
}

func TestDecodeResponseBodyInvalidJSON(t *testing.T) {
	t.Parallel()
	_, err := DecodeResponseBody([]byte(`not-json`))
	var envErr *ResponseBodyError
	if !errors.As(err, &envErr) {
		t.Fatalf("want ResponseBodyError got %T %v", err, err)
	}
}

func TestUnmarshalDataSkipNull(t *testing.T) {
	t.Parallel()
	env := &ResponseBody{Code: 0, Data: []byte("null")}
	var out map[string]any
	if err := env.UnmarshalData(&out); err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Fatalf("null data should skip, got %#v", out)
	}
}
