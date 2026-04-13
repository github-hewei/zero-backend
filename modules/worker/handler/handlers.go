package handler

import (
	"context"
	"encoding/json"

	"zero-backend/pkg/logger"
	"zero-backend/pkg/queue"
)

// ExampleHandler 示例任务处理器，处理 "example" 类型的任务
// 可作为模板参考，实际使用时替换为具体业务处理器
type ExampleHandler struct {
	logger logger.Logger
}

// NewExampleHandler 创建示例任务处理器
func NewExampleHandler(log logger.Logger) *ExampleHandler {
	return &ExampleHandler{logger: log}
}

// Handle 处理示例任务
func (h *ExampleHandler) Handle(ctx context.Context, task *queue.Task) error {
	var payload map[string]any
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		h.logger.Error("failed to unmarshal example task payload",
			"task_id", task.ID, "error", err)
		return err
	}

	h.logger.Info("processing example task",
		"task_id", task.ID, "payload", payload)
	return nil
}
