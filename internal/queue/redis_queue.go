package queue

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisQueue 基于Redis的工作队列实现
type RedisQueue struct {
	client *redis.Client
	config QueueConfig
	prefix string
}

// NewRedisQueue 创建新的Redis队列
func NewRedisQueue(client *redis.Client, config QueueConfig) *RedisQueue {
	if config.Name == "" {
		config = DefaultConfig()
	}

	return &RedisQueue{
		client: client,
		config: config,
		prefix: "ZAG:QUEUE:" + config.Name + ":",
	}
}

// Enqueue 立即入队任务
func (q *RedisQueue) Enqueue(ctx context.Context, task *Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	// 设置任务状态
	task.Status = TaskStatusPending
	task.CreatedAt = time.Now().Unix()

	// 序列化任务
	data, err := task.Marshal()
	if err != nil {
		return fmt.Errorf("marshal task failed: %w", err)
	}

	// 获取队列键
	queueKey := q.getImmediateQueueKey()

	// 将任务推入队列
	if err := q.client.LPush(ctx, queueKey, data).Err(); err != nil {
		return fmt.Errorf("push to queue failed: %w", err)
	}

	// 更新队列统计
	q.incrementStats(ctx, "enqueued")

	return nil
}

// EnqueueWithDelay 延迟入队任务
func (q *RedisQueue) EnqueueWithDelay(ctx context.Context, task *Task, delay time.Duration) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	// 设置延迟时间
	task.DelayUntil = time.Now().Add(delay).Unix()
	task.Status = TaskStatusPending
	task.CreatedAt = time.Now().Unix()

	// 序列化任务
	data, err := task.Marshal()
	if err != nil {
		return fmt.Errorf("marshal task failed: %w", err)
	}

	// 获取延迟队列键
	delayQueueKey := q.getDelayedQueueKey()

	// 将任务添加到有序集合（按执行时间排序）
	score := float64(task.DelayUntil)
	if err := q.client.ZAdd(ctx, delayQueueKey, redis.Z{
		Score:  score,
		Member: data,
	}).Err(); err != nil {
		return fmt.Errorf("add to delayed queue failed: %w", err)
	}

	// 更新队列统计
	q.incrementStats(ctx, "delayed")

	return nil
}

// Dequeue 出队任务
func (q *RedisQueue) Dequeue(ctx context.Context) (*Task, error) {
	// 首先检查延迟队列是否有到期的任务
	if err := q.moveDelayedTasks(ctx); err != nil {
		return nil, fmt.Errorf("move delayed tasks failed: %w", err)
	}

	// 获取立即执行队列键
	queueKey := q.getImmediateQueueKey()
	processingKey := q.getProcessingQueueKey()

	// 使用BRPOPLPUSH实现可靠出队（原子操作）
	// 从队列尾部弹出并添加到处理中队列
	result, err := q.client.BRPopLPush(ctx, queueKey, processingKey, 30*time.Second).Result()
	if err != nil {
		if err == redis.Nil {
			// 队列为空
			return nil, nil
		}
		return nil, fmt.Errorf("dequeue failed: %w", err)
	}

	// 反序列化任务
	task, err := UnmarshalTask([]byte(result))
	if err != nil {
		// 如果反序列化失败，从处理中队列移除
		q.client.LRem(ctx, processingKey, 1, result)
		return nil, fmt.Errorf("unmarshal task failed: %w", err)
	}

	// 更新任务状态
	task.Status = TaskStatusProcessing
	task.StartedAt = time.Now().Unix()

	// 将更新后的任务保存回处理中队列
	updatedData, err := task.Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshal updated task failed: %w", err)
	}

	// 替换处理中队列中的任务数据
	if err := q.client.LSet(ctx, processingKey, 0, updatedData).Err(); err != nil {
		return nil, fmt.Errorf("update processing task failed: %w", err)
	}

	// 设置任务可见性超时
	timeoutKey := q.getTaskTimeoutKey(task.ID)
	if err := q.client.Set(ctx, timeoutKey, "1", time.Duration(q.config.VisibilityTimeout)*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("set task timeout failed: %w", err)
	}

	// 更新队列统计
	q.incrementStats(ctx, "dequeued")

	return task, nil
}

// Ack 确认任务完成
func (q *RedisQueue) Ack(ctx context.Context, task *Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	// 从处理中队列移除
	processingKey := q.getProcessingQueueKey()

	// 在更新任务状态之前序列化任务，确保与处理中队列中的数据匹配
	// 保存当前状态
	originalStatus := task.Status
	originalCompletedAt := task.CompletedAt

	// 临时恢复为处理中状态以便匹配
	task.Status = TaskStatusProcessing
	task.CompletedAt = 0

	// 序列化用于匹配
	data, err := task.Marshal()
	if err != nil {
		// 恢复原始状态
		task.Status = originalStatus
		task.CompletedAt = originalCompletedAt
		return fmt.Errorf("marshal task failed: %w", err)
	}

	// 从处理中队列移除
	if err := q.client.LRem(ctx, processingKey, 1, data).Err(); err != nil {
		// 恢复原始状态
		task.Status = originalStatus
		task.CompletedAt = originalCompletedAt
		return fmt.Errorf("remove from processing queue failed: %w", err)
	}

	// 现在更新任务状态为完成
	task.Status = TaskStatusCompleted
	task.CompletedAt = time.Now().Unix()

	// 清理超时键
	timeoutKey := q.getTaskTimeoutKey(task.ID)
	q.client.Del(ctx, timeoutKey)

	// 更新队列统计
	q.incrementStats(ctx, "completed")

	return nil
}

// Nack 确认任务失败
func (q *RedisQueue) Nack(ctx context.Context, task *Task, failureErr error) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	// 保存原始错误信息
	errorMsg := ""
	if failureErr != nil {
		errorMsg = failureErr.Error()
	}

	// 更新任务状态
	task.Status = TaskStatusFailed
	task.Error = errorMsg
	task.RetryCount++

	// 检查是否超过最大重试次数
	if task.RetryCount > task.MaxRetries {
		// 移动到死信队列
		return q.MoveToDeadLetter(ctx, task, "max retries exceeded")
	}

	// 计算重试延迟
	delay := q.calculateRetryDelay(task.RetryCount)

	// 从处理中队列移除
	processingKey := q.getProcessingQueueKey()

	// 在更新任务状态之前序列化任务，确保与处理中队列中的数据匹配
	// 保存当前状态
	originalStatus := task.Status
	originalError := task.Error
	originalRetryCount := task.RetryCount

	// 临时恢复为处理中状态以便匹配
	task.Status = TaskStatusProcessing
	task.Error = ""
	task.RetryCount = task.RetryCount - 1 // 因为上面已经递增了，所以这里减1

	// 序列化用于匹配
	data, marshalErr := task.Marshal()
	if marshalErr != nil {
		// 恢复原始状态
		task.Status = originalStatus
		task.Error = originalError
		task.RetryCount = originalRetryCount
		return fmt.Errorf("marshal task failed: %w", marshalErr)
	}

	// 从处理中队列移除
	if err := q.client.LRem(ctx, processingKey, 1, data).Err(); err != nil {
		// 恢复原始状态
		task.Status = originalStatus
		task.Error = originalError
		task.RetryCount = originalRetryCount
		return fmt.Errorf("remove from processing queue failed: %w", err)
	}

	// 恢复任务状态并更新
	task.Status = TaskStatusFailed
	task.Error = errorMsg
	task.RetryCount = originalRetryCount + 1

	// 清理超时键
	timeoutKey := q.getTaskTimeoutKey(task.ID)
	q.client.Del(ctx, timeoutKey)

	// 重新入队（延迟重试）
	return q.EnqueueWithDelay(ctx, task, delay)
}

// Retry 重试任务
func (q *RedisQueue) Retry(ctx context.Context, task *Task, delay time.Duration) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	// 更新重试计数
	task.RetryCount++
	task.Status = TaskStatusPending
	task.Error = ""

	// 重新入队
	return q.EnqueueWithDelay(ctx, task, delay)
}

// MoveToDeadLetter 移动到死信队列
func (q *RedisQueue) MoveToDeadLetter(ctx context.Context, task *Task, reason string) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	if !q.config.EnableDeadLetter {
		// 如果不启用死信队列，直接丢弃任务
		return nil
	}

	// 更新任务状态
	task.Status = TaskStatusDeadLetter
	task.Error = reason
	task.CompletedAt = time.Now().Unix()

	// 序列化任务
	data, err := task.Marshal()
	if err != nil {
		return fmt.Errorf("marshal task failed: %w", err)
	}

	// 获取死信队列键
	deadLetterKey := q.getDeadLetterQueueKey()

	// 添加到死信队列
	if err := q.client.LPush(ctx, deadLetterKey, data).Err(); err != nil {
		return fmt.Errorf("add to dead letter queue failed: %w", err)
	}

	// 从处理中队列移除（如果存在）
	processingKey := q.getProcessingQueueKey()
	q.client.LRem(ctx, processingKey, 1, data)

	// 清理超时键
	timeoutKey := q.getTaskTimeoutKey(task.ID)
	q.client.Del(ctx, timeoutKey)

	// 更新队列统计
	q.incrementStats(ctx, "dead_letter")

	return nil
}

// GetStats 获取队列统计信息
func (q *RedisQueue) GetStats(ctx context.Context) (*QueueStats, error) {
	stats := &QueueStats{
		Name: q.config.Name,
	}

	// 获取各个队列的长度
	immediateKey := q.getImmediateQueueKey()
	delayedKey := q.getDelayedQueueKey()
	processingKey := q.getProcessingQueueKey()
	deadLetterKey := q.getDeadLetterQueueKey()

	var err error
	stats.Pending, err = q.client.LLen(ctx, immediateKey).Result()
	if err != nil {
		return nil, err
	}

	stats.Delayed, err = q.client.ZCard(ctx, delayedKey).Result()
	if err != nil {
		return nil, err
	}

	stats.Processing, err = q.client.LLen(ctx, processingKey).Result()
	if err != nil {
		return nil, err
	}

	stats.DeadLetter, err = q.client.LLen(ctx, deadLetterKey).Result()
	if err != nil {
		return nil, err
	}

	// 从统计哈希中获取其他统计信息
	statsKey := q.getStatsKey()
	statsMap, err := q.client.HGetAll(ctx, statsKey).Result()
	if err != nil {
		return nil, err
	}

	// 解析统计信息
	if val, ok := statsMap["enqueued"]; ok {
		stats.TotalEnqueued, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, ok := statsMap["dequeued"]; ok {
		stats.TotalDequeued, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, ok := statsMap["completed"]; ok {
		stats.Completed, _ = strconv.ParseInt(val, 10, 64)
	}
	if val, ok := statsMap["failed"]; ok {
		stats.Failed, _ = strconv.ParseInt(val, 10, 64)
	}

	return stats, nil
}

// Purge 清空队列
func (q *RedisQueue) Purge(ctx context.Context) error {
	keys := []string{
		q.getImmediateQueueKey(),
		q.getDelayedQueueKey(),
		q.getProcessingQueueKey(),
		q.getDeadLetterQueueKey(),
		q.getStatsKey(),
	}

	for _, key := range keys {
		if err := q.client.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("delete key %s failed: %w", key, err)
		}
	}

	return nil
}

// Close 关闭队列
func (q *RedisQueue) Close() error {
	// Redis客户端由外部管理，这里不需要关闭
	return nil
}

// GetClient 获取Redis客户端（用于需要直接访问客户端的场景）
func (q *RedisQueue) GetClient() *redis.Client {
	return q.client
}

// moveDelayedTasks 将到期的延迟任务移动到立即执行队列
func (q *RedisQueue) moveDelayedTasks(ctx context.Context) error {
	delayedKey := q.getDelayedQueueKey()
	immediateKey := q.getImmediateQueueKey()

	// 获取当前时间之前的所有任务
	now := float64(time.Now().Unix())

	// 获取到期的任务
	tasks, err := q.client.ZRangeByScore(ctx, delayedKey, &redis.ZRangeBy{
		Min:   "0",
		Max:   strconv.FormatFloat(now, 'f', -1, 64),
		Count: 100, // 每次最多移动100个任务
	}).Result()

	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	// 将任务添加到立即执行队列
	for _, taskData := range tasks {
		if err := q.client.LPush(ctx, immediateKey, taskData).Err(); err != nil {
			return err
		}
		// 从延迟队列移除
		if err := q.client.ZRem(ctx, delayedKey, taskData).Err(); err != nil {
			return err
		}
	}

	return nil
}

// calculateRetryDelay 计算重试延迟
func (q *RedisQueue) calculateRetryDelay(retryCount int) time.Duration {
	switch q.config.RetryDelay {
	case RetryDelayFixed:
		return 5 * time.Second
	case RetryDelayExponential:
		// 指数退避：2^retryCount 秒，最大1小时
		delay := time.Duration(1<<uint(retryCount)) * time.Second
		if delay > time.Hour {
			return time.Hour
		}
		return delay
	case RetryDelayRandom:
		// 随机延迟：1-10秒
		return time.Duration(1+time.Now().UnixNano()%10) * time.Second
	default:
		return 5 * time.Second
	}
}

// incrementStats 增加统计计数
func (q *RedisQueue) incrementStats(ctx context.Context, field string) {
	statsKey := q.getStatsKey()
	q.client.HIncrBy(ctx, statsKey, field, 1)
}

// getImmediateQueueKey 获取立即执行队列键
func (q *RedisQueue) getImmediateQueueKey() string {
	return q.prefix + "IMMEDIATE"
}

// getDelayedQueueKey 获取延迟队列键
func (q *RedisQueue) getDelayedQueueKey() string {
	return q.prefix + "DELAYED"
}

// getProcessingQueueKey 获取处理中队列键
func (q *RedisQueue) getProcessingQueueKey() string {
	return q.prefix + "PROCESSING"
}

// getDeadLetterQueueKey 获取死信队列键
func (q *RedisQueue) getDeadLetterQueueKey() string {
	return q.prefix + "DEAD"
}

// getStatsKey 获取统计键
func (q *RedisQueue) getStatsKey() string {
	return q.prefix + "STATS"
}

// getTaskTimeoutKey 获取任务超时键
func (q *RedisQueue) getTaskTimeoutKey(taskID string) string {
	return q.prefix + "TIMEOUT:" + taskID
}

// getLockKey 获取分布式锁键
func (q *RedisQueue) getLockKey(lockName string) string {
	return q.prefix + "LOCK:" + lockName
}
