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

// Client 是游戏服签名 HTTP 客户端。只支持服务器签名，不要 import 主仓内部包。
type Client struct {
	httpClient       *http.Client
	url              string
	appID            string
	organizationID   string
	serverSignSecret string
	serverVersion    string
	now              func() time.Time
	nonce            func() (string, error)
}

// NewClient 创建客户端。WithURL、WithAppID、WithServerSignSecret 必填。
// X-Server-Version 固定为 v3.5.0。
func NewClient(opts ...Option) (*Client, error) {
	s := &settings{
		Now:   func() time.Time { return time.Now() },
		Nonce: randomNonce,
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
	if s.ServerSignSecret == "" {
		return nil, &ConfigError{Msg: "WithServerSignSecret 必填"}
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
	return &Client{
		httpClient:       hc,
		url:              s.URL,
		appID:            s.AppID,
		organizationID:   s.OrganizationID,
		serverSignSecret: s.ServerSignSecret,
		serverVersion:    DefaultServerVersion,
		now:              s.Now,
		nonce:            s.Nonce,
	}, nil
}

func (c *Client) AppID() string { return c.appID }

// ServerSignSecret 返回服务器 SDK 签名秘钥（Hub app.auth_secret）。
func (c *Client) ServerSignSecret() string { return c.serverSignSecret }

// DoJSON 发送 JSON 请求并解码新 body。POST 的 body 使用 json.Encoder（带尾随换行）参与签算。
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
	if c.organizationID != "" {
		req.Header.Set(headerOrganizationID, c.organizationID)
	}
	ts := c.now().Unix()
	nonce, err := c.nonce()
	if err != nil {
		return &TransportError{Err: err}
	}
	signer := HMACSHA256Signer{OrganizationID: c.organizationID, AppID: c.appID, Secret: c.serverSignSecret}
	req.Header.Set(headerAuthorization, signer.Authorization(requestURI, bodyStr, ts, nonce))
	return nil
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
