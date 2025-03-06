package credential

import (
	"fmt"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core/auth"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core/consts"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/utils"
	"time"
)

type PlatformCredential struct {
	Signer auth.Signer
}

func (c *PlatformCredential) CreateAuthorizationHeader(uri, signData string) (string, error) {

	if c.Signer == nil {
		return "", fmt.Errorf("not init signer")
	}

	timestamp := time.Now().Unix()

	nonce := utils.GenerateNonce()

	message := fmt.Sprintf(consts.SignatureFormat, uri, timestamp, nonce, signData)

	signature := c.Signer.Sign(message)

	authorization := fmt.Sprintf(
		consts.HeaderAuthorizationValue,
		signature.OrganizationId, signature.AppId, nonce, timestamp, signature.Signature)

	return authorization, nil
}

func (c *PlatformCredential) TokenSign(token string, timestamp int64) (string, error) {

	if c.Signer == nil {
		return "", fmt.Errorf("not init signer")
	}

	signature, err := c.Signer.UserSign(token, timestamp)

	return signature, err
}
