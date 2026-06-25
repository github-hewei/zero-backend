package providers

import (
	"zero-backend/modules/admin/controller"
	controller2 "zero-backend/modules/api/controller"

	"github.com/google/wire"
)

// AdminControllersProviderSet 提供管理端控制器依赖集合
var AdminControllersProviderSet = wire.NewSet(
	controller.NewAuthController,
	controller.NewRbacMenuController,
	controller.NewRbacApiController,
	controller.NewRbacRoleController,
	controller.NewRbacUserController,
	controller.NewRbacStoreController,
	controller.NewUserController,
	controller.NewSettingController,
	controller.NewSettingDefaultController,
	controller.NewHealthController,
	wire.Struct(new(controller.Controllers), "*"),
)

// ApiControllersProviderSet 提供API控制器依赖集合
var ApiControllersProviderSet = wire.NewSet(
	controller2.NewAuthController,
	controller2.NewSettingController,
	wire.Struct(new(controller2.Controllers), "*"),
)
