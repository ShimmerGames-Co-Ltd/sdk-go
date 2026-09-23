package core

import (
	"net/http"
	"strings"
	"time"
)

// Option 配置 Client。
type Option func(*settings) error

type settings struct {
	URL            string
	AppID          string
	OrganizationID string
	AuthSecret     string
	UserToken      string
	HTTPClient     *http.Client
	ServerVersion  string
	PathPrefix     string
	SignRequestURI func(method, requestURI string) string
	Now            func() time.Time
	Nonce          func() (string, error)
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

// WithAuthSecret 启用游戏服 SERVER_SIGNATURE（Authorization HMAC）。
func WithAuthSecret(secret string) Option {
	return func(s *settings) error {
		s.AuthSecret = secret
		return nil
	}
}

// WithUserToken 启用玩家 X-User-Token。
func WithUserToken(token string) Option {
	return func(s *settings) error {
		s.UserToken = token
		return nil
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(s *settings) error {
		s.HTTPClient = c
		return nil
	}
}

func WithServerVersion(v string) Option {
	return func(s *settings) error {
		s.ServerVersion = strings.TrimSpace(v)
		return nil
	}
}

// WithPathPrefix 仅影响签算 URI（例如 Kong 未 strip 时的 /chat 前缀），不会改实际请求 URL。
func WithPathPrefix(prefix string) Option {
	return func(s *settings) error {
		if prefix == "" {
			s.PathPrefix = ""
			return nil
		}
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}
		s.PathPrefix = strings.TrimRight(prefix, "/")
		return nil
	}
}

// WithSignRequestURI 完全覆盖签算 URI。SDK 不会设置 X-Original-URI。
func WithSignRequestURI(fn func(method, requestURI string) string) Option {
	return func(s *settings) error {
		s.SignRequestURI = fn
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
