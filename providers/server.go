package providers

import (
	"zero-backend/modules/admin/server"
	server2 "zero-backend/modules/api/server"

	"github.com/google/wire"
)

// AdminServerProviderSet 提供管理端服务依赖集合
var AdminServerProviderSet = wire.NewSet(server.NewHTTPServer, server.NewGin)

// ApiServerProviderSet 提供API服务依赖集合
var ApiServerProviderSet = wire.NewSet(server2.NewHTTPServer, server2.NewGin)
