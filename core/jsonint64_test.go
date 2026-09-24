package core

import (
	"encoding/json"
	"testing"
)

func TestJSONInt64NumberAndString(t *testing.T) {
	t.Parallel()
	body := struct {
		N JSONInt64 `json:"n"`
	}{N: 42}
	raw, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `{"n":42}` {
		t.Fatalf("encode=%s", raw)
	}
	for _, in := range []string{`{"n":42}`, `{"n":"42"}`} {
		var got struct {
			N JSONInt64 `json:"n"`
		}
		if err := json.Unmarshal([]byte(in), &got); err != nil {
			t.Fatal(err)
		}
		if got.N != 42 {
			t.Fatalf("%s => %d", in, got.N)
		}
	}
}
