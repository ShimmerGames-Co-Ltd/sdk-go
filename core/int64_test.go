package core

import (
	"encoding/json"
	"testing"
)

// 向量：超过 JS Number.MAX_SAFE_INTEGER 的 int64，禁止经 float64。
const (
	unsafeSeq = int64(9007199254740993)
	snowflake = int64(1234567890123456789)
)

func TestInt64ResponseBodyKeepsExact(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"code":0,"data":{"seq":9007199254740993,"id":1234567890123456789}}`)
	env, err := DecodeResponseBody(raw)
	if err != nil {
		t.Fatal(err)
	}
	var out struct {
		Seq int64 `json:"seq"`
		ID  int64 `json:"id"`
	}
	if err := env.UnmarshalData(&out); err != nil {
		t.Fatal(err)
	}
	if out.Seq != unsafeSeq || out.ID != snowflake {
		t.Fatalf("got seq=%d id=%d", out.Seq, out.ID)
	}
}

func TestInt64MustNotUseFloat64Map(t *testing.T) {
	t.Parallel()
	var m map[string]any
	if err := json.Unmarshal([]byte(`{"seq":9007199254740993}`), &m); err != nil {
		t.Fatal(err)
	}
	f, ok := m["seq"].(float64)
	if !ok {
		t.Fatalf("encoding/json 解到 any 是 %T，本测试只用来对照", m["seq"])
	}
	if int64(f) == unsafeSeq {
		t.Fatal("float64 碰巧能表示该值，换更大向量")
	}
	// 对照：同一 JSON 解到 int64 必须精确（SDK 公共 API 只用结构体）。
	var typed struct {
		Seq int64 `json:"seq"`
	}
	if err := json.Unmarshal([]byte(`{"seq":9007199254740993}`), &typed); err != nil {
		t.Fatal(err)
	}
	if typed.Seq != unsafeSeq {
		t.Fatalf("int64 字段丢精度: %d", typed.Seq)
	}
}
