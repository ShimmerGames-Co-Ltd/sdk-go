package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const nonceAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nonceLength = 32

// HMACSHA256Signer 计算 Authorization 头所用的 HMAC-SHA256 签名。
// 旧 sdk-go 类型名误叫 SHA256WithRSASign，实际算法一直是 HMAC，本类型禁止再出现 RSA 字样。
type HMACSHA256Signer struct {
	OrganizationID string
	AppID          string
	Secret         string
}

// SignMessage 对「URI+timestamp+nonce+body」原文做 HMAC-SHA256 再 Base64。
func SignMessage(secret, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// AuthorizationHeader 按服务端 parseSignature 格式组装 Authorization。
func AuthorizationHeader(orgID, appID, nonce string, timestamp int64, signature string) string {
	return fmt.Sprintf(authorizationValue, orgID, appID, nonce, timestamp, signature)
}

// SignAuthorization 生成一条完整 Authorization 头。
func SignAuthorization(orgID, appID, secret, requestURI, body string, timestamp int64, nonce string) string {
	msg := fmt.Sprintf(SignatureFormat, requestURI, timestamp, nonce, body)
	sig := SignMessage(secret, msg)
	return AuthorizationHeader(orgID, appID, nonce, timestamp, sig)
}

// Authorization 用 signer 字段生成头。
func (s *HMACSHA256Signer) Authorization(requestURI, body string, timestamp int64, nonce string) string {
	if s == nil {
		return ""
	}
	return SignAuthorization(s.OrganizationID, s.AppID, s.Secret, requestURI, body, timestamp, nonce)
}

// UserVerifySign 计算 Auth ServerVerify 的 body.sign：HMAC-SHA1 + Base64 + URL Escape。
// 原文为升序 KV：appid=&token=&ts=。不要把本算法用在 Chat Authorization 上。
func UserVerifySign(appID, token, secret string, ts int64) string {
	kv := map[string]string{
		"appid": appID,
		"token": token,
		"ts":    fmt.Sprintf("%d", ts),
	}
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+kv[k])
	}
	content := strings.Join(parts, "&")
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(content))
	b64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return url.QueryEscape(b64)
}

func randomNonce() (string, error) {
	b := make([]byte, nonceLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	out := make([]byte, nonceLength)
	for i := 0; i < nonceLength; i++ {
		out[i] = nonceAlphabet[int(b[i])%len(nonceAlphabet)]
	}
	return string(out), nil
}
