package providers

import (
	"zero-backend/internal/service"
	service2 "zero-backend/modules/admin/service"
	service3 "zero-backend/modules/api/service"

	"github.com/google/wire"
)

// ServiceProviderSet 提供服务层依赖集合
var ServiceProviderSet = wire.NewSet(
	service.NewRbacMenuService,
	service.NewRbacApiService,
	service.NewRbacRoleService,
	service.NewRbacUserService,
	service.NewRbacStoreService,
	service.NewUploadGroupService,
	service.NewUploadFileService,
	service.NewUserService,
	service.NewSettingService,
	service.NewSettingDefaultService,
	service.NewArticleCategoryService,
	service.NewArticleService,
	service.NewRegionService,
)

// AdminServiceProviderSet 提供管理端服务层依赖集合
var AdminServiceProviderSet = wire.NewSet(service2.NewAuthService)

// ApiServiceProviderSet 提供管理端服务层依赖集合
var ApiServiceProviderSet = wire.NewSet(service3.NewAuthService)
