package workqueue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
	"zero-backend/internal/logger"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	rdbReadyQueueKey = "ready_queue"   // 就绪队列键名
	rdbTempQueueKey  = "temp_queue_%d" // 临时队列名称键名
	rdbDeadQueueKey  = "dead_queue"    // 死信队列键名
	rdbDelayQueueKey = "delay_queue"   // 延迟队列键名
	rdbLockKey       = "lock:task:%s"  // 分布式锁键名
)

var (
	ErrQueueAlreadyExists = errors.New("queue already exists")
	ErrQueueNotFound      = errors.New("queue not found")
	ErrWorkerStopped      = errors.New("worker is stopped")
)

// HandlerFunc 处理函数
type HandlerFunc func(*Task) (time.Duration, error)

// Task 任务
type Task struct {
	ID             string    `json:"id"`              // 任务唯一ID
	Payload        any       `json:"payload"`         // 原始数据
	CreatedAt      time.Time `json:"created_at"`      // 创建时间
	ExecuteAt      time.Time `json:"execute_at"`      // 执行时间
	ProcessedCount int       `json:"processed_count"` // 处理次数
}

// Worker 队列管理器
type Worker struct {
	mu      *sync.RWMutex
	wg      *sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	queues  map[string]*Queue
	stopped bool

	domain             string
	rdb                *redis.Client
	redisTimeout       time.Duration
	delayScanInterval  time.Duration
	backOffMinDuration time.Duration
	backOffMaxDuration time.Duration
	backOffFactor      float64
}

// Option 配置项
type Option func(*Worker)

// WithDomain 设置队列名称前缀
func WithDomain(domain string) Option {
	return func(w *Worker) {
		w.domain = domain
	}
}

// WithRedisTimeout 设置Redis操作超时时间
func WithRedisTimeout(timeout time.Duration) Option {
	return func(w *Worker) {
		w.redisTimeout = timeout
	}
}

// WithDelayScanInterval 设置延迟队列扫描间隔
func WithDelayScanInterval(interval time.Duration) Option {
	return func(w *Worker) {
		w.delayScanInterval = interval
	}
}

// WithBackOffConfig 设置退避策略配置
func WithBackOffConfig(min, max time.Duration, factor float64) Option {
	return func(w *Worker) {
		w.backOffMinDuration = min
		w.backOffMaxDuration = max
		w.backOffFactor = factor
	}
}

// NewWorker 创建队列管理器
func NewWorker(ctx context.Context, rdb *redis.Client, opts ...Option) *Worker {
	ctx, cancel := context.WithCancel(ctx)
	worker := &Worker{
		mu:                 &sync.RWMutex{},
		wg:                 &sync.WaitGroup{},
		ctx:                ctx,
		cancel:             cancel,
		queues:             make(map[string]*Queue),
		stopped:            false,
		domain:             "workqueue",
		rdb:                rdb,
		redisTimeout:       3 * time.Second,
		delayScanInterval:  5 * time.Second,
		backOffMinDuration: 100 * time.Millisecond,
		backOffMaxDuration: 5 * time.Second,
		backOffFactor:      2,
	}

	for _, opt := range opts {
		opt(worker)
	}

	return worker
}

// RegisterQueue 注册队列
func (w *Worker) RegisterQueue(name string, concurrency int, maxRetries int, handler HandlerFunc) error {
	if concurrency <= 0 {
		return errors.New("concurrency must be positive")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.stopped {
		return ErrWorkerStopped
	}
	if _, exists := w.queues[name]; exists {
		return ErrQueueAlreadyExists
	}

	ctx, cancel := context.WithCancel(w.ctx)
	q := &Queue{
		w:           w,
		mu:          &sync.RWMutex{},
		wg:          &sync.WaitGroup{},
		ctx:         ctx,
		cancel:      cancel,
		name:        name,
		concurrency: concurrency,
		maxRetries:  maxRetries,
		handler:     handler,
	}
	w.queues[name] = q

	w.wg.Add(1)
	go q.startProcessing()

	return nil
}

// UnregisterQueue 注销队列
func (w *Worker) UnregisterQueue(name string) error {
	w.mu.RLock()
	q, exists := w.queues[name]
	w.mu.RUnlock()

	if !exists {
		return ErrQueueNotFound
	}

	q.stop()
	return nil
}

// Shutdown 停止整个Worker
func (w *Worker) Shutdown(ctx context.Context) error {
	w.mu.Lock()
	if w.stopped {
		w.mu.Unlock()
		return nil
	}

	w.stopped = true
	w.cancel()
	w.mu.Unlock()

	// 等待队列全部处理完成
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

// EnqueueTask 添加任务
func (w *Worker) EnqueueTask(name string, payload any, delay time.Duration) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	q, exists := w.queues[name]
	if !exists {
		return ErrQueueNotFound
	}

	t := &Task{
		ID:             uuid.New().String(),
		Payload:        payload,
		CreatedAt:      time.Now(),
		ExecuteAt:      time.Now().Add(delay),
		ProcessedCount: 0,
	}

	return q.enqueueTask(t)
}

// Queue 任务队列
type Queue struct {
	w       *Worker
	mu      *sync.RWMutex
	wg      *sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	stopped bool

	name        string
	concurrency int
	maxRetries  int
	handler     HandlerFunc
}

// startProcessing 启动队列处理
func (q *Queue) startProcessing() {
	q.mu.Lock()
	if q.stopped {
		q.mu.Unlock()
		q.w.wg.Done()
		return
	}

	q.mu.Unlock()

	// 启动并发处理
	q.wg.Add(q.concurrency + 1)
	go q.processDelayQueue()

	for i := 0; i < q.concurrency; i++ {
		go q.workerProcess(i)
	}

	// 等待停止信号
	<-q.ctx.Done()
	q.stop()
}

// processDelayQueue 处理延迟队列
func (q *Queue) processDelayQueue() {
	defer q.wg.Done()

	ticker := time.NewTicker(q.w.delayScanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-q.ctx.Done():
			return
		case <-ticker.C:
			q.processDueTasks()
		}
	}
}

// processDueTasks 处理到期的任务
func (q *Queue) processDueTasks() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	ctx, cancel = context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	// 获取当前时间作为最大score
	maxScore := float64(time.Now().Unix())

	// 获取所有已到期的任务
	tasks, err := q.w.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:     q.redisKey(rdbDelayQueueKey),
		Start:   "0",
		Stop:    fmt.Sprintf("%f", maxScore),
		ByScore: true,
	}).Result()
	if err != nil {
		logger.Ctx(q.ctx).Err(err, "fetch delay tasks error")
		return
	}

	if len(tasks) == 0 {
		return
	}

	// 移动任务到就绪队列
	ctx, cancel = context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	if _, err := q.w.rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.ZRem(ctx, q.redisKey(rdbDelayQueueKey), tasks)
		pipe.LPush(ctx, q.redisKey(rdbReadyQueueKey), tasks)
		return nil
	}); err != nil {
		logger.Ctx(q.ctx).Err(err, "move tasks to ready queue error")
		return
	}

	logger.Ctx(q.ctx).Info("moved tasks to ready queue", "count", len(tasks))
}

// workerProcess 处理任务
func (q *Queue) workerProcess(id int) {
	defer q.wg.Done()

	// 指数退避策略
	backoff := NewExponentialBackoff(q.w.backOffMinDuration, q.w.backOffMaxDuration, q.w.backOffFactor)
	tempQueueKey := q.redisKey(fmt.Sprintf(rdbTempQueueKey, id))

	for {
		select {
		case <-q.ctx.Done():
			return
		default:
			// 从队列中取出任务
			taskStr, err := q.fetchTask(tempQueueKey)
			if err != nil {
				if errors.Is(err, redis.Nil) {
					backoff.Wait()

					// 队列中暂无任务，回收临时队列任务
					q.reclaimTempQueueTasks(tempQueueKey)
					continue
				}

				logger.Ctx(q.ctx).Err(err, "fetch task error")
				backoff.Wait()
				continue
			}

			backoff.Reset()
			q.processTask(taskStr, tempQueueKey)
		}
	}
}

// fetchTask 从队列中取出任务
func (q *Queue) fetchTask(tempQueueKey string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	// 将任务从就绪队列中取出移入临时队列中
	taskStr, err := q.w.rdb.LMove(ctx, q.redisKey(rdbReadyQueueKey), tempQueueKey, "RIGHT", "LEFT").Result()
	if err != nil {
		return "", err
	}

	return taskStr, nil
}

// processTask 处理任务
func (q *Queue) processTask(taskStr string, tempQueueKey string) {
	// 解析任务
	task, err := q.parseTask(taskStr)
	if err != nil || task == nil {
		// 解析失败，删除任务
		q.removeTask(taskStr, tempQueueKey)
		logger.Ctx(q.ctx).Err(err, "parse task error, task removed")
		return
	}

	// 尝试获取分布式锁，防止重复处理
	lockKey := q.redisKey(fmt.Sprintf(rdbLockKey, task.ID))
	ctx, cancel := context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	locked, err := q.w.rdb.SetNX(ctx, lockKey, "1", 30*time.Second).Result()
	if err != nil {
		logger.Ctx(q.ctx).Err(err, "acquire lock error")
		// 获取锁失败，重新入队
		q.requeueTask(tempQueueKey)
		return
	}
	if !locked {
		// 已被其他消费者处理，删除任务
		q.removeTask(taskStr, tempQueueKey)
		logger.Ctx(q.ctx).Info("task already processed by another worker", "task_id", task.ID)
		return
	}

	// 执行处理任务函数
	if delay, err := q.handler(task); err != nil {
		task.ProcessedCount++
		task.ExecuteAt = time.Now().Add(delay)

		// 删除临时队列中的任务
		q.removeTask(taskStr, tempQueueKey)

		if task.ProcessedCount >= q.maxRetries {
			// 达到最大处理次数，添加到死信队列
			if err := q.enqueueDeadTask(task); err != nil {
				logger.Ctx(q.ctx).Err(err, "enqueue dead task error")
				// 死信队列写入失败，记录日志但不丢失任务
			}
			logger.Ctx(q.ctx).Info("task moved to dead queue",
				"task_id", task.ID,
				"processed_count", task.ProcessedCount,
			)
		} else {
			// 任务重新写入队列重新处理
			if err := q.enqueueTask(task); err != nil {
				logger.Ctx(q.ctx).Err(err, "enqueue task error")
				// 入队失败，任务可能丢失，记录错误
			}
			logger.Ctx(q.ctx).Info("task requeued for retry",
				"task_id", task.ID,
				"processed_count", task.ProcessedCount,
				"next_retry_at", task.ExecuteAt,
			)
		}
	} else {
		// 任务处理成功，删除任务
		q.removeTask(taskStr, tempQueueKey)
		logger.Ctx(q.ctx).Info("task processed successfully", "task_id", task.ID)
	}

	// 释放锁
	q.w.rdb.Del(ctx, lockKey)
}

// requeueTask 重新入队任务
func (q *Queue) requeueTask(tempQueueKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	_, err := q.w.rdb.LMove(ctx, tempQueueKey, q.redisKey(rdbReadyQueueKey), "RIGHT", "LEFT").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Ctx(q.ctx).Err(err, "requeue task error")
	}
}

// removeTask 从队列删除任务
func (q *Queue) removeTask(task string, queueName string) {
	ctx, cancel := context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	_, err := q.w.rdb.LRem(ctx, queueName, 0, task).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return
		}

		logger.Ctx(q.ctx).Err(err, "remove task error")
		return
	}
}

// enqueueDeadTask 将任务添加到死信队列
func (q *Queue) enqueueDeadTask(task *Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	taskStr, err := q.serializeTask(task)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	if err := q.w.rdb.LPush(ctx, q.redisKey(rdbDeadQueueKey), taskStr).Err(); err != nil {
		return err
	}

	return nil
}

// reclaimTempQueueTasks 回收临时队列任务
func (q *Queue) reclaimTempQueueTasks(tempQueueKey string) {
	ctx, cancel := context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	count := 0
	for {
		_, err := q.w.rdb.LMove(ctx, tempQueueKey, q.redisKey(rdbReadyQueueKey), "RIGHT", "LEFT").Result()
		if err != nil {
			if !errors.Is(err, redis.Nil) {
				logger.Ctx(q.ctx).Err(err, "reclaim temp queue task error")
			}
			break
		}
		count++
	}

	if count > 0 {
		logger.Ctx(q.ctx).Info("reclaimed tasks from temp queue", "count", count)
	}
}

// parseTask 解析任务
func (q *Queue) parseTask(data string) (*Task, error) {
	task := &Task{}

	if err := json.Unmarshal([]byte(data), task); err != nil {
		return nil, fmt.Errorf("parse task error: %w", err)
	}

	return task, nil
}

// serializeTask 序列化任务
func (q *Queue) serializeTask(task *Task) (string, error) {
	if task == nil {
		return "", errors.New("task cannot be nil")
	}

	taskBytes, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("serialize task error: %w", err)
	}

	return string(taskBytes), nil
}

// stop 停止队列
func (q *Queue) stop() {
	q.mu.Lock()
	if q.stopped {
		q.mu.Unlock()
		return
	}
	q.stopped = true
	q.cancel()
	q.mu.Unlock()

	// 等待所有协程处理完成
	done := make(chan struct{})
	go func() {
		q.wg.Wait()
		close(done)
	}()

	<-done

	// 队列停止，队列数量减1
	q.w.wg.Done()
}

// enqueueTask 添加任务到队列
func (q *Queue) enqueueTask(task *Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	// 序列化任务
	taskStr, err := q.serializeTask(task)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), q.w.redisTimeout)
	defer cancel()

	if task.ExecuteAt.After(time.Now().Add(time.Second)) {
		// 写入延迟队列
		if err := q.w.rdb.ZAdd(ctx, q.redisKey(rdbDelayQueueKey), redis.Z{
			Score:  float64(task.ExecuteAt.Unix()),
			Member: taskStr,
		}).Err(); err != nil {
			return fmt.Errorf("failed to add to delay queue: %w", err)
		}
		logger.Ctx(q.ctx).Info("task enqueued to delay queue",
			"task_id", task.ID,
			"execute_at", task.ExecuteAt,
		)
	} else {
		// 写入就绪队列
		if err := q.w.rdb.LPush(ctx, q.redisKey(rdbReadyQueueKey), taskStr).Err(); err != nil {
			return fmt.Errorf("failed to add to ready queue: %w", err)
		}
		logger.Ctx(q.ctx).Info("task enqueued to ready queue", "task_id", task.ID)
	}

	return nil
}

// redisKey 获取redis键名
func (q *Queue) redisKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", q.w.domain, q.name, key)
}

// ExponentialBackoff 指数退避工具
type ExponentialBackoff struct {
	min    time.Duration
	max    time.Duration
	factor float64
	next   time.Duration
}

// NewExponentialBackoff 创建指数退避工具
func NewExponentialBackoff(min, max time.Duration, factor float64) *ExponentialBackoff {
	return &ExponentialBackoff{
		min:    min,
		max:    max,
		factor: factor,
		next:   min,
	}
}

// Wait 阻塞等待
func (b *ExponentialBackoff) Wait() {
	time.Sleep(b.next)
	b.next = time.Duration(float64(b.next) * b.factor)
	if b.next > b.max {
		b.next = b.max
	}
}

// Reset 重置
func (b *ExponentialBackoff) Reset() {
	b.next = b.min
}
