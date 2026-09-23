package officialweb

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

// SignSortedQSMD5 实现官充游戏入站 S2：剔除 sign → key 字典序 → querystring → 拼 secret → MD5 hex。
// 数组按重复 key 序列化（param=1&param=2），与中台 gameintegration 一致。
func SignSortedQSMD5(params map[string]any, secret string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var b strings.Builder
	first := true
	for _, k := range keys {
		for _, v := range stringifyParamValues(params[k]) {
			if !first {
				b.WriteByte('&')
			}
			first = false
			b.WriteString(url.QueryEscape(k))
			b.WriteByte('=')
			b.WriteString(url.QueryEscape(v))
		}
	}
	sum := md5.Sum([]byte(b.String() + secret))
	return hex.EncodeToString(sum[:])
}

func stringifyParamValues(v any) []string {
	if v == nil {
		return []string{""}
	}
	switch x := v.(type) {
	case string:
		return []string{x}
	case []string:
		out := make([]string, len(x))
		copy(out, x)
		return out
	case []int64:
		out := make([]string, len(x))
		for i, n := range x {
			out[i] = strconv.FormatInt(n, 10)
		}
		return out
	case []int:
		out := make([]string, len(x))
		for i, n := range x {
			out[i] = strconv.Itoa(n)
		}
		return out
	case []any:
		out := make([]string, len(x))
		for i, item := range x {
			out[i] = fmt.Sprint(item)
		}
		return out
	case json.Number:
		return []string{x.String()}
	default:
		return []string{fmt.Sprint(x)}
	}
}
