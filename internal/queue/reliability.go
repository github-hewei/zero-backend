package queue

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Deduplicator 防重复处理器
type Deduplicator struct {
	client *redis.Client
	prefix string
}

// NewDeduplicator 创建防重复处理器
func NewDeduplicator(client *redis.Client) *Deduplicator {
	return &Deduplicator{
		client: client,
		prefix: "ZAG:DEDUP:",
	}
}

// GenerateTaskKey 生成任务唯一键
func (d *Deduplicator) GenerateTaskKey(queueName string, payload []byte) string {
	hash := sha256.Sum256(payload)
	return queueName + ":" + hex.EncodeToString(hash[:])
}

// IsDuplicate 检查任务是否重复
func (d *Deduplicator) IsDuplicate(ctx context.Context, queueName string, payload []byte) (bool, error) {
	taskKey := d.GenerateTaskKey(queueName, payload)
	redisKey := d.prefix + taskKey

	// 使用SET NX命令检查是否已存在
	result, err := d.client.SetNX(ctx, redisKey, "1", 24*time.Hour).Result()
	if err != nil {
		return false, fmt.Errorf("check duplicate failed: %w", err)
	}

	// 如果设置成功，说明不是重复任务
	// 如果设置失败，说明是重复任务
	return !result, nil
}

// MarkAsProcessed 标记任务为已处理
func (d *Deduplicator) MarkAsProcessed(ctx context.Context, queueName string, payload []byte) error {
	taskKey := d.GenerateTaskKey(queueName, payload)
	redisKey := d.prefix + taskKey

	// 延长键的过期时间
	return d.client.Expire(ctx, redisKey, 24*time.Hour).Err()
}

// CleanupOldKeys 清理旧的防重复键
func (d *Deduplicator) CleanupOldKeys(ctx context.Context) error {
	// 使用SCAN迭代所有防重复键
	iter := d.client.Scan(ctx, 0, d.prefix+"*", 100).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// 批量删除
	if len(keys) > 0 {
		return d.client.Del(ctx, keys...).Err()
	}

	return nil
}

// TaskRecovery 任务恢复器
type TaskRecovery struct {
	client *redis.Client
}

// NewTaskRecovery 创建任务恢复器
func NewTaskRecovery(client *redis.Client) *TaskRecovery {
	return &TaskRecovery{
		client: client,
	}
}

// RecoverStalledTasks 恢复停滞的任务
func (r *TaskRecovery) RecoverStalledTasks(ctx context.Context, queue *RedisQueue) (int, error) {
	processingKey := queue.getProcessingQueueKey()

	// 获取所有处理中的任务
	tasks, err := r.client.LRange(ctx, processingKey, 0, -1).Result()
	if err != nil {
		return 0, err
	}

	recoveredCount := 0

	for _, taskData := range tasks {
		task, err := UnmarshalTask([]byte(taskData))
		if err != nil {
			// 无法解析的任务，从队列中移除
			r.client.LRem(ctx, processingKey, 1, taskData)
			continue
		}

		// 检查任务是否超时
		timeoutKey := queue.getTaskTimeoutKey(task.ID)
		exists, err := r.client.Exists(ctx, timeoutKey).Result()
		if err != nil {
			continue
		}

		if exists == 0 {
			// 任务已超时，需要恢复
			if err := r.recoverTask(ctx, queue, task, taskData); err != nil {
				fmt.Printf("Recover task %s failed: %v\n", task.ID, err)
				continue
			}
			recoveredCount++
		}
	}

	return recoveredCount, nil
}

// recoverTask 恢复单个任务
func (r *TaskRecovery) recoverTask(ctx context.Context, queue *RedisQueue, task *Task, taskData string) error {
	processingKey := queue.getProcessingQueueKey()

	// 检查任务是否超过最大重试次数
	if task.RetryCount >= task.MaxRetries {
		// 移动到死信队列
		if err := queue.MoveToDeadLetter(ctx, task, "max retries exceeded after recovery"); err != nil {
			return err
		}
		// 从处理中队列移除
		r.client.LRem(ctx, processingKey, 1, taskData)
		return nil
	}

	// 增加重试计数
	task.RetryCount++
	task.Status = TaskStatusPending
	task.Error = ""

	// 序列化更新后的任务
	updatedData, err := task.Marshal()
	if err != nil {
		return err
	}

	// 从处理中队列移除旧任务
	if err := r.client.LRem(ctx, processingKey, 1, taskData).Err(); err != nil {
		return err
	}

	// 重新入队
	immediateKey := queue.getImmediateQueueKey()
	if err := r.client.LPush(ctx, immediateKey, updatedData).Err(); err != nil {
		return err
	}

	fmt.Printf("Recovered task %s (retry %d/%d)\n", task.ID, task.RetryCount, task.MaxRetries)
	return nil
}

// BatchRecovery 批量恢复所有队列的停滞任务
func (r *TaskRecovery) BatchRecovery(ctx context.Context, queueManager *QueueManager) (map[string]int, error) {
	stats, err := queueManager.GetQueueStats(ctx)
	if err != nil {
		return nil, err
	}

	results := make(map[string]int)

	for queueName := range stats {
		queue := queueManager.GetOrCreateQueue(queueName, DefaultConfig())
		recovered, err := r.RecoverStalledTasks(ctx, queue)
		if err != nil {
			return nil, fmt.Errorf("recover tasks for queue %s failed: %w", queueName, err)
		}
		results[queueName] = recovered
	}

	return results, nil
}

// GracefulShutdown 优雅关闭处理器
type GracefulShutdown struct {
	workerPools map[string]*WorkerPool
	queues      map[string]*RedisQueue
	timeout     time.Duration
}

// NewGracefulShutdown 创建优雅关闭处理器
func NewGracefulShutdown(timeout time.Duration) *GracefulShutdown {
	return &GracefulShutdown{
		workerPools: make(map[string]*WorkerPool),
		queues:      make(map[string]*RedisQueue),
		timeout:     timeout,
	}
}

// RegisterWorkerPool 注册工作线程池
func (g *GracefulShutdown) RegisterWorkerPool(queueName string, pool *WorkerPool) {
	g.workerPools[queueName] = pool
}

// RegisterQueue 注册队列
func (g *GracefulShutdown) RegisterQueue(queueName string, queue *RedisQueue) {
	g.queues[queueName] = queue
}

// Shutdown 执行优雅关闭
func (g *GracefulShutdown) Shutdown(ctx context.Context) error {
	// 设置超时上下文
	shutdownCtx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	// 1. 停止接收新任务（标记队列为只读）
	fmt.Println("Stopping new task acceptance...")

	// 2. 停止所有工作线程池
	fmt.Println("Stopping worker pools...")
	var errs []error
	for queueName, pool := range g.workerPools {
		if err := pool.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("stop worker pool %s failed: %w", queueName, err))
		} else {
			fmt.Printf("Worker pool %s stopped\n", queueName)
		}
	}

	// 3. 等待处理中的任务完成
	fmt.Println("Waiting for in-progress tasks to complete...")
	time.Sleep(5 * time.Second)

	// 4. 恢复所有停滞的任务
	fmt.Println("Recovering stalled tasks...")
	recovery := NewTaskRecovery(g.getRedisClient())
	for queueName, queue := range g.queues {
		recovered, err := recovery.RecoverStalledTasks(shutdownCtx, queue)
		if err != nil {
			errs = append(errs, fmt.Errorf("recover stalled tasks for queue %s failed: %w", queueName, err))
		} else {
			fmt.Printf("Recovered %d stalled tasks from queue %s\n", recovered, queueName)
		}
	}

	// 5. 清理资源
	fmt.Println("Cleaning up resources...")

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	fmt.Println("Graceful shutdown completed")
	return nil
}

// getRedisClient 获取Redis客户端（假设所有队列使用同一个客户端）
func (g *GracefulShutdown) getRedisClient() *redis.Client {
	if len(g.queues) == 0 {
		return nil
	}

	// 获取第一个队列的客户端
	for _, queue := range g.queues {
		return queue.GetClient()
	}

	return nil
}
