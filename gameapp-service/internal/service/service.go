package service

import "github.com/google/wire"

// ProviderSet is service providers.
// 服务层的依赖注入
// 服务层的依赖注入，需要在服务层的代码中使用 wire.NewSet 注册依赖
var ProviderSet = wire.NewSet(NewGreeterService, NewGameAppService)
