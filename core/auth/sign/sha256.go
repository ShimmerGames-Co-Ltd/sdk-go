package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"gitlab.shimmer.lab/platform/sdk-go/core/auth"
	"gitlab.shimmer.lab/platform/sdk-go/utils/sign"
	"net/url"
)

type SHA256WithRSASign struct {
	OrganizationId string
	AppId          string
	Secret         string
}

func (s *SHA256WithRSASign) Sign(message string) *auth.SignatureInfo {

	mac := hmac.New(sha256.New, []byte(s.Secret))
	mac.Write([]byte(message))

	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return &auth.SignatureInfo{
		OrganizationId: s.OrganizationId,
		AppId:          s.AppId,
		Secret:         s.Secret,
		Signature:      signature,
	}
}

func (s *SHA256WithRSASign) UserSign(token string, timestamp int64) (string, error) {

	signature, err := sign.Generate(
		map[string]any{
			"appid": s.AppId,
			"ts":    timestamp,
			"token": token,
		},
		sign.WithSecret(s.Secret),
		sign.WithUseKv(true),
		sign.WithJoiner("&"),
		sign.WithSortKeys(sign.ASC),
		sign.WithSignAlg(sign.HmacSha1WithBase64))

	signature = url.QueryEscape(signature)

	return signature, err
}
