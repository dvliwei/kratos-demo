package service

import "github.com/google/wire"

// ProviderSet is service providers.
// 注册服务
var ProviderSet = wire.NewSet(
	NewGreeterService,
	// 注册网关服务
	NewGatewayService,
)
