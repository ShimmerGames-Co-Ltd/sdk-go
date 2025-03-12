# SDK-GO

平台SDK集成库 (对应Golang版本)

## 功能介绍
1. 接口SDK
2. HTTP 客户端支持, 请求签名

## 开始
```shell
go get -u github.com/ShimmerGames-Co-Ltd/sdk-go
```

### 创建请求实例
```go

    ctx := context.Background()
        sdkOpts := []core.ClientOption{
	// 选择国内平台地址
        option.WithRegion(option.RegionCN),
        // 选择全球平台地址
        // option.WithRegion(option.RegionGlobal),
	// 自定义平台地址
        // option.WithUrl("https://example.platform.com"), 
        // 服务侧验证信息
        option.WithAuthCipher("YOUR_ORGANIZATIONID", "YOUR_APPID", "APP_SECRET"),
	}

	client, err := core.NewClient(ctx, sdkOpts...)
	if err != nil {
		err = fmt.Errorf("create sdk client failed. `%s`", err)
		return
	}
```

####  验证用户登录
```go
        // 初始化认证服务
        svc := authorization.AuthUserService{Client: client}
        // 发起请求
        resp, result, err := svc.UserAuthorize(ctx, authorization.UserAuthorizeRequest{Token: message.Token})
        if err != nil {
            log.Errorf("authorize failed. `%s`", err)
            return
        }
```

####  验证订单
```go
        // 初始化认证服务
        svc := iap.PaymentService{Client: client}
        // 发起请求
        resp, result, err := svc.VerifyPayment(ctx, iap.VerifyPaymentRequest{OrderId: message.OrderId})
        if err != nil {
            log.Errorf("verify failed. `%s`", err)
            return
        }
```