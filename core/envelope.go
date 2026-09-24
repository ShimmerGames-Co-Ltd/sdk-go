package core

import (
	"encoding/json"
	"fmt"
)

// ResponseBody 是新 body：HTTP 200 + {code,data,message}（common.Response 形态）。
type ResponseBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// DecodeResponseBody 解析新 body。HTTP 状态必须已确认为 200。
func DecodeResponseBody(body []byte) (*ResponseBody, error) {
	var env ResponseBody
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, &ResponseBodyError{Err: err, Body: string(body)}
	}
	return &env, nil
}

// UnmarshalData 把 data 解进 out；data 为空或 null 时跳过。
func (e *ResponseBody) UnmarshalData(out any) error {
	if e == nil || out == nil {
		return nil
	}
	if len(e.Data) == 0 || string(e.Data) == "null" {
		return nil
	}
	if err := json.Unmarshal(e.Data, out); err != nil {
		return fmt.Errorf("sdk-go: decode data: %w", err)
	}
	return nil
}
