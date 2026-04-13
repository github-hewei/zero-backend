package handler

import (
	"context"
	"fmt"

	"zero-backend/pkg/logger"
	"zero-backend/pkg/queue"
)

// Registry 任务处理器注册中心，根据任务类型分发到对应的处理器
// 实现 queue.Handler 接口，作为工作线程池的统一处理器入口
type Registry struct {
	handlers map[string]queue.Handler
	logger   logger.Logger
}

// NewRegistry 创建处理器注册中心
func NewRegistry(log logger.Logger) *Registry {
	return &Registry{
		handlers: make(map[string]queue.Handler),
		logger:   log,
	}
}

// Register 注册任务处理器
func (r *Registry) Register(taskType string, handler queue.Handler) {
	r.handlers[taskType] = handler
	r.logger.Info("registered task handler", "task_type", taskType)
}

// Handle 实现 queue.Handler 接口，根据任务类型分发到对应的处理器
func (r *Registry) Handle(ctx context.Context, task *queue.Task) error {
	h, ok := r.handlers[task.Type]
	if !ok {
		return fmt.Errorf("no handler registered for task type: %s", task.Type)
	}
	return h.Handle(ctx, task)
}
