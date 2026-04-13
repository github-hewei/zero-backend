package command

import (
	"encoding/json"
	"fmt"

	"zero-backend/internal/constants"
	"zero-backend/pkg/queue"

	"github.com/spf13/cobra"
)

// QueueCommand 队列命令
type QueueCommand struct {
	*cobra.Command
	manager *queue.QueueManager
}

// NewQueueCommand 创建队列命令
func NewQueueCommand(manager *queue.QueueManager) *QueueCommand {
	cmd := &QueueCommand{
		Command: &cobra.Command{
			Use:   "queue",
			Short: "队列管理",
			Long:  `队列管理命令行工具，用于向队列中写入测试数据`,
		},
		manager: manager,
	}

	cmd.AddCommand(newQueuePublishCommand(manager))
	return cmd
}

// newQueuePublishCommand 创建发布测试任务子命令
func newQueuePublishCommand(manager *queue.QueueManager) *cobra.Command {
	var (
		queueName string
		taskType  string
		count     int
		message   string
	)

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "发布测试任务",
		Long:  `向指定队列发布测试任务，用于验证 Worker 模块是否正常消费`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			for i := 0; i < count; i++ {
				payload, _ := json.Marshal(map[string]any{
					"index":   i + 1,
					"message": message,
				})

				task := queue.NewTask(queueName, taskType, payload)
				if err := manager.EnqueueTask(ctx, queueName, task); err != nil {
					return fmt.Errorf("enqueue task %d failed: %w", i+1, err)
				}

				cmd.Printf("✓ Published task [%d] id=%s type=%s queue=%s\n",
					i+1, task.ID, taskType, queueName)
			}

			cmd.Printf("\nTotal %d task(s) published to queue [%s]\n", count, queueName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&queueName, "queue", "q", constants.QueueDefaultKey, "队列名称")
	cmd.Flags().StringVarP(&taskType, "type", "t", "example", "任务类型")
	cmd.Flags().IntVarP(&count, "count", "n", 1, "发布任务数量")
	cmd.Flags().StringVarP(&message, "message", "m", "hello from cli", "任务消息内容")

	return cmd
}
