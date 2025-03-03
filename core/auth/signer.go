package auth

type SignatureInfo struct {
	OrganizationId string
	AppId          string
	Secret         string
	Signature      string
}

// Signer 数字签名生成器
type Signer interface {
	Sign(message string) *SignatureInfo
	UserSign(token string, timestamp int64) (string, error)
}
