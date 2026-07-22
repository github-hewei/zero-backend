package worker

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-kit/queue"
)

const (
	queueDefaultKey = "default"
	queueTestKey    = "test"
)

// Server 队列工作服务，管理队列消费和工作线程池。
type Server struct {
	manager *queue.QueueManager
	logger  logger.Logger
}

// NewServer 创建队列工作服务。
func NewServer(manager *queue.QueueManager, registry *Registry, log logger.Logger) *Server {
	config := queue.DefaultConfig().WithName(queueDefaultKey)
	pool, err := manager.RegisterWorkerPool(queueDefaultKey, registry, config)
	if err != nil {
		log.Err(err, "Failed to register default worker pool")
	} else {
		pool.WithLogger(log)
	}
	return &Server{manager: manager, logger: log}
}

// Run 启动队列工作服务，阻塞直到收到退出信号。
func (s *Server) Run() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.manager.StartAllWorkerPools(ctx); err != nil {
		s.logger.Err(err, "Failed to start worker pools")
		return
	}

	s.logger.Info("Worker started, waiting for tasks...")

	<-sig
	s.logger.Info("Shutting down worker...")

	if err := s.manager.StopAllWorkerPools(); err != nil {
		s.logger.Err(err, "Failed to stop worker pools")
	}

	s.logger.Info("Worker stopped")
}
