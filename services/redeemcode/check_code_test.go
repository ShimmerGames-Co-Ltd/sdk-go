package redeemcode

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/ShimmerGames-Co-Ltd/sdk-go/core"
	"github.com/ShimmerGames-Co-Ltd/sdk-go/core/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testOrganizationId = "10000000000"
	testAppId          = "shim_test"
	testSecret         = "test-secret"
)

// 以下响应体为 shimo 修复后 (services/redeemcode/docs/修复方案-旧版body兑换码字段名兼容.md)
// 在 X-Server-Version: v1.0.0 下对 CheckCodeReply 的输出: 旧版扁平 body, 字段名与旧 define.Code 一致 (camelCase)
const (
	checkCodeSuccessBody = `{"code":{"ID":1234567890123456789,"activeTime":0,"appId":"shim_test","channel":"","code":"ABC","createTime":1700000000,` +
		`"effectTime":0,"expireTime":0,"groupId":-7,"maxUseCount":-1,"reward":"gold*10","status":1,"type":1,"useCount":3},` +
		`"record":{"appId":"shim_test","channel":"","code":"ABC","codeId":9,"groupId":-7,"redeemTime":1700000001,"rewarded":true,"userId":"u1"},"result":2}`
	checkCodeNotExistBody = `{"code":null,"record":null,"result":1}`
	checkCodeErrorBody    = `{"code":400,"details":[],"message":"user_id 不能为空"}`
)

var authorizationPattern = regexp.MustCompile(`(\w+)="([^"]*)"`)

type capturedRequest struct {
	method string
	uri    string
	header http.Header
	body   string
}

func newCheckCodeTestService(t *testing.T, status int, respBody string) (*RedeemCodeService, *capturedRequest) {
	t.Helper()

	captured := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		captured.method = r.Method
		captured.uri = r.RequestURI
		captured.header = r.Header.Clone()
		captured.body = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(respBody))
	}))
	t.Cleanup(srv.Close)

	client, err := core.NewClient(context.Background(),
		option.WithUrl(srv.URL),
		option.WithAuthCipher(testOrganizationId, testAppId, testSecret))
	require.NoError(t, err)

	return &RedeemCodeService{Client: client}, captured
}

func parseAuthorization(value string) map[string]string {
	fields := map[string]string{}
	for _, m := range authorizationPattern.FindAllStringSubmatch(value, -1) {
		fields[m[1]] = m[2]
	}
	return fields
}

func TestCheckCode_SendsSignedRequestToNewPath(t *testing.T) {
	svc, captured := newCheckCodeTestService(t, http.StatusOK, checkCodeSuccessBody)

	_, _, err := svc.CheckCode(context.Background(), CheckCodeRequest{Code: "ABC", RoleId: 9})
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, captured.method)
	assert.Equal(t, "/redeemcode/v1/server/codes:check", captured.uri)
	// shimo 从 X-App-ID 读取应用ID (Header 大小写不敏感), 且需旧版 body
	assert.Equal(t, testAppId, captured.header.Get("X-App-ID"))
	assert.Equal(t, "v1.0.0", captured.header.Get("X-Server-Version"))

	var body map[string]any
	require.NoError(t, json.Unmarshal([]byte(captured.body), &body))
	assert.Equal(t, map[string]any{"code": "ABC", "role_id": float64(9)}, body)

	// 与 shimo 服务端验签一致: HMAC-SHA256(uri + timestamp + nonce + body) 的 base64
	auth := parseAuthorization(captured.header.Get("Authorization"))
	assert.Equal(t, testOrganizationId, auth["organization_id"])
	assert.Equal(t, testAppId, auth["appid"])
	mac := hmac.New(sha256.New, []byte(testSecret))
	_, _ = fmt.Fprintf(mac, "%s%s%s%s", captured.uri, auth["timestamp"], auth["nonce"], captured.body)
	assert.Equal(t, base64.StdEncoding.EncodeToString(mac.Sum(nil)), auth["signature"])
}

func TestCheckCode_DecodesReply(t *testing.T) {
	svc, _ := newCheckCodeTestService(t, http.StatusOK, checkCodeSuccessBody)

	resp, result, err := svc.CheckCode(context.Background(), CheckCodeRequest{Code: "ABC", RoleId: 1})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, CodeResultReceivedBySelf, resp.Result)
	require.NotNil(t, resp.Code)
	assert.Equal(t, Code{
		ID: 1234567890123456789, Code: "ABC", Reward: "gold*10", GroupId: -7, AppId: testAppId,
		Type: CodeTypeDisposable, CreateTime: 1700000000, MaxUseCount: -1, UseCount: 3, Status: CodeStatusActive,
	}, *resp.Code)
	require.NotNil(t, resp.Record)
	assert.Equal(t, RedeemCodeRecord{
		UserId: "u1", CodeId: 9, AppId: testAppId, Code: "ABC", Rewarded: true, GroupId: -7, RedeemTime: 1700000001,
	}, *resp.Record)
}

func TestCheckCode_DecodesNotExistWithNullCode(t *testing.T) {
	svc, _ := newCheckCodeTestService(t, http.StatusOK, checkCodeNotExistBody)

	resp, _, err := svc.CheckCode(context.Background(), CheckCodeRequest{Code: "NOPE", RoleId: 9})

	require.NoError(t, err)
	assert.Equal(t, CodeResultNotExist, resp.Result)
	assert.Nil(t, resp.Code)
	assert.Nil(t, resp.Record)
}

func TestCheckCode_ReturnsPlatErrorOnErrorStatus(t *testing.T) {
	svc, _ := newCheckCodeTestService(t, http.StatusBadRequest, checkCodeErrorBody)

	resp, result, err := svc.CheckCode(context.Background(), CheckCodeRequest{Code: "ABC", RoleId: 1})

	assert.Nil(t, resp)
	assert.NotNil(t, result)
	var platErr *core.PlatError
	require.True(t, errors.As(err, &platErr))
	assert.Equal(t, http.StatusBadRequest, platErr.StatusCode)
	assert.Equal(t, 400, platErr.Code)
	assert.Equal(t, "user_id 不能为空", platErr.Message)
}

func TestCheckCode_RejectsInvalidRequestWithoutSending(t *testing.T) {
	cases := map[string]CheckCodeRequest{
		"empty code":    {RoleId: 1},
		"empty role id": {Code: "ABC"},
	}
	for name, req := range cases {
		t.Run(name, func(t *testing.T) {
			svc, captured := newCheckCodeTestService(t, http.StatusOK, checkCodeSuccessBody)

			resp, result, err := svc.CheckCode(context.Background(), req)

			require.Error(t, err)
			assert.Nil(t, resp)
			assert.Nil(t, result)
			assert.Empty(t, captured.uri)
		})
	}
}
