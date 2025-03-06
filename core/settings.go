package core

import (
	"github.com/shimmergames/sdk-go/core/auth"
	"net/http"
)

type ClientSettings struct {

	// HttpClient Http客户端实例
	HttpClient *http.Client

	// Region 服务地区
	Region string

	// Signer 签名器
	Signer auth.Signer

	// 组织ID OrganizationId
	OrganizationId string

	// AppId 应用ID
	AppId string

	// Url 连接地址
	Url string
}
