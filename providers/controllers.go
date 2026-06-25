package providers

import (
	controller2 "zero-backend/modules/api/controller"

	"github.com/google/wire"
)

// ApiControllersProviderSet 提供API控制器依赖集合
var ApiControllersProviderSet = wire.NewSet(
	controller2.NewAuthController,
	wire.Struct(new(controller2.Controllers), "*"),
)
