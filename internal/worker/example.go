package worker

import (
	"context"
	"encoding/json"

	"github.com/241x/zero-kit/queue"

	"github.com/241x/zero-kit/logger"
)

// ExampleHandler 示例任务处理器。
type ExampleHandler struct {
	L logger.Logger
}

// Handle 处理示例任务。
func (h *ExampleHandler) Handle(ctx context.Context, task *queue.Task) error {
	var payload map[string]any
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		h.L.Error("failed to unmarshal example task payload", "task_id", task.ID, "error", err)
		return err
	}
	h.L.Info("processing example task", "task_id", task.ID, "payload", payload)
	return nil
}
