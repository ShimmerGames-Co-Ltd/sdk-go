package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type SortType int

const (
	ASC SortType = iota
	DESC
)

type Alg int

const (
	Empty Alg = iota
	Md5
	HmacSha1WithBase64
	HmacSha256
)

var AlgFuncMap map[Alg]func(content, secret string, opt *option) string

func init() {
	AlgFuncMap = map[Alg]func(content, secret string, opt *option) string{
		Empty:              signEmpty,
		HmacSha1WithBase64: signHmacSha1WithBase64,
		HmacSha256:         signHmacSha256,
		Md5:                signMd5,
	}
}

func signHmacSha1WithBase64(content, secret string, opt *option) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func signHmacSha256(content, secret string, opt *option) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(content))
	return string(mac.Sum(nil))
}

func signMd5(content, secret string, opt *option) string {
	bts := md5.Sum([]byte(strings.Join([]string{content, secret}, opt.Joiner)))
	return hex.EncodeToString(bts[:])
}

func signEmpty(content, _ string, _ *option) string {
	return content
}

type option struct {
	Joiner        string
	SkipKeys      []string
	SkipNull      bool
	SkipNullValue bool
	UseKv         bool
	SortKeys      SortType
	SignAlg       Alg
	Secret        string
	NullStr       string
	UseNullStr    bool
}

type SignOption func(sg *option)

func WithJoiner(joiner string) SignOption {
	return func(sg *option) {
		sg.Joiner = joiner
	}
}

func WithSkipKeys(skipKeys ...string) SignOption {
	return func(sg *option) {
		sg.SkipKeys = append(sg.SkipKeys, skipKeys...)
	}
}

func WithSkipNull(skipNull bool) SignOption {
	return func(sg *option) {
		sg.SkipNull = skipNull
	}
}

func WithSkipNullValue(skipNullValue bool) SignOption {
	return func(sg *option) {
		sg.SkipNullValue = skipNullValue
	}
}

func WithUseKv(useKv bool) SignOption {
	return func(sg *option) {
		sg.UseKv = useKv
	}
}

func WithSortKeys(sortKeys SortType) SignOption {
	return func(sg *option) {
		sg.SortKeys = sortKeys
	}
}

func WithSignAlg(signAlg Alg) SignOption {
	return func(sg *option) {
		sg.SignAlg = signAlg
	}
}

func WithNullStr(nullstr string) SignOption {
	return func(sg *option) {
		sg.NullStr = nullstr
		sg.UseNullStr = true
	}
}

func WithSecret(secret string) SignOption {
	return func(sg *option) {
		sg.Secret = secret
	}
}

func newOptions(opts ...SignOption) *option {
	sg := &option{
		SkipKeys: []string{},
	}
	for _, opt := range opts {
		opt(sg)
	}
	return sg
}

// isNull 判断值是否为空
func isNull(value interface{}) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		return v == ""
	case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, uint8, float32, float64:
		return v == 0
	case bool:
		return !v
	}
	return false
}

func ReverseStrings(arr []string) {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		arr[i], arr[n-1-i] = arr[n-1-i], arr[i]
	}
}

func Contains(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

// Generate 生成签名
func Generate(params map[string]interface{}, opts ...SignOption) (string, error) {
	sg := newOptions(opts...)
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}

	if sg.SortKeys == ASC {
		sort.Strings(keys)
	} else if sg.SortKeys == DESC {
		ReverseStrings(keys)
	}

	// 为保证兼容性, 取消 go 1.21 下slices包的使用
	// if sg.SortKeys == ASC {
	//	 slices.Sort(keys)
	// } else if sg.SortKeys == DESC {
	//	 slices.Reverse(keys)
	// }

	ss := []string{}
	for _, k := range keys {
		if Contains(sg.SkipKeys, k) {
			continue
		}
		value := params[k]
		var v string
		if isNull(value) {
			if sg.SkipNull {
				continue
			}
			if sg.UseNullStr {
				v = sg.NullStr
			} else {
				v = fmt.Sprintf("%v", params[k])
			}
		} else {
			v = fmt.Sprintf("%v", params[k])
		}
		if sg.UseKv {
			ss = append(ss, fmt.Sprintf("%s=%s", k, v))
		} else {
			ss = append(ss, v)
		}
	}
	content := strings.Join(ss, sg.Joiner)
	if f, exist := AlgFuncMap[sg.SignAlg]; !exist {
		return "", fmt.Errorf("sign alg %v not exist", sg.SignAlg)
	} else {
		return f(content, sg.Secret, sg), nil
	}
}
