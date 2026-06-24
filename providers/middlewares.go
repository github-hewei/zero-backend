package providers

import (
	middleware1 "zero-backend/modules/admin/middleware"
	middleware2 "zero-backend/modules/api/middleware"

	"github.com/google/wire"
)

// AdminMiddlewaresProviderSet 提供管理端中间件依赖集合
var AdminMiddlewaresProviderSet = wire.NewSet(
	middleware1.NewAuthMiddleware,
	wire.Struct(new(middleware1.Middlewares), "*"),
)

// ApiMiddlewaresProviderSet 提供开放接口中间件依赖集合
var ApiMiddlewaresProviderSet = wire.NewSet(
	middleware2.NewAuthMiddleware,
	wire.Struct(new(middleware2.Middlewares), "*"),
)
