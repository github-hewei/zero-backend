package providers

import (
	"zero-backend/internal/middleware"
	middleware2 "zero-backend/modules/admin/middleware"
	middleware3 "zero-backend/modules/api/middleware"

	"github.com/google/wire"
)

// MiddlewaresProviderSet 提供中间件依赖集合
var MiddlewaresProviderSet = wire.NewSet(
	middleware.NewBeforeMiddleware,
	middleware.NewCorsMiddleware,
	wire.Struct(new(middleware.Middlewares), "*"),
)

// AdminMiddlewaresProviderSet 提供管理端中间件依赖集合
var AdminMiddlewaresProviderSet = wire.NewSet(
	middleware2.NewAuthMiddleware,
	wire.Struct(new(middleware2.Middlewares), "*"),
)

// ApiMiddlewaresProviderSet 提供开放接口中间件依赖集合
var ApiMiddlewaresProviderSet = wire.NewSet(
	middleware3.NewAuthMiddleware,
	wire.Struct(new(middleware3.Middlewares), "*"),
)
