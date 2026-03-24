package queue

import (
	"context"
	"time"
)

// Handler 任务处理器接口
type Handler interface {
	// Handle 处理任务
	Handle(ctx context.Context, task *Task) error
}

// HandlerFunc 任务处理器函数类型
type HandlerFunc func(ctx context.Context, task *Task) error

// Handle 实现Handler接口
func (f HandlerFunc) Handle(ctx context.Context, task *Task) error {
	return f(ctx, task)
}

// Queue 工作队列接口
type Queue interface {
	// Enqueue 立即入队任务
	Enqueue(ctx context.Context, task *Task) error

	// EnqueueWithDelay 延迟入队任务
	EnqueueWithDelay(ctx context.Context, task *Task, delay time.Duration) error

	// Dequeue 出队任务
	Dequeue(ctx context.Context) (*Task, error)

	// Ack 确认任务完成
	Ack(ctx context.Context, task *Task) error

	// Nack 确认任务失败
	Nack(ctx context.Context, task *Task, err error) error

	// Retry 重试任务
	Retry(ctx context.Context, task *Task, delay time.Duration) error

	// MoveToDeadLetter 移动到死信队列
	MoveToDeadLetter(ctx context.Context, task *Task, reason string) error

	// GetStats 获取队列统计信息
	GetStats(ctx context.Context) (*QueueStats, error)

	// Purge 清空队列
	Purge(ctx context.Context) error

	// Close 关闭队列
	Close() error
}

// QueueStats 队列统计信息
type QueueStats struct {
	// 队列名称
	Name string `json:"name"`

	// 等待处理的任务数
	Pending int64 `json:"pending"`

	// 处理中的任务数
	Processing int64 `json:"processing"`

	// 延迟队列中的任务数
	Delayed int64 `json:"delayed"`

	// 死信队列中的任务数
	DeadLetter int64 `json:"dead_letter"`

	// 已完成的任务数
	Completed int64 `json:"completed"`

	// 失败的任务数
	Failed int64 `json:"failed"`

	// 总入队任务数
	TotalEnqueued int64 `json:"total_enqueued"`

	// 总出队任务数
	TotalDequeued int64 `json:"total_dequeued"`

	// 平均处理时间（毫秒）
	AvgProcessingTime int64 `json:"avg_processing_time"`

	// 当前并发数
	CurrentConcurrency int `json:"current_concurrency"`
}

// Worker 工作线程接口
type Worker interface {
	// Start 启动工作线程
	Start(ctx context.Context) error

	// Stop 停止工作线程
	Stop() error

	// IsRunning 是否正在运行
	IsRunning() bool

	// GetStats 获取工作线程统计信息
	GetStats() *WorkerStats
}

// WorkerStats 工作线程统计信息
type WorkerStats struct {
	// 工作线程ID
	ID string `json:"id"`

	// 是否正在运行
	Running bool `json:"running"`

	// 已处理任务数
	ProcessedTasks int64 `json:"processed_tasks"`

	// 成功任务数
	SuccessfulTasks int64 `json:"successful_tasks"`

	// 失败任务数
	FailedTasks int64 `json:"failed_tasks"`

	// 当前处理中的任务
	CurrentTask *Task `json:"current_task"`

	// 启动时间
	StartedAt time.Time `json:"started_at"`

	// 最后活动时间
	LastActivity time.Time `json:"last_activity"`
}
