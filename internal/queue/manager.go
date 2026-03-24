package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// QueueManager 队列管理器
type QueueManager struct {
	client      *redis.Client
	queues      map[string]*RedisQueue
	workerPools map[string]*WorkerPool
	mutex       sync.RWMutex
}

// NewQueueManager 创建队列管理器
func NewQueueManager(client *redis.Client) *QueueManager {
	return &QueueManager{
		client:      client,
		queues:      make(map[string]*RedisQueue),
		workerPools: make(map[string]*WorkerPool),
	}
}

// GetOrCreateQueue 获取或创建队列
func (m *QueueManager) GetOrCreateQueue(name string, config QueueConfig) *RedisQueue {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if queue, exists := m.queues[name]; exists {
		return queue
	}

	config.Name = name
	queue := NewRedisQueue(m.client, config)
	m.queues[name] = queue

	return queue
}

// RegisterWorkerPool 注册工作线程池
func (m *QueueManager) RegisterWorkerPool(queueName string, handler Handler, config QueueConfig) (*WorkerPool, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查队列是否存在
	queue, exists := m.queues[queueName]
	if !exists {
		// 自动创建队列
		config.Name = queueName
		queue = NewRedisQueue(m.client, config)
		m.queues[queueName] = queue
	}

	// 创建工作线程池
	pool := NewWorkerPool(queue, handler, config)
	m.workerPools[queueName] = pool

	return pool, nil
}

// StartAllWorkerPools 启动所有工作线程池
func (m *QueueManager) StartAllWorkerPools(ctx context.Context) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var errs []error
	for queueName, pool := range m.workerPools {
		if err := pool.Start(ctx); err != nil {
			errs = append(errs, fmt.Errorf("start worker pool for queue %s failed: %w", queueName, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// StopAllWorkerPools 停止所有工作线程池
func (m *QueueManager) StopAllWorkerPools() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var errs []error
	for queueName, pool := range m.workerPools {
		if err := pool.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("stop worker pool for queue %s failed: %w", queueName, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// GetQueueStats 获取队列统计信息
func (m *QueueManager) GetQueueStats(ctx context.Context) (map[string]*QueueStats, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string]*QueueStats)
	for name, queue := range m.queues {
		queueStats, err := queue.GetStats(ctx)
		if err != nil {
			return nil, fmt.Errorf("get stats for queue %s failed: %w", name, err)
		}
		stats[name] = queueStats
	}

	return stats, nil
}

// GetWorkerPoolStats 获取工作线程池统计信息
func (m *QueueManager) GetWorkerPoolStats() map[string][]*WorkerStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	stats := make(map[string][]*WorkerStats)
	for queueName, pool := range m.workerPools {
		stats[queueName] = pool.GetWorkerStats()
	}

	return stats
}

// EnqueueTask 入队任务到指定队列
func (m *QueueManager) EnqueueTask(ctx context.Context, queueName string, task *Task) error {
	queue := m.GetOrCreateQueue(queueName, DefaultConfig())
	return queue.Enqueue(ctx, task)
}

// EnqueueTaskWithDelay 延迟入队任务到指定队列
func (m *QueueManager) EnqueueTaskWithDelay(ctx context.Context, queueName string, task *Task, delay time.Duration) error {
	queue := m.GetOrCreateQueue(queueName, DefaultConfig())
	return queue.EnqueueWithDelay(ctx, task, delay)
}

// CleanupStalledTasks 清理停滞的任务
func (m *QueueManager) CleanupStalledTasks(ctx context.Context) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var errs []error
	for queueName, queue := range m.queues {
		if err := m.cleanupQueueStalledTasks(ctx, queue); err != nil {
			errs = append(errs, fmt.Errorf("cleanup stalled tasks for queue %s failed: %w", queueName, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// cleanupQueueStalledTasks 清理单个队列的停滞任务
func (m *QueueManager) cleanupQueueStalledTasks(ctx context.Context, queue *RedisQueue) error {
	// 获取处理中队列的所有任务
	processingKey := queue.getProcessingQueueKey()
	tasks, err := m.client.LRange(ctx, processingKey, 0, -1).Result()
	if err != nil {
		return err
	}

	stalledCount := 0

	for _, taskData := range tasks {
		task, err := UnmarshalTask([]byte(taskData))
		if err != nil {
			// 无法解析的任务，从队列中移除
			m.client.LRem(ctx, processingKey, 1, taskData)
			continue
		}

		// 检查任务是否超时（超过可见性超时时间）
		timeoutKey := queue.getTaskTimeoutKey(task.ID)
		exists, err := m.client.Exists(ctx, timeoutKey).Result()
		if err != nil {
			continue
		}

		if exists == 0 {
			// 超时键不存在，说明任务已超时
			// 将任务重新放回队列
			task.RetryCount++
			task.Status = TaskStatusPending

			updatedData, err := task.Marshal()
			if err != nil {
				continue
			}

			// 从处理中队列移除
			m.client.LRem(ctx, processingKey, 1, taskData)

			// 重新入队
			immediateKey := queue.getImmediateQueueKey()
			m.client.LPush(ctx, immediateKey, updatedData)

			stalledCount++
		}
	}

	if stalledCount > 0 {
		fmt.Printf("Cleaned up %d stalled tasks\n", stalledCount)
	}

	return nil
}

// Monitor 启动队列监控
func (m *QueueManager) Monitor(ctx context.Context) {
	go m.monitorLoop(ctx)
}

// monitorLoop 监控循环
func (m *QueueManager) monitorLoop(ctx context.Context) {
	cleanupTicker := time.NewTicker(5 * time.Minute)
	statsTicker := time.NewTicker(1 * time.Minute)

	defer cleanupTicker.Stop()
	defer statsTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-cleanupTicker.C:
			// 定期清理停滞任务
			if err := m.CleanupStalledTasks(ctx); err != nil {
				fmt.Printf("Cleanup stalled tasks failed: %v\n", err)
			}
		case <-statsTicker.C:
			// 定期打印统计信息
			stats, err := m.GetQueueStats(ctx)
			if err != nil {
				fmt.Printf("Get queue stats failed: %v\n", err)
				continue
			}

			for queueName, queueStats := range stats {
				fmt.Printf("Queue %s: pending=%d, processing=%d, delayed=%d, dead=%d\n",
					queueName, queueStats.Pending, queueStats.Processing,
					queueStats.Delayed, queueStats.DeadLetter)
			}
		}
	}
}
