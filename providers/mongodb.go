package providers

import (
	"zero-backend/internal/storage/mongodb"

	"github.com/google/wire"
)

// MongoDBProviderSet 提供MongoDB数据库依赖集合
var MongoDBProviderSet = wire.NewSet(
	mongodb.NewConn,
	wire.FieldsOf(new(*mongodb.Conn), "Client", "DB"),
)
