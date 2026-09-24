package core

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// JSONInt64 同时接受 JSON number 与十进制 string。
// 新 body 的 protojson 把 int64 编成 string；请求侧仍按 number 发出。
type JSONInt64 int64

func (n *JSONInt64) MarshalJSON() ([]byte, error) {
	if n == nil {
		return []byte("0"), nil
	}
	return []byte(strconv.FormatInt(int64(*n), 10)), nil
}

func (n *JSONInt64) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || string(b) == "null" {
		*n = 0
		return nil
	}
	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*n = 0
			return nil
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*n = JSONInt64(v)
		return nil
	}
	v, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*n = JSONInt64(v)
	return nil
}
