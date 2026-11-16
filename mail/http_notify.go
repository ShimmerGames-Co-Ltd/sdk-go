package mail

import (
	"crypto/hmac"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/v3/core"
)

const (
	HeaderNotifyAppID     = "X-App-ID"
	HeaderNotifySignature = "X-Mail-Signature"
	// HTTPNotifyMaxBodyBytes 限制单次通知正文，与当前Mail通知接收契约一致。
	HTTPNotifyMaxBodyBytes = 1 << 20
)

func (e *HTTPNotifyError) Error() string {
	return "sdk-go: mail HTTP notification: " + e.Message
}

type httpNotifyHandler struct {
	appID      string
	authSecret string
	onNotify   HTTPNotifyHandlerFunc
}

// NewHTTPNotifyHandler 可直接挂到net/http或包装给现有Web框架，路由由项目选择。
// authSecret必须为该App的auth_secret，不是sdk_secret；本函数不负责重试、Sync或去重。
func NewHTTPNotifyHandler(appID, authSecret string, onNotify HTTPNotifyHandlerFunc) (http.Handler, error) {
	if err := validateHTTPNotifyConfig(appID, authSecret); err != nil {
		return nil, err
	}
	if onNotify == nil {
		return nil, &core.ConfigError{
			Msg: "mail HTTP notification handler不能为空",
		}
	}
	return &httpNotifyHandler{
		appID:      strings.TrimSpace(appID),
		authSecret: authSecret,
		onNotify:   onNotify,
	}, nil
}

func (h *httpNotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	notification, err := ParseHTTPNotify(r, h.appID, h.authSecret)
	if err != nil {
		WriteHTTPNotifyError(w, err)
		return
	}
	if err := h.onNotify(r.Context(), r, notification); err != nil {
		WriteHTTPNotifyError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ParseHTTPNotify 供已有路由或动态配置接入：消费有界请求体，先验原文签名再解码。
// 校验X-App-ID、HMAC-SHA256/Base64、二选一通知及最小信封字段；业务处理仍由项目负责。
func ParseHTTPNotify(r *http.Request, appID, authSecret string) (HTTPNotifyRequest, error) {
	if r.Method != http.MethodPost {
		return HTTPNotifyRequest{}, &HTTPNotifyError{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "POST required",
		}
	}
	if err := validateHTTPNotifyConfig(appID, authSecret); err != nil {
		return HTTPNotifyRequest{}, err
	}
	if strings.TrimSpace(r.Header.Get(HeaderNotifyAppID)) != strings.TrimSpace(appID) {
		return HTTPNotifyRequest{}, &HTTPNotifyError{
			StatusCode: http.StatusUnauthorized,
			Message:    "app mismatch",
		}
	}
	signature := r.Header.Get(HeaderNotifySignature)
	if signature == "" {
		return HTTPNotifyRequest{}, &HTTPNotifyError{
			StatusCode: http.StatusUnauthorized,
			Message:    "notification signature required",
		}
	}
	raw, err := io.ReadAll(io.LimitReader(r.Body, HTTPNotifyMaxBodyBytes+1))
	if err != nil {
		return HTTPNotifyRequest{}, &HTTPNotifyError{
			StatusCode: http.StatusBadRequest,
			Message:    "cannot read notification body",
		}
	}
	if len(raw) > HTTPNotifyMaxBodyBytes {
		return HTTPNotifyRequest{}, &HTTPNotifyError{
			StatusCode: http.StatusRequestEntityTooLarge,
			Message:    "notification body exceeds 1 MiB",
		}
	}
	expected := core.SignMessage(authSecret, string(raw))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return HTTPNotifyRequest{}, &HTTPNotifyError{
			StatusCode: http.StatusUnauthorized,
			Message:    "invalid notification signature",
		}
	}
	var notification NewMailNotify
	if err := json.Unmarshal(raw, &notification); err != nil {
		return HTTPNotifyRequest{}, &HTTPNotifyError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid notification JSON",
		}
	}
	if err := validateHTTPNotification(notification); err != nil {
		return HTTPNotifyRequest{}, err
	}
	return HTTPNotifyRequest{
		NewMailNotify: notification,
		AppID:         strings.TrimSpace(appID),
		RawBody:       raw,
	}, nil
}

// WriteHTTPNotifyError 将协议错误写成非2xx；其它错误固定503，不向外泄漏业务错误内容。
func WriteHTTPNotifyError(w http.ResponseWriter, err error) {
	status := http.StatusServiceUnavailable
	message := "notification not accepted"
	var protocolError *HTTPNotifyError
	if errors.As(err, &protocolError) && protocolError.StatusCode >= 400 && protocolError.StatusCode <= 599 {
		status = protocolError.StatusCode
		message = protocolError.Message
	}
	if status == http.StatusMethodNotAllowed {
		w.Header().Set("Allow", http.MethodPost)
	}
	http.Error(w, message, status)
}

func validateHTTPNotifyConfig(appID, authSecret string) error {
	if strings.TrimSpace(appID) == "" || strings.TrimSpace(authSecret) == "" {
		return &core.ConfigError{
			Msg: "mail HTTP notification需要appID和auth_secret",
		}
	}
	return nil
}

func validateHTTPNotification(notification NewMailNotify) error {
	if (notification.PlayerMail == nil) == (notification.ServerMail == nil) {
		return &HTTPNotifyError{
			StatusCode: http.StatusBadRequest,
			Message:    "exactly one notification kind required",
		}
	}
	var serverID string
	var headerCount int
	if notification.PlayerMail != nil {
		serverID = notification.PlayerMail.ServerID
		headerCount = len(notification.PlayerMail.Headers)
		if strings.TrimSpace(notification.PlayerMail.PlayerID) == "" {
			return &HTTPNotifyError{
				StatusCode: http.StatusBadRequest,
				Message:    "player_id required",
			}
		}
	} else {
		serverID = notification.ServerMail.ServerID
		headerCount = len(notification.ServerMail.Headers)
	}
	serverID = strings.TrimSpace(serverID)
	if serverID == "" || serverID == "*" {
		return &HTTPNotifyError{
			StatusCode: http.StatusBadRequest,
			Message:    "concrete server_id required",
		}
	}
	if headerCount != 1 {
		return &HTTPNotifyError{
			StatusCode: http.StatusBadRequest,
			Message:    "exactly one mail header required",
		}
	}
	var mailID int64
	if notification.PlayerMail != nil {
		mailID = int64(notification.PlayerMail.Headers[0].ID)
	} else {
		mailID = int64(notification.ServerMail.Headers[0].ID)
	}
	if mailID <= 0 {
		return &HTTPNotifyError{
			StatusCode: http.StatusBadRequest,
			Message:    "positive mail header id required",
		}
	}
	return nil
}
