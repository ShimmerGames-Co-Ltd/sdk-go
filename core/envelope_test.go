package core

import (
	"errors"
	"testing"
)

func TestDecodeEnvelopeBusinessCode(t *testing.T) {
	t.Parallel()
	env, err := DecodeEnvelope([]byte(`{"code":14003,"message":"账号已被禁言"}`))
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

func TestDecodeEnvelopeInvalidJSON(t *testing.T) {
	t.Parallel()
	_, err := DecodeEnvelope([]byte(`not-json`))
	var envErr *EnvelopeError
	if !errors.As(err, &envErr) {
		t.Fatalf("want EnvelopeError got %T %v", err, err)
	}
}

func TestUnmarshalDataSkipNull(t *testing.T) {
	t.Parallel()
	env := &Envelope{Code: 0, Data: []byte("null")}
	var out map[string]any
	if err := env.UnmarshalData(&out); err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Fatalf("null data should skip, got %#v", out)
	}
}
