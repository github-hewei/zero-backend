package providers

import (
	"zero-backend/internal/repository"

	"github.com/google/wire"
)

// RepositoryProviderSet 提供仓储层依赖集合
var RepositoryProviderSet = wire.NewSet(
	repository.NewRbacMenuRepository,
	repository.NewRbacApiRepository,
	repository.NewRbacRoleRepository,
	repository.NewRbacUserRepository,
	repository.NewRbacStoreRepository,
	repository.NewRbacRoleMenuRepository,
	repository.NewRbacUserRoleRepository,
	repository.NewRbacMenuApiRepository,
	repository.NewUploadGroupRepository,
	repository.NewUploadFileRepository,
	repository.NewUserRepository,
	repository.NewSettingRepository,
	repository.NewSettingDefaultRepository,
	repository.NewRegionRepository,
	repository.NewUserPointsLogRepository,
	repository.NewArticleCategoryRepository,
	repository.NewArticleRepository,
)
