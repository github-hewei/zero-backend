package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"zero-backend/internal/constants"
	"zero-backend/pkg/logger"
	"zero-backend/pkg/queue"

	"zero-backend/modules/worker/handler"
)

// WorkerServer 队列工作服务，管理队列消费和工作线程池
type WorkerServer struct {
	manager *queue.QueueManager
	logger  logger.Logger
}

// NewWorkerServer 创建队列工作服务
func NewWorkerServer(manager *queue.QueueManager, registry *handler.Registry, log logger.Logger) *WorkerServer {
	// 注册默认队列工作线程池
	config := queue.DefaultConfig().WithName(constants.QueueDefaultKey)
	pool, err := manager.RegisterWorkerPool(constants.QueueDefaultKey, registry, config)
	if err != nil {
		log.Err(err, "Failed to register default worker pool")
	} else {
		pool.WithLogger(log)
	}

	return &WorkerServer{
		manager: manager,
		logger:  log,
	}
}

// Run 启动队列工作服务，阻塞直到收到退出信号
func (s *WorkerServer) Run() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动所有工作线程池
	if err := s.manager.StartAllWorkerPools(ctx); err != nil {
		s.logger.Err(err, "Failed to start worker pools")
		return
	}

	s.logger.Info("Worker started, waiting for tasks...")

	// 等待退出信号
	<-sig
	s.logger.Info("Shutting down worker...")

	// 优雅停止所有工作线程池
	if err := s.manager.StopAllWorkerPools(); err != nil {
		s.logger.Err(err, "Failed to stop worker pools")
	}

	s.logger.Info("Worker stopped")
}
