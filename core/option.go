package core

import (
	"net/http"
	"strings"
	"time"
)

// Option 配置 Client。
type Option func(*settings) error

type settings struct {
	URL              string
	AppID            string
	OrganizationID   string
	ServerSignSecret string
	HTTPClient       *http.Client
	Now              func() time.Time
	Nonce            func() (string, error)
}

func WithURL(u string) Option {
	return func(s *settings) error {
		s.URL = strings.TrimRight(strings.TrimSpace(u), "/")
		return nil
	}
}

func WithAppID(id string) Option {
	return func(s *settings) error {
		s.AppID = strings.TrimSpace(id)
		return nil
	}
}

func WithOrganizationID(id string) Option {
	return func(s *settings) error {
		s.OrganizationID = strings.TrimSpace(id)
		return nil
	}
}

// WithServerSignSecret 设置服务器 SDK 签名秘钥。
// 对应 Hub app.auth_secret，不是客户端 SDK 的 sdk_secret。
// 同一密钥用于 Authorization（HMAC-SHA256）和 Auth Verify 请求体（HMAC-SHA1）。
func WithServerSignSecret(secret string) Option {
	return func(s *settings) error {
		s.ServerSignSecret = secret
		return nil
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(s *settings) error {
		s.HTTPClient = c
		return nil
	}
}

func withNow(fn func() time.Time) Option {
	return func(s *settings) error {
		s.Now = fn
		return nil
	}
}

func withNonce(fn func() (string, error)) Option {
	return func(s *settings) error {
		s.Nonce = fn
		return nil
	}
}
