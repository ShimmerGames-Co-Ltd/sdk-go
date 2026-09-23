package core

import (
	"errors"
	"fmt"
	"net/http"
)

// TransportError 表示网络/客户端层失败（未拿到完整 HTTP 响应）。
type TransportError struct {
	Err error
}

func (e *TransportError) Error() string {
	if e == nil || e.Err == nil {
		return "sdk-go: transport error"
	}
	return fmt.Sprintf("sdk-go: transport: %v", e.Err)
}

func (e *TransportError) Unwrap() error { return e.Err }

// HTTPError 表示拿到了 HTTP 响应但状态码不是 200。
// 业务失败是 HTTP 200 + body.code，不会落到本类型。
type HTTPError struct {
	StatusCode int
	Body       string
	Header     http.Header
}

func (e *HTTPError) Error() string {
	if e == nil {
		return "sdk-go: http error"
	}
	return fmt.Sprintf("sdk-go: http %d: %s", e.StatusCode, e.Body)
}

// APIError 表示 HTTP 200 且业务 body code != 0。
type APIError struct {
	Code    int
	Message string
	Body    string
}

func (e *APIError) Error() string {
	if e == nil {
		return "sdk-go: api error"
	}
	if e.Message != "" {
		return fmt.Sprintf("sdk-go: api code=%d message=%s", e.Code, e.Message)
	}
	return fmt.Sprintf("sdk-go: api code=%d", e.Code)
}

// ResponseBodyError 表示 HTTP 200 但响应不是合法 {code,data,message} 新 body。
type ResponseBodyError struct {
	Err  error
	Body string
}

func (e *ResponseBodyError) Error() string {
	if e == nil {
		return "sdk-go: response body error"
	}
	return fmt.Sprintf("sdk-go: response body: %v", e.Err)
}

func (e *ResponseBodyError) Unwrap() error { return e.Err }

// EnvelopeError 旧名别名。
type EnvelopeError = ResponseBodyError

// AsAPIError 提取业务 body 错误。
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	ok := errors.As(err, &apiErr)
	return apiErr, ok
}

// AsHTTPError 提取非 200 HTTP 错误。
func AsHTTPError(err error) (*HTTPError, bool) {
	var httpErr *HTTPError
	ok := errors.As(err, &httpErr)
	return httpErr, ok
}

// ConfigError 表示构造 Client 时参数不合法。
type ConfigError struct {
	Msg string
}

func (e *ConfigError) Error() string {
	if e == nil {
		return "sdk-go: config error"
	}
	return "sdk-go: " + e.Msg
}
