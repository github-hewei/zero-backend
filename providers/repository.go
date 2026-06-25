package providers

import (
	"zero-backend/internal/repository"
	"zero-backend/modules/rbac"

	"github.com/google/wire"
)

// RepositoryProviderSet 提供仓储层依赖集合
var RepositoryProviderSet = wire.NewSet(
	repository.NewUserRepository,
	rbac.NewRbacApiRepository,
)
