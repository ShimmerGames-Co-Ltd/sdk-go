package option

import (
	"gitlab.shimmer.lab/platform/sdk-go/core"
	"gitlab.shimmer.lab/platform/sdk-go/core/auth/sign"
	"gitlab.shimmer.lab/platform/sdk-go/core/consts"
	"net/http"
)

type RegionType string

const (
	RegionCN     RegionType = consts.PlatformRegionCN
	RegionGlobal RegionType = consts.PlatformRegionGlobal
)

type withHttpClientOption struct {
	Client *http.Client
}

func (w withHttpClientOption) Apply(o *core.ClientSettings) error {
	o.HttpClient = w.Client
	return nil
}

func WithHttpClient(client *http.Client) core.ClientOption {
	return withRegionOption{
		settings: core.ClientSettings{
			HttpClient: client,
		},
	}
}

type withRegionOption struct {
	settings core.ClientSettings
}

func (w withRegionOption) Apply(o *core.ClientSettings) error {
	o.Region = w.settings.Region
	return nil
}

func WithRegion(region RegionType) core.ClientOption {
	return withRegionOption{
		settings: core.ClientSettings{
			Region: string(region),
		},
	}
}

type withAuthCipherOption struct {
	settings core.ClientSettings
}

func (w withAuthCipherOption) Apply(o *core.ClientSettings) error {
	o.Signer = w.settings.Signer
	o.AppId = w.settings.AppId
	o.OrganizationId = w.settings.OrganizationId
	return nil
}

func WithAuthCipher(organizationId, appid, secret string) core.ClientOption {
	return withAuthCipherOption{
		settings: core.ClientSettings{
			Signer: &sign.SHA256WithRSASign{
				OrganizationId: organizationId,
				AppId:          appid,
				Secret:         secret,
			},
			AppId:          appid,
			OrganizationId: organizationId,
		},
	}
}

type withUrlOption struct {
	settings core.ClientSettings
}

func (w withUrlOption) Apply(o *core.ClientSettings) error {
	o.Url = w.settings.Url
	return nil
}

func WithUrl(url string) core.ClientOption {
	return withUrlOption{
		settings: core.ClientSettings{
			Url: url,
		},
	}
}
