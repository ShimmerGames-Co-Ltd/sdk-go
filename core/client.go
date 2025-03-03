package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.shimmer.lab/platform/sdk-go/core/auth"
	"gitlab.shimmer.lab/platform/sdk-go/core/auth/credential"
	"gitlab.shimmer.lab/platform/sdk-go/core/consts"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
)

type APIResult struct {
	Request  *http.Request
	Response *http.Response
	PlatErr  *PlatError
}

type ClientOption interface {
	Apply(settings *ClientSettings) error
}

type Client struct {
	httpClient     *http.Client
	credential     auth.Credential
	url            string
	organizationId string
	appid          string
}

func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {

	settings, err := initSettings(opts)
	if err != nil {
		return nil, fmt.Errorf("init sdk setting failed. `%s`", err)
	}

	client := initClient(ctx, settings)

	return client, nil
}

func initSettings(opts []ClientOption) (settings *ClientSettings, err error) {

	var clientSettings ClientSettings

	for _, opt := range opts {
		if err = opt.Apply(&clientSettings); err != nil {
			return nil, err
		}
	}

	return &clientSettings, nil
}

func initClient(_ context.Context, settings *ClientSettings) *Client {

	var defaultUrl string

	if settings.Url == "" {
		switch settings.Region {
		case consts.PlatformRegionCN:
			defaultUrl = consts.PlatformCNApiServer
		case consts.PlatformRegionGlobal:
			defaultUrl = consts.PlatformGlobalApiServer
		default:
			defaultUrl = consts.PlatformGlobalApiServer
		}
	} else {
		defaultUrl = settings.Url
	}

	client := &Client{
		httpClient:     settings.HttpClient,
		url:            defaultUrl,
		credential:     &credential.PlatformCredential{Signer: settings.Signer},
		organizationId: settings.OrganizationId,
		appid:          settings.AppId,
	}

	if client.httpClient == nil {
		client.httpClient = &http.Client{Timeout: consts.HttpDefaultTimeout}
	}

	return client
}

func (client *Client) TokenSign(token string, timestamp int64) (string, error) {
	return client.credential.TokenSign(token, timestamp)
}

func (client *Client) RequestPost(ctx context.Context,
	path string,
	headers http.Header,
	body interface{}) (
	result *APIResult, err error) {

	uri, err := url.Parse(client.URL() + path)
	if err != nil {
		return
	}

	contentType := consts.HeaderContentTypeApplicationJSON

	bodyBuf := &bytes.Buffer{}
	if err = json.NewEncoder(bodyBuf).Encode(body); err != nil {
		return
	}

	return client.doRequestPost(ctx, uri.String(), contentType, headers, bodyBuf, bodyBuf.String())
}

func (client *Client) doRequestPost(ctx context.Context,
	requestURL string,
	contentType string,
	header http.Header,
	body io.Reader,
	signData string) (
	result *APIResult, err error) {

	var request *http.Request

	if request, err = http.NewRequestWithContext(ctx, http.MethodPost, requestURL, body); err != nil {
		return nil, err
	}

	for key, values := range header {
		for _, v := range values {
			request.Header.Add(key, v)
		}
	}

	request.Header.Set(consts.HeaderAppId, client.appid)
	request.Header.Set(consts.HeaderOrganizationId, client.organizationId)
	request.Header.Set(consts.HeaderContentType, contentType)
	request.Header.Set(consts.HeaderAccept, "*/*")
	request.Header.Set(consts.HeaderServiceSDKVersion, consts.HeaderServiceSDKVersionValue)
	request.Header.Set(consts.HeaderUserAgent, fmt.Sprintf(consts.HeaderUserAgentFormat, runtime.GOOS, runtime.Version()))

	var authorization string
	if authorization, err = client.credential.CreateAuthorizationHeader(request.URL.RequestURI(), signData); err != nil {
		return nil, fmt.Errorf("generate authorization err:%s", err.Error())
	}

	request.Header.Set(consts.HeaderAuthorization, authorization)

	result, err = client.doHTTP(request)
	if err != nil {
		return result, err
	}

	if err = CheckResponse(result.Response); err != nil {
		return result, err
	}

	return result, nil
}

func (client *Client) URL() string {
	return client.url
}

func (client *Client) doHTTP(req *http.Request) (result *APIResult, err error) {

	result = &APIResult{
		Request: req,
	}

	result.Response, err = client.httpClient.Do(req)

	return result, err
}

func CheckResponse(resp *http.Response) error {

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("invalid response, read body error: %w", err)
	}

	_ = resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewBuffer(data))

	apiError := &PlatError{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(data),
	}

	// 忽略 JSON 解析错误，均返回 apiError
	_ = json.Unmarshal(data, apiError)

	return apiError
}

func UnmarshalResponse(httpResp *http.Response, resp interface{}) error {
	body, err := ioutil.ReadAll(httpResp.Body)
	_ = httpResp.Body.Close()

	if err != nil {
		return err
	}

	httpResp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	err = json.Unmarshal(body, resp)
	if err != nil {
		return err
	}
	return nil
}
