package queue

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRedisClient 创建模拟Redis客户端
func mockRedisClient(t *testing.T) *redis.Client {
	// 使用不同的DB避免测试间干扰
	db := time.Now().UnixNano() % 15
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.91.100:6379",
		Password: "123456",
		DB:       int(db),
	})

	// 测试连接
	ctx := context.Background()
	err := client.Ping(ctx).Err()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
	}

	// 清理数据库
	err = client.FlushDB(ctx).Err()
	if err != nil {
		t.Logf("Flush DB failed: %v", err)
	}

	return client
}

// cleanupRedis 清理Redis测试数据
func cleanupRedis(t *testing.T, client *redis.Client) {
	if client == nil {
		return
	}

	ctx := context.Background()
	// 不关闭客户端，让测试函数自己处理
	// 只记录清理错误，不失败测试
	err := client.FlushDB(ctx).Err()
	if err != nil && err.Error() != "redis: client is closed" {
		t.Logf("Cleanup Redis failed: %v", err)
	}
}

func TestTaskCreation(t *testing.T) {
	task := NewTask("test-queue", "test-type", []byte("test payload"))

	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "test-queue", task.Queue)
	assert.Equal(t, "test-type", task.Type)
	assert.Equal(t, []byte("test payload"), task.Payload)
	assert.Equal(t, 3, task.MaxRetries)
	assert.Equal(t, 0, task.RetryCount)
	assert.Equal(t, TaskStatusPending, task.Status)
}

func TestTaskWithDelay(t *testing.T) {
	task := NewTask("test-queue", "test-type", []byte("payload"))

	// 设置5秒延迟
	task.WithDelay(5 * time.Second)

	now := time.Now().Unix()
	expectedDelayUntil := now + 5

	// 允许1秒的误差
	assert.InDelta(t, expectedDelayUntil, task.DelayUntil, 1)
}

func TestTaskWithMetadata(t *testing.T) {
	task := NewTask("test-queue", "test-type", []byte("payload"))

	task.WithMetadata("key1", "value1")
	task.WithMetadata("key2", "value2")

	assert.Equal(t, "value1", task.Metadata["key1"])
	assert.Equal(t, "value2", task.Metadata["key2"])
}

func TestTaskMarshaling(t *testing.T) {
	task := NewTask("test-queue", "test-type", []byte("test payload"))
	task.WithMetadata("test-key", "test-value")

	// 序列化
	data, err := task.Marshal()
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// 反序列化
	decodedTask, err := UnmarshalTask(data)
	require.NoError(t, err)

	assert.Equal(t, task.ID, decodedTask.ID)
	assert.Equal(t, task.Queue, decodedTask.Queue)
	assert.Equal(t, task.Type, decodedTask.Type)
	assert.Equal(t, task.Payload, decodedTask.Payload)
	assert.Equal(t, task.Metadata["test-key"], decodedTask.Metadata["test-key"])
}

func TestQueueConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "default", config.Name)
	assert.Equal(t, 10, config.MaxConcurrency)
	assert.Equal(t, 300, config.VisibilityTimeout)
	assert.Equal(t, 3, config.MaxRetries)
	assert.Equal(t, RetryDelayExponential, config.RetryDelay)
	assert.True(t, config.EnableDeadLetter)

	// 测试配置修改
	modifiedConfig := config.
		WithName("test-queue").
		WithMaxConcurrency(5).
		WithMaxRetries(5)

	assert.Equal(t, "test-queue", modifiedConfig.Name)
	assert.Equal(t, 5, modifiedConfig.MaxConcurrency)
	assert.Equal(t, 5, modifiedConfig.MaxRetries)
}

func TestRedisQueueBasicOperations(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	config := DefaultConfig().WithName("test-queue")
	queue := NewRedisQueue(client, config)

	// 测试入队
	task := NewTask("test-queue", "test-type", []byte("test payload"))
	err := queue.Enqueue(ctx, task)
	require.NoError(t, err)

	// 测试出队
	dequeuedTask, err := queue.Dequeue(ctx)
	require.NoError(t, err)
	require.NotNil(t, dequeuedTask)

	assert.Equal(t, task.ID, dequeuedTask.ID)
	assert.Equal(t, TaskStatusProcessing, dequeuedTask.Status)

	// 测试确认完成
	err = queue.Ack(ctx, dequeuedTask)
	require.NoError(t, err)

	// 验证任务已从处理中队列移除
	stats, err := queue.GetStats(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.Processing)
	assert.Equal(t, int64(1), stats.Completed)
}

func TestRedisQueueDelayedTask(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	config := DefaultConfig().WithName("test-delayed-queue")
	queue := NewRedisQueue(client, config)

	// 入队延迟任务
	task := NewTask("test-delayed-queue", "delayed-type", []byte("delayed payload"))
	err := queue.EnqueueWithDelay(ctx, task, 2*time.Second)
	require.NoError(t, err)

	// 立即尝试出队（应该没有任务，因为延迟未到）
	dequeuedTask, err := queue.Dequeue(ctx)
	require.NoError(t, err)
	assert.Nil(t, dequeuedTask)

	// 等待延迟时间
	time.Sleep(3 * time.Second)

	// 现在应该可以出队
	dequeuedTask, err = queue.Dequeue(ctx)
	require.NoError(t, err)
	require.NotNil(t, dequeuedTask)

	assert.Equal(t, task.ID, dequeuedTask.ID)

	// 确认任务
	err = queue.Ack(ctx, dequeuedTask)
	require.NoError(t, err)
}

func TestRedisQueueRetry(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	config := DefaultConfig().WithName("test-retry-queue")
	queue := NewRedisQueue(client, config)

	// 入队任务
	task := NewTask("test-retry-queue", "retry-type", []byte("retry payload"))
	task.WithMaxRetries(2)

	err := queue.Enqueue(ctx, task)
	require.NoError(t, err)

	// 出队任务
	dequeuedTask, err := queue.Dequeue(ctx)
	require.NoError(t, err)
	require.NotNil(t, dequeuedTask)

	// 模拟任务失败
	err = queue.Nack(ctx, dequeuedTask, assert.AnError)
	require.NoError(t, err)

	// 检查任务是否被重试
	stats, err := queue.GetStats(ctx)
	require.NoError(t, err)

	// 任务应该在延迟队列中等待重试
	assert.Equal(t, int64(1), stats.Delayed)
}

func TestRedisQueueDeadLetter(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	config := DefaultConfig().WithName("test-deadletter-queue").WithDeadLetter(true, 1)
	queue := NewRedisQueue(client, config)

	// 创建只允许重试1次的任务
	task := NewTask("test-deadletter-queue", "deadletter-type", []byte("deadletter payload"))
	task.WithMaxRetries(1)

	err := queue.Enqueue(ctx, task)
	require.NoError(t, err)

	// 出队并失败两次（超过最大重试次数）
	dequeuedTask, err := queue.Dequeue(ctx)
	require.NoError(t, err)
	require.NotNil(t, dequeuedTask)

	// 第一次失败
	err = queue.Nack(ctx, dequeuedTask, assert.AnError)
	require.NoError(t, err)

	// 等待重试（指数退避：第一次重试延迟2秒）
	time.Sleep(3 * time.Second)

	// 再次出队（重试）
	dequeuedTask, err = queue.Dequeue(ctx)
	require.NoError(t, err)
	require.NotNil(t, dequeuedTask)

	// 第二次失败（应该进入死信队列）
	err = queue.Nack(ctx, dequeuedTask, assert.AnError)
	require.NoError(t, err)

	// 检查死信队列
	stats, err := queue.GetStats(ctx)
	require.NoError(t, err)

	assert.Equal(t, int64(1), stats.DeadLetter)
}

func TestWorkerPool(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	config := DefaultConfig().WithName("test-worker-pool").WithMaxConcurrency(2)
	queue := NewRedisQueue(client, config)

	// 创建处理器
	processedTasks := 0
	handler := HandlerFunc(func(ctx context.Context, task *Task) error {
		processedTasks++
		return nil
	})

	// 创建工作线程池
	workerPool := NewWorkerPool(queue, handler, config)

	// 启动工作线程池
	err := workerPool.Start(ctx)
	require.NoError(t, err)

	// 入队一些任务
	for i := 0; i < 5; i++ {
		task := NewTask("test-worker-pool", "worker-type", []byte("worker payload"))
		err := queue.Enqueue(ctx, task)
		require.NoError(t, err)
	}

	// 等待任务被处理
	time.Sleep(2 * time.Second)

	// 检查处理的任务数量
	assert.GreaterOrEqual(t, processedTasks, 1)

	// 停止工作线程池
	err = workerPool.Stop()
	require.NoError(t, err)

	// 检查工作线程池状态
	assert.False(t, workerPool.IsRunning())
}

func TestQueueManager(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	// 创建队列管理器
	manager := NewQueueManager(client)

	// 获取或创建队列
	testQueue := manager.GetOrCreateQueue("manager-test", DefaultConfig())
	assert.NotNil(t, testQueue)

	// 入队任务
	task := NewTask("manager-test", "manager-type", []byte("manager payload"))
	err := manager.EnqueueTask(ctx, "manager-test", task)
	require.NoError(t, err)

	// 获取队列统计信息
	stats, err := manager.GetQueueStats(ctx)
	require.NoError(t, err)

	queueStats, exists := stats["manager-test"]
	require.True(t, exists)
	assert.Equal(t, int64(1), queueStats.Pending)
}

func TestDeduplicator(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	dedup := NewDeduplicator(client)
	payload := []byte("unique payload")

	// 第一次检查应该不是重复
	isDuplicate, err := dedup.IsDuplicate(ctx, "test-queue", payload)
	require.NoError(t, err)
	assert.False(t, isDuplicate)

	// 标记为已处理
	err = dedup.MarkAsProcessed(ctx, "test-queue", payload)
	require.NoError(t, err)

	// 第二次检查应该是重复
	isDuplicate, err = dedup.IsDuplicate(ctx, "test-queue", payload)
	require.NoError(t, err)
	assert.True(t, isDuplicate)
}

func TestLockManager(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	lockManager := NewLockManager(client)

	// 获取锁
	lock, err := lockManager.AcquireLock(ctx, "test-lock", 10*time.Second)
	require.NoError(t, err)
	require.NotNil(t, lock)

	// 检查锁是否被持有
	isLocked, err := lock.IsLocked(ctx)
	require.NoError(t, err)
	assert.True(t, isLocked)

	// 释放锁
	err = lockManager.ReleaseLock(ctx, "test-lock")
	require.NoError(t, err)

	// 再次检查锁状态
	isLocked, err = lock.IsLocked(ctx)
	require.NoError(t, err)
	assert.False(t, isLocked)
}

// 集成测试：完整的队列工作流程
func TestIntegrationWorkflow(t *testing.T) {
	client := mockRedisClient(t)
	defer cleanupRedis(t, client)
	defer client.Close()

	ctx := context.Background()

	// 1. 创建队列管理器
	manager := NewQueueManager(client)

	// 2. 定义处理器
	processedCount := 0
	handler := HandlerFunc(func(ctx context.Context, task *Task) error {
		processedCount++
		t.Logf("Processing task %s: %s", task.ID, string(task.Payload))
		return nil
	})

	// 3. 注册工作线程池
	_, err := manager.RegisterWorkerPool("integration-test", handler,
		DefaultConfig().WithMaxConcurrency(3))
	require.NoError(t, err)

	// 4. 启动工作线程池
	err = manager.StartAllWorkerPools(ctx)
	require.NoError(t, err)

	// 5. 入队多个任务
	for i := 0; i < 10; i++ {
		task := NewTask("integration-test", "integration-type",
			[]byte("integration payload"))

		if i%3 == 0 {
			// 每3个任务延迟1秒
			err = manager.EnqueueTaskWithDelay(ctx, "integration-test", task, 1*time.Second)
		} else {
			err = manager.EnqueueTask(ctx, "integration-test", task)
		}
		require.NoError(t, err)
	}

	// 6. 等待任务被处理
	time.Sleep(3 * time.Second)

	// 7. 检查处理结果
	assert.GreaterOrEqual(t, processedCount, 5)

	// 8. 获取统计信息
	stats, err := manager.GetQueueStats(ctx)
	require.NoError(t, err)

	queueStats, exists := stats["integration-test"]
	require.True(t, exists)
	t.Logf("Queue stats: pending=%d, processing=%d, completed=%d, delayed=%d",
		queueStats.Pending, queueStats.Processing, queueStats.Completed, queueStats.Delayed)

	// 9. 停止所有工作线程池
	err = manager.StopAllWorkerPools()
	require.NoError(t, err)

	t.Log("Integration test completed successfully")
}
