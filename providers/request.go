package providers

import (
	"zero-backend/internal/request"

	"github.com/google/wire"
)

// RequestProviderSet 提供请求相关依赖集合
var RequestProviderSet = wire.NewSet(request.NewRequest, request.NewValidate, request.NewTrans)
