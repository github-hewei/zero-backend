package providers

import (
	"zero-backend/modules/rbac"
	"zero-backend/modules/user"

	"github.com/google/wire"
)

// RepositoryProviderSet 提供仓储层依赖集合
var RepositoryProviderSet = wire.NewSet(
	user.NewRepository,
	rbac.NewRbacApiRepository,
)
