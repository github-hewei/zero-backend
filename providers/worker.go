package providers

import (
	"zero-backend/modules/worker/handler"
	"zero-backend/modules/worker/server"
	"zero-backend/pkg/logger"
	"zero-backend/pkg/queue"

	"github.com/google/wire"
)

// ProvideRegistry 创建任务处理器注册中心并注册所有任务处理器
func ProvideRegistry(log logger.Logger, exampleHandler *handler.ExampleHandler) *handler.Registry {
	registry := handler.NewRegistry(log)
	registry.Register("example", exampleHandler)
	return registry
}

// WorkerProviderSet 提供 Worker 模块依赖集合
var WorkerProviderSet = wire.NewSet(
	ProvideRegistry,
	handler.NewExampleHandler,
	queue.NewQueueManager,
	server.NewWorkerServer,
)
