package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Coordinator 分布式协调器
type Coordinator struct {
	client       *redis.Client
	instanceID   string
	queues       map[string]*RedisQueue
	lockManager  *LockManager
	heartbeatKey string
	alive        bool
	mutex        sync.RWMutex
}

// NewCoordinator 创建分布式协调器
func NewCoordinator(client *redis.Client) *Coordinator {
	instanceID := fmt.Sprintf("instance-%d-%s", time.Now().UnixNano(), randomString(6))

	return &Coordinator{
		client:       client,
		instanceID:   instanceID,
		queues:       make(map[string]*RedisQueue),
		lockManager:  NewLockManager(client),
		heartbeatKey: "ZAG:COORDINATOR:HEARTBEAT",
		alive:        false,
	}
}

// Start 启动协调器
func (c *Coordinator) Start(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.alive {
		return errors.New("coordinator is already running")
	}

	c.alive = true

	// 启动心跳
	go c.heartbeatLoop(ctx)

	// 启动实例发现
	go c.discoveryLoop(ctx)

	// 启动负载均衡
	go c.loadBalanceLoop(ctx)

	fmt.Printf("Coordinator started with instance ID: %s\n", c.instanceID)
	return nil
}

// Stop 停止协调器
func (c *Coordinator) Stop(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.alive {
		return nil
	}

	c.alive = false

	// 释放所有锁
	if err := c.lockManager.ReleaseAllLocks(ctx); err != nil {
		return fmt.Errorf("release all locks failed: %w", err)
	}

	// 移除心跳
	c.client.HDel(ctx, c.heartbeatKey, c.instanceID)

	fmt.Println("Coordinator stopped")
	return nil
}

// RegisterQueue 注册队列到协调器
func (c *Coordinator) RegisterQueue(queueName string, queue *RedisQueue) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.queues[queueName] = queue
}

// AcquireQueueLock 获取队列锁（用于确保同一队列只有一个实例在处理）
func (c *Coordinator) AcquireQueueLock(ctx context.Context, queueName string) (*RedisLock, error) {
	lockName := fmt.Sprintf("queue:%s", queueName)
	return c.lockManager.AcquireLock(ctx, lockName, 30*time.Second)
}

// IsQueueLeader 检查当前实例是否是队列的领导者
func (c *Coordinator) IsQueueLeader(ctx context.Context, queueName string) (bool, error) {
	leaderKey := fmt.Sprintf("ZAG:COORDINATOR:LEADER:%s", queueName)

	// 获取当前领导者
	currentLeader, err := c.client.Get(ctx, leaderKey).Result()
	if err != nil && err != redis.Nil {
		return false, err
	}

	// 如果没有领导者或领导者是自己，尝试成为领导者
	if currentLeader == "" || currentLeader == c.instanceID {
		// 尝试设置领导者
		success, err := c.client.SetNX(ctx, leaderKey, c.instanceID, 60*time.Second).Result()
		if err != nil {
			return false, err
		}
		return success, nil
	}

	return currentLeader == c.instanceID, nil
}

// heartbeatLoop 心跳循环
func (c *Coordinator) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.sendHeartbeat(ctx)
		}
	}
}

// sendHeartbeat 发送心跳
func (c *Coordinator) sendHeartbeat(ctx context.Context) {
	c.mutex.RLock()
	if !c.alive {
		c.mutex.RUnlock()
		return
	}
	c.mutex.RUnlock()

	// 更新心跳时间戳
	timestamp := time.Now().Unix()
	err := c.client.HSet(ctx, c.heartbeatKey, c.instanceID, timestamp).Err()
	if err != nil {
		fmt.Printf("Send heartbeat failed: %v\n", err)
		return
	}

	// 设置过期时间（自动清理失效实例）
	c.client.Expire(ctx, c.heartbeatKey, 30*time.Second)
}

// discoveryLoop 实例发现循环
func (c *Coordinator) discoveryLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.discoverInstances(ctx)
		}
	}
}

// discoverInstances 发现其他实例
func (c *Coordinator) discoverInstances(ctx context.Context) {
	instances, err := c.client.HGetAll(ctx, c.heartbeatKey).Result()
	if err != nil {
		fmt.Printf("Discover instances failed: %v\n", err)
		return
	}

	now := time.Now().Unix()
	aliveInstances := make([]string, 0, len(instances))

	for instanceID, timestampStr := range instances {
		var timestamp int64
		fmt.Sscanf(timestampStr, "%d", &timestamp)

		// 检查实例是否存活（30秒内有心跳）
		if now-timestamp < 30 {
			aliveInstances = append(aliveInstances, instanceID)
		} else {
			// 移除失效实例
			c.client.HDel(ctx, c.heartbeatKey, instanceID)
		}
	}

	if len(aliveInstances) > 0 {
		fmt.Printf("Discovered %d alive instances: %v\n", len(aliveInstances), aliveInstances)
	}
}

// loadBalanceLoop 负载均衡循环
func (c *Coordinator) loadBalanceLoop(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.balanceQueues(ctx)
		}
	}
}

// balanceQueues 平衡队列负载
func (c *Coordinator) balanceQueues(ctx context.Context) {
	c.mutex.RLock()
	queues := make([]string, 0, len(c.queues))
	for queueName := range c.queues {
		queues = append(queues, queueName)
	}
	c.mutex.RUnlock()

	// 获取所有存活实例
	instances, err := c.client.HGetAll(ctx, c.heartbeatKey).Result()
	if err != nil {
		return
	}

	if len(instances) <= 1 {
		// 只有一个实例，不需要负载均衡
		return
	}

	// 对每个队列进行负载均衡
	for _, queueName := range queues {
		c.balanceQueue(ctx, queueName, instances)
	}
}

// balanceQueue 平衡单个队列的负载
func (c *Coordinator) balanceQueue(ctx context.Context, queueName string, instances map[string]string) {
	// 获取队列统计信息
	queue, exists := c.queues[queueName]
	if !exists {
		return
	}

	stats, err := queue.GetStats(ctx)
	if err != nil {
		return
	}

	// 如果队列负载过高，考虑重新分配
	totalLoad := stats.Pending + stats.Processing + stats.Delayed
	if totalLoad < 100 {
		// 负载不高，不需要重新分配
		return
	}

	// 检查当前实例是否是队列领导者
	isLeader, err := c.IsQueueLeader(ctx, queueName)
	if err != nil || !isLeader {
		return
	}

	fmt.Printf("Queue %s needs load balancing: total load = %d\n", queueName, totalLoad)

	// 这里可以实现更复杂的负载均衡算法
	// 例如：将部分任务转移到其他实例
}

// DistributedQueueManager 分布式队列管理器
type DistributedQueueManager struct {
	coordinator *Coordinator
	manager     *QueueManager
	config      DistributedConfig
}

// DistributedConfig 分布式配置
type DistributedConfig struct {
	// 是否启用分布式协调
	Enabled bool

	// 实例自动发现
	AutoDiscovery bool

	// 负载均衡
	LoadBalancing bool

	// 故障转移
	Failover bool
}

// NewDistributedQueueManager 创建分布式队列管理器
func NewDistributedQueueManager(client *redis.Client, config DistributedConfig) *DistributedQueueManager {
	coordinator := NewCoordinator(client)
	manager := NewQueueManager(client)

	return &DistributedQueueManager{
		coordinator: coordinator,
		manager:     manager,
		config:      config,
	}
}

// Start 启动分布式队列管理器
func (d *DistributedQueueManager) Start(ctx context.Context) error {
	if d.config.Enabled {
		if err := d.coordinator.Start(ctx); err != nil {
			return fmt.Errorf("start coordinator failed: %w", err)
		}
	}

	// 启动队列管理器监控
	d.manager.Monitor(ctx)

	return nil
}

// Stop 停止分布式队列管理器
func (d *DistributedQueueManager) Stop(ctx context.Context) error {
	if d.config.Enabled {
		if err := d.coordinator.Stop(ctx); err != nil {
			return fmt.Errorf("stop coordinator failed: %w", err)
		}
	}

	// 停止所有工作线程池
	if err := d.manager.StopAllWorkerPools(); err != nil {
		return fmt.Errorf("stop worker pools failed: %w", err)
	}

	return nil
}

// GetQueueManager 获取队列管理器
func (d *DistributedQueueManager) GetQueueManager() *QueueManager {
	return d.manager
}

// GetCoordinator 获取协调器
func (d *DistributedQueueManager) GetCoordinator() *Coordinator {
	return d.coordinator
}

// EnqueueWithCoordination 带协调的入队
func (d *DistributedQueueManager) EnqueueWithCoordination(ctx context.Context, queueName string, task *Task) error {
	if d.config.Enabled {
		// 获取队列锁确保一致性
		lock, err := d.coordinator.AcquireQueueLock(ctx, queueName)
		if err != nil {
			return fmt.Errorf("acquire queue lock failed: %w", err)
		}
		defer lock.Unlock(ctx)
	}

	return d.manager.EnqueueTask(ctx, queueName, task)
}

// GetClusterStats 获取集群统计信息
func (d *DistributedQueueManager) GetClusterStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取队列统计
	queueStats, err := d.manager.GetQueueStats(ctx)
	if err != nil {
		return nil, err
	}
	stats["queues"] = queueStats

	// 获取工作线程统计
	workerStats := d.manager.GetWorkerPoolStats()
	stats["workers"] = workerStats

	if d.config.Enabled {
		// 获取实例信息
		instances, err := d.coordinator.client.HGetAll(ctx, d.coordinator.heartbeatKey).Result()
		if err != nil {
			return nil, err
		}
		stats["instances"] = instances

		// 获取当前实例ID
		stats["instance_id"] = d.coordinator.instanceID
	}

	return stats, nil
}
