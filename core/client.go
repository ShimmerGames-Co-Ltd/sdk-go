package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// AuthMode 构造期选定，禁止每请求切换。
type AuthMode int

const (
	AuthUnspecified AuthMode = iota
	AuthServerSignature
	AuthUserToken
)

// Client 是签名 + 新信封 HTTP 客户端。游戏服不要 import 主仓内部包，只依赖本 module。
type Client struct {
	httpClient     *http.Client
	url            string
	appID          string
	organizationID string
	authSecret     string
	userToken      string
	serverVersion  string
	pathPrefix     string
	signRequestURI func(method, requestURI string) string
	now            func() time.Time
	nonce          func() (string, error)
	mode           AuthMode
}

// NewClient 创建客户端。WithURL 与 WithAppID 必填；鉴权二选一：WithAuthSecret 或 WithUserToken。
func NewClient(opts ...Option) (*Client, error) {
	s := &settings{
		ServerVersion: DefaultServerVersion,
		Now:           func() time.Time { return time.Now() },
		Nonce:         randomNonce,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	if s.URL == "" {
		return nil, &ConfigError{Msg: "WithURL 必填，禁止写死生产域名"}
	}
	if s.AppID == "" {
		return nil, &ConfigError{Msg: "WithAppID 必填"}
	}
	if s.AuthSecret != "" && s.UserToken != "" {
		return nil, &ConfigError{Msg: "WithAuthSecret 与 WithUserToken 不能同时设置"}
	}
	if s.AuthSecret == "" && s.UserToken == "" {
		return nil, &ConfigError{Msg: "必须设置 WithAuthSecret 或 WithUserToken"}
	}
	mode := AuthUserToken
	if s.AuthSecret != "" {
		mode = AuthServerSignature
	}
	hc := s.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: DefaultHTTPTimeout}
	}
	if s.Now == nil {
		s.Now = time.Now
	}
	if s.Nonce == nil {
		s.Nonce = randomNonce
	}
	if s.ServerVersion == "" {
		s.ServerVersion = DefaultServerVersion
	}
	return &Client{
		httpClient:     hc,
		url:            s.URL,
		appID:          s.AppID,
		organizationID: s.OrganizationID,
		authSecret:     s.AuthSecret,
		userToken:      s.UserToken,
		serverVersion:  s.ServerVersion,
		pathPrefix:     s.PathPrefix,
		signRequestURI: s.SignRequestURI,
		now:            s.Now,
		nonce:          s.Nonce,
		mode:           mode,
	}, nil
}

func (c *Client) AppID() string      { return c.appID }
func (c *Client) AuthMode() AuthMode { return c.mode }
func (c *Client) AuthSecret() string { return c.authSecret }

// DoJSON 发送 JSON 请求并解码新信封。POST 的 body 使用 json.Encoder（带尾随换行）参与签算。
func (c *Client) DoJSON(ctx context.Context, method, path string, body any, out any) error {
	if ctx == nil {
		ctx = context.Background()
	}
	var bodyStr string
	var bodyReader io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return fmt.Errorf("sdk-go: encode body: %w", err)
		}
		bodyStr = buf.String()
		bodyReader = bytes.NewReader(buf.Bytes())
	}

	fullURL, requestURI, err := c.buildURL(path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return &TransportError{Err: err}
	}
	if err := c.applyHeaders(req, method, requestURI, bodyStr); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &TransportError{Err: err}
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return &TransportError{Err: err}
	}

	if resp.StatusCode != http.StatusOK {
		return &HTTPError{StatusCode: resp.StatusCode, Body: string(raw), Header: resp.Header.Clone()}
	}
	env, err := DecodeResponseBody(raw)
	if err != nil {
		return err
	}
	if env.Code != 0 {
		return &APIError{Code: env.Code, Message: env.Message, Body: string(raw)}
	}
	return env.UnmarshalData(out)
}

func (c *Client) applyHeaders(req *http.Request, method, requestURI, bodyStr string) error {
	req.Header.Set(headerAccept, "*/*")
	req.Header.Set(headerAppID, c.appID)
	req.Header.Set(headerServerVersion, c.serverVersion)
	req.Header.Set(headerUserAgent, fmt.Sprintf(headerUserAgentFormat, runtime.GOOS, runtime.Version()))
	if req.Body != nil && method != http.MethodGet && method != http.MethodHead {
		req.Header.Set(headerContentType, headerContentJSON)
	}
	switch c.mode {
	case AuthServerSignature:
		if c.organizationID != "" {
			req.Header.Set(headerOrganizationID, c.organizationID)
		}
		ts := c.now().Unix()
		nonce, err := c.nonce()
		if err != nil {
			return &TransportError{Err: err}
		}
		signURI := c.signURI(method, requestURI)
		signer := HMACSHA256Signer{OrganizationID: c.organizationID, AppID: c.appID, Secret: c.authSecret}
		req.Header.Set(headerAuthorization, signer.Authorization(signURI, bodyStr, ts, nonce))
	case AuthUserToken:
		req.Header.Set(headerUserToken, c.userToken)
	}
	return nil
}

func (c *Client) signURI(method, requestURI string) string {
	if c.signRequestURI != nil {
		return c.signRequestURI(method, requestURI)
	}
	if c.pathPrefix == "" {
		return requestURI
	}
	path := requestURI
	query := ""
	if i := strings.Index(requestURI, "?"); i >= 0 {
		path, query = requestURI[:i], requestURI[i:]
	}
	return c.pathPrefix + path + query
}

func (c *Client) buildURL(path string) (fullURL, requestURI string, err error) {
	if path == "" {
		return "", "", &ConfigError{Msg: "path 不能为空"}
	}
	rel, err := url.Parse(path)
	if err != nil {
		return "", "", &ConfigError{Msg: "非法 path: " + path}
	}
	if _, err := url.Parse(c.url); err != nil {
		return "", "", &ConfigError{Msg: "非法 WithURL"}
	}
	joined := strings.TrimRight(c.url, "/") + "/" + strings.TrimLeft(rel.Path, "/")
	if rel.RawQuery != "" {
		joined += "?" + rel.RawQuery
	}
	u, err := url.Parse(joined)
	if err != nil {
		return "", "", &ConfigError{Msg: "拼接 URL 失败"}
	}
	return u.String(), u.RequestURI(), nil
}
