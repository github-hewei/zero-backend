package queue

import (
	"encoding/json"
	"time"
)

// Task 表示一个工作队列任务
type Task struct {
	// 任务ID
	ID string `json:"id"`

	// 队列名称
	Queue string `json:"queue"`

	// 任务数据
	Payload []byte `json:"payload"`

	// 任务类型
	Type string `json:"type"`

	// 最大重试次数
	MaxRetries int `json:"max_retries"`

	// 当前重试次数
	RetryCount int `json:"retry_count"`

	// 延迟执行时间（Unix时间戳，秒）
	DelayUntil int64 `json:"delay_until"`

	// 任务创建时间
	CreatedAt int64 `json:"created_at"`

	// 任务开始处理时间
	StartedAt int64 `json:"started_at"`

	// 任务完成时间
	CompletedAt int64 `json:"completed_at"`

	// 任务状态
	Status TaskStatus `json:"status"`

	// 错误信息
	Error string `json:"error"`

	// 元数据
	Metadata map[string]string `json:"metadata"`
}

// TaskStatus 任务状态
type TaskStatus string

const (
	// TaskStatusPending 等待处理
	TaskStatusPending TaskStatus = "pending"
	// TaskStatusProcessing 处理中
	TaskStatusProcessing TaskStatus = "processing"
	// TaskStatusCompleted 已完成
	TaskStatusCompleted TaskStatus = "completed"
	// TaskStatusFailed 已失败
	TaskStatusFailed TaskStatus = "failed"
	// TaskStatusDeadLetter 死信
	TaskStatusDeadLetter TaskStatus = "dead_letter"
)

// NewTask 创建新任务
func NewTask(queue, taskType string, payload []byte) *Task {
	return &Task{
		ID:         generateTaskID(),
		Queue:      queue,
		Type:       taskType,
		Payload:    payload,
		MaxRetries: 3,
		RetryCount: 0,
		DelayUntil: 0,
		CreatedAt:  time.Now().Unix(),
		Status:     TaskStatusPending,
		Metadata:   make(map[string]string),
	}
}

// WithDelay 设置延迟执行
func (t *Task) WithDelay(delay time.Duration) *Task {
	t.DelayUntil = time.Now().Add(delay).Unix()
	return t
}

// WithMaxRetries 设置最大重试次数
func (t *Task) WithMaxRetries(maxRetries int) *Task {
	t.MaxRetries = maxRetries
	return t
}

// WithMetadata 设置元数据
func (t *Task) WithMetadata(key, value string) *Task {
	if t.Metadata == nil {
		t.Metadata = make(map[string]string)
	}
	t.Metadata[key] = value
	return t
}

// Marshal 序列化任务为JSON
func (t *Task) Marshal() ([]byte, error) {
	return json.Marshal(t)
}

// UnmarshalTask 从JSON反序列化任务
func UnmarshalTask(data []byte) (*Task, error) {
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

// generateTaskID 生成任务ID
func generateTaskID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
