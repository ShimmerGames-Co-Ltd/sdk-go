package auth

type Credential interface {
	CreateAuthorizationHeader(uri, signData string) (string, error)
	TokenSign(token string, timestamp int64) (string, error)
}
