package providers

import (
	middleware2 "zero-backend/modules/api/middleware"

	"github.com/google/wire"
)

// ApiMiddlewaresProviderSet 提供开放接口中间件依赖集合
var ApiMiddlewaresProviderSet = wire.NewSet(
	middleware2.NewAuthMiddleware,
	wire.Struct(new(middleware2.Middlewares), "*"),
)
