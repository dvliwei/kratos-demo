package server

import (
	"github.com/google/wire"
)

// ProviderSet is server providers.
// 服务提供者集合，包含了GRPC和HTTP服务器的构造函数
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)
