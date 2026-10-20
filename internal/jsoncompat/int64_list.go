package jsoncompat

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Int64List 在 SDK 内部兼容 protojson 的字符串整数与 JSON 数字。
// 底层使用 []int64，转换回业务字段时无需逐元素复制。
type Int64List []int64

func (s *Int64List) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var out Int64List
	if raw != nil {
		out = make(Int64List, len(raw))
	}
	for i, elem := range raw {
		value := string(elem)
		if len(elem) > 0 && elem[0] == '"' {
			if err := json.Unmarshal(elem, &value); err != nil {
				return fmt.Errorf("int64 list[%d]: %w", i, err)
			}
		}
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("int64 list[%d]: %w", i, err)
		}
		out[i] = n
	}
	*s = out
	return nil
}
