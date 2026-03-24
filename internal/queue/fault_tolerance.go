package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// FaultToleranceManager 容错管理器
type FaultToleranceManager struct {
	client           *redis.Client
	recovery         *TaskRecovery
	deduplicator     *Deduplicator
	gracefulShutdown *GracefulShutdown
	monitoring       bool
	mutex            sync.RWMutex
}

// NewFaultToleranceManager 创建容错管理器
func NewFaultToleranceManager(client *redis.Client) *FaultToleranceManager {
	return &FaultToleranceManager{
		client:           client,
		recovery:         NewTaskRecovery(client),
		deduplicator:     NewDeduplicator(client),
		gracefulShutdown: NewGracefulShutdown(30 * time.Second),
		monitoring:       false,
	}
}

// StartMonitoring 启动容错监控
func (f *FaultToleranceManager) StartMonitoring(ctx context.Context) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if f.monitoring {
		return errors.New("monitoring is already running")
	}

	f.monitoring = true

	// 启动定期恢复检查
	go f.recoveryMonitor(ctx)

	// 启动防重复清理
	go f.deduplicationCleanup(ctx)

	// 启动死信队列监控
	go f.deadLetterMonitor(ctx)

	fmt.Println("Fault tolerance monitoring started")
	return nil
}

// StopMonitoring 停止容错监控
func (f *FaultToleranceManager) StopMonitoring() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.monitoring = false
	fmt.Println("Fault tolerance monitoring stopped")
}

// recoveryMonitor 恢复监控
func (f *FaultToleranceManager) recoveryMonitor(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.mutex.RLock()
			if !f.monitoring {
				f.mutex.RUnlock()
				return
			}
			f.mutex.RUnlock()

			// 这里可以添加具体的恢复逻辑
			// 例如：检查并恢复停滞的任务
			fmt.Println("Recovery monitor tick")
		}
	}
}

// deduplicationCleanup 防重复清理
func (f *FaultToleranceManager) deduplicationCleanup(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.mutex.RLock()
			if !f.monitoring {
				f.mutex.RUnlock()
				return
			}
			f.mutex.RUnlock()

			if err := f.deduplicator.CleanupOldKeys(ctx); err != nil {
				fmt.Printf("Deduplication cleanup failed: %v\n", err)
			} else {
				fmt.Println("Deduplication cleanup completed")
			}
		}
	}
}

// deadLetterMonitor 死信队列监控
func (f *FaultToleranceManager) deadLetterMonitor(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.mutex.RLock()
			if !f.monitoring {
				f.mutex.RUnlock()
				return
			}
			f.mutex.RUnlock()

			// 检查死信队列并报警
			f.checkDeadLetterQueues(ctx)
		}
	}
}

// checkDeadLetterQueues 检查死信队列
func (f *FaultToleranceManager) checkDeadLetterQueues(ctx context.Context) {
	// 扫描所有死信队列
	pattern := "ZAG:QUEUE:*:DEAD"
	iter := f.client.Scan(ctx, 0, pattern, 100).Iterator()

	var deadLetterQueues []string
	for iter.Next(ctx) {
		deadLetterQueues = append(deadLetterQueues, iter.Val())
	}

	if err := iter.Err(); err != nil {
		fmt.Printf("Scan dead letter queues failed: %v\n", err)
		return
	}

	// 检查每个死信队列的长度
	for _, queueKey := range deadLetterQueues {
		length, err := f.client.LLen(ctx, queueKey).Result()
		if err != nil {
			fmt.Printf("Check dead letter queue %s failed: %v\n", queueKey, err)
			continue
		}

		if length > 0 {
			// 死信队列中有任务，发出警告
			fmt.Printf("WARNING: Dead letter queue %s has %d tasks\n", queueKey, length)

			// 这里可以添加报警逻辑，例如发送邮件、Slack消息等
			if length > 100 {
				fmt.Printf("CRITICAL: Dead letter queue %s has more than 100 tasks!\n", queueKey)
			}
		}
	}
}

// CrashRecovery 崩溃恢复
type CrashRecovery struct {
	client *redis.Client
}

// NewCrashRecovery 创建崩溃恢复器
func NewCrashRecovery(client *redis.Client) *CrashRecovery {
	return &CrashRecovery{
		client: client,
	}
}

// RecoverFromCrash 从崩溃中恢复
func (c *CrashRecovery) RecoverFromCrash(ctx context.Context, queueManager *QueueManager) (map[string]int, error) {
	fmt.Println("Starting crash recovery...")

	// 1. 恢复所有停滞的任务
	recovery := NewTaskRecovery(c.client)
	recoveredTasks, err := recovery.BatchRecovery(ctx, queueManager)
	if err != nil {
		return nil, fmt.Errorf("batch recovery failed: %w", err)
	}

	// 2. 清理孤儿锁
	orphanLocks, err := c.cleanupOrphanLocks(ctx)
	if err != nil {
		return nil, fmt.Errorf("cleanup orphan locks failed: %w", err)
	}

	// 3. 重置过期的领导者选举
	resetLeaders, err := c.resetExpiredLeaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("reset expired leaders failed: %w", err)
	}

	// 汇总恢复结果
	results := make(map[string]int)
	for queue, count := range recoveredTasks {
		results[fmt.Sprintf("queue_%s_recovered", queue)] = count
	}
	results["orphan_locks_cleaned"] = orphanLocks
	results["expired_leaders_reset"] = resetLeaders

	fmt.Printf("Crash recovery completed: %+v\n", results)
	return results, nil
}

// cleanupOrphanLocks 清理孤儿锁
func (c *CrashRecovery) cleanupOrphanLocks(ctx context.Context) (int, error) {
	pattern := "ZAG:LOCK:*"
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()

	var lockKeys []string
	for iter.Next(ctx) {
		lockKeys = append(lockKeys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return 0, err
	}

	cleanedCount := 0

	for _, lockKey := range lockKeys {
		// 获取锁的TTL
		ttl, err := c.client.TTL(ctx, lockKey).Result()
		if err != nil {
			continue
		}

		// 如果TTL为-2（键不存在）或-1（没有设置过期时间），跳过
		if ttl == -2*time.Second || ttl == -1*time.Second {
			continue
		}

		// 检查锁是否已经过期但未被清理
		// 这里可以根据实际情况添加更复杂的检查逻辑
		if ttl < 0 {
			// 删除孤儿锁
			if err := c.client.Del(ctx, lockKey).Err(); err != nil {
				continue
			}
			cleanedCount++
		}
	}

	return cleanedCount, nil
}

// resetExpiredLeaders 重置过期的领导者
func (c *CrashRecovery) resetExpiredLeaders(ctx context.Context) (int, error) {
	pattern := "ZAG:COORDINATOR:LEADER:*"
	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()

	var leaderKeys []string
	for iter.Next(ctx) {
		leaderKeys = append(leaderKeys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return 0, err
	}

	resetCount := 0

	for _, leaderKey := range leaderKeys {
		// 获取领导者的TTL
		ttl, err := c.client.TTL(ctx, leaderKey).Result()
		if err != nil {
			continue
		}

		// 如果领导者已过期，删除它
		if ttl == -2*time.Second || ttl < 0 {
			if err := c.client.Del(ctx, leaderKey).Err(); err != nil {
				continue
			}
			resetCount++
		}
	}

	return resetCount, nil
}

// DeadLetterManager 死信队列管理器
type DeadLetterManager struct {
	client *redis.Client
}

// NewDeadLetterManager 创建死信队列管理器
func NewDeadLetterManager(client *redis.Client) *DeadLetterManager {
	return &DeadLetterManager{
		client: client,
	}
}

// InspectDeadLetterQueue 检查死信队列
func (d *DeadLetterManager) InspectDeadLetterQueue(ctx context.Context, queueName string, limit int64) ([]*Task, error) {
	queueKey := fmt.Sprintf("ZAG:QUEUE:%s:DEAD", queueName)

	// 获取死信队列中的任务
	taskData, err := d.client.LRange(ctx, queueKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	tasks := make([]*Task, 0, len(taskData))
	for _, data := range taskData {
		task, err := UnmarshalTask([]byte(data))
		if err != nil {
			// 无法解析的任务，跳过
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// RetryDeadLetterTask 重试死信队列中的任务
func (d *DeadLetterManager) RetryDeadLetterTask(ctx context.Context, queueName string, taskID string) error {
	queueKey := fmt.Sprintf("ZAG:QUEUE:%s:DEAD", queueName)

	// 获取死信队列中的所有任务
	taskData, err := d.client.LRange(ctx, queueKey, 0, -1).Result()
	if err != nil {
		return err
	}

	// 查找指定任务
	for _, data := range taskData {
		task, err := UnmarshalTask([]byte(data))
		if err != nil {
			continue
		}

		if task.ID == taskID {
			// 从死信队列移除
			if err := d.client.LRem(ctx, queueKey, 1, data).Err(); err != nil {
				return err
			}

			// 重置任务状态
			task.Status = TaskStatusPending
			task.RetryCount = 0
			task.Error = ""

			// 重新入队
			immediateKey := fmt.Sprintf("ZAG:QUEUE:%s:IMMEDIATE", queueName)
			updatedData, err := task.Marshal()
			if err != nil {
				return err
			}

			return d.client.LPush(ctx, immediateKey, updatedData).Err()
		}
	}

	return errors.New("task not found in dead letter queue")
}

// PurgeDeadLetterQueue 清空死信队列
func (d *DeadLetterManager) PurgeDeadLetterQueue(ctx context.Context, queueName string) error {
	queueKey := fmt.Sprintf("ZAG:QUEUE:%s:DEAD", queueName)
	return d.client.Del(ctx, queueKey).Err()
}

// ExportDeadLetterTasks 导出死信队列任务
func (d *DeadLetterManager) ExportDeadLetterTasks(ctx context.Context, queueName string) ([]byte, error) {
	tasks, err := d.InspectDeadLetterQueue(ctx, queueName, 1000)
	if err != nil {
		return nil, err
	}

	// 将任务序列化为JSON
	type ExportTask struct {
		ID          string            `json:"id"`
		Queue       string            `json:"queue"`
		Type        string            `json:"type"`
		Payload     string            `json:"payload"`
		RetryCount  int               `json:"retry_count"`
		MaxRetries  int               `json:"max_retries"`
		CreatedAt   int64             `json:"created_at"`
		StartedAt   int64             `json:"started_at"`
		CompletedAt int64             `json:"completed_at"`
		Status      TaskStatus        `json:"status"`
		Error       string            `json:"error"`
		Metadata    map[string]string `json:"metadata"`
	}

	exportTasks := make([]ExportTask, len(tasks))
	for i, task := range tasks {
		exportTasks[i] = ExportTask{
			ID:          task.ID,
			Queue:       task.Queue,
			Type:        task.Type,
			Payload:     string(task.Payload),
			RetryCount:  task.RetryCount,
			MaxRetries:  task.MaxRetries,
			CreatedAt:   task.CreatedAt,
			StartedAt:   task.StartedAt,
			CompletedAt: task.CompletedAt,
			Status:      task.Status,
			Error:       task.Error,
			Metadata:    task.Metadata,
		}
	}

	// 使用JSON序列化
	return json.MarshalIndent(exportTasks, "", "  ")
}
