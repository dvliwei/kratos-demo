package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
// 业务逻辑层的依赖注入
// 业务逻辑层的依赖注入，需要在业务逻辑层的代码中使用 wire.NewSet 注册依赖
var ProviderSet = wire.NewSet(NewGreeterUsecase, NewGameAppUseCase)
