package cli

import (
	"encoding/json"
	"fmt"

	"github.com/241x/zero-kit/queue"
	"github.com/spf13/cobra"
)

const (
	queueDefaultKey = "default"
	queueTestKey    = "test"
)

// QueueCmd 队列管理命令。
func QueueCmd(manager *queue.QueueManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queue",
		Short: "队列管理",
		Long:  `队列管理命令行工具，用于向队列中写入测试数据`,
	}
	cmd.AddCommand(queuePublishCmd(manager))
	return cmd
}

// queuePublishCmd 发布测试任务命令。
func queuePublishCmd(manager *queue.QueueManager) *cobra.Command {
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

	cmd.Flags().StringVarP(&queueName, "queue", "q", queueDefaultKey, "队列名称")
	cmd.Flags().StringVarP(&taskType, "type", "t", "example", "任务类型")
	cmd.Flags().IntVarP(&count, "count", "n", 1, "发布任务数量")
	cmd.Flags().StringVarP(&message, "message", "m", "hello from cli", "任务消息内容")

	return cmd
}
