package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// RedisWorker 基于Redis的工作线程
type RedisWorker struct {
	id      string
	queue   *RedisQueue
	handler Handler
	config  QueueConfig

	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	running    bool
	stats      *WorkerStats
	statsMutex sync.RWMutex
}

// NewRedisWorker 创建新的Redis工作线程
func NewRedisWorker(queue *RedisQueue, handler Handler, config QueueConfig) *RedisWorker {
	workerID := fmt.Sprintf("worker-%s-%d", config.Name, time.Now().UnixNano())

	return &RedisWorker{
		id:      workerID,
		queue:   queue,
		handler: handler,
		config:  config,
		stats: &WorkerStats{
			ID:              workerID,
			Running:         false,
			ProcessedTasks:  0,
			SuccessfulTasks: 0,
			FailedTasks:     0,
			StartedAt:       time.Time{},
			LastActivity:    time.Time{},
		},
	}
}

// Start 启动工作线程
func (w *RedisWorker) Start(ctx context.Context) error {
	if w.running {
		return errors.New("worker is already running")
	}

	w.ctx, w.cancel = context.WithCancel(ctx)
	w.running = true
	w.stats.StartedAt = time.Now()
	w.stats.Running = true

	// 启动工作协程
	w.wg.Add(1)
	go w.run()

	return nil
}

// Stop 停止工作线程
func (w *RedisWorker) Stop() error {
	if !w.running {
		return nil
	}

	w.cancel()
	w.wg.Wait()

	w.statsMutex.Lock()
	w.running = false
	w.stats.Running = false
	w.stats.LastActivity = time.Now()
	w.statsMutex.Unlock()

	return nil
}

// IsRunning 是否正在运行
func (w *RedisWorker) IsRunning() bool {
	w.statsMutex.RLock()
	defer w.statsMutex.RUnlock()
	return w.running
}

// GetStats 获取工作线程统计信息
func (w *RedisWorker) GetStats() *WorkerStats {
	w.statsMutex.RLock()
	defer w.statsMutex.RUnlock()

	// 返回副本
	stats := *w.stats
	return &stats
}

// run 工作线程主循环
func (w *RedisWorker) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			w.processTask()
		}
	}
}

// processTask 处理单个任务
func (w *RedisWorker) processTask() {
	ctx, cancel := context.WithTimeout(w.ctx, 30*time.Second)
	defer cancel()

	// 从队列获取任务
	task, err := w.queue.Dequeue(ctx)
	if err != nil {
		// 记录错误但不停止工作线程
		fmt.Printf("Worker %s: dequeue failed: %v\n", w.id, err)
		return
	}

	if task == nil {
		// 队列为空，等待下一次轮询
		return
	}

	// 更新统计信息
	w.statsMutex.Lock()
	w.stats.CurrentTask = task
	w.stats.LastActivity = time.Now()
	w.statsMutex.Unlock()

	// 处理任务
	startTime := time.Now()
	processErr := w.handler.Handle(ctx, task)
	processingTime := time.Since(startTime)

	// 更新统计信息
	w.statsMutex.Lock()
	w.stats.ProcessedTasks++
	w.stats.CurrentTask = nil

	if processErr == nil {
		w.stats.SuccessfulTasks++
		// 确认任务完成
		if err := w.queue.Ack(ctx, task); err != nil {
			fmt.Printf("Worker %s: ack failed: %v\n", w.id, err)
		}
	} else {
		w.stats.FailedTasks++
		// 确认任务失败
		if err := w.queue.Nack(ctx, task, processErr); err != nil {
			fmt.Printf("Worker %s: nack failed: %v\n", w.id, err)
		}
	}
	w.statsMutex.Unlock()

	// 记录处理结果
	if processErr != nil {
		fmt.Printf("Worker %s: task %s failed after %v: %v\n",
			w.id, task.ID, processingTime, processErr)
	} else {
		fmt.Printf("Worker %s: task %s completed in %v\n",
			w.id, task.ID, processingTime)
	}
}

// WorkerPool 工作线程池
type WorkerPool struct {
	queue   *RedisQueue
	handler Handler
	config  QueueConfig
	workers []*RedisWorker
	running bool
	mutex   sync.RWMutex
}

// NewWorkerPool 创建工作线程池
func NewWorkerPool(queue *RedisQueue, handler Handler, config QueueConfig) *WorkerPool {
	return &WorkerPool{
		queue:   queue,
		handler: handler,
		config:  config,
		workers: make([]*RedisWorker, 0, config.MaxConcurrency),
	}
}

// Start 启动工作线程池
func (p *WorkerPool) Start(ctx context.Context) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.running {
		return errors.New("worker pool is already running")
	}

	// 创建工作线程
	for i := 0; i < p.config.MaxConcurrency; i++ {
		worker := NewRedisWorker(p.queue, p.handler, p.config)
		if err := worker.Start(ctx); err != nil {
			// 停止已启动的工作线程
			for _, w := range p.workers {
				w.Stop()
			}
			return fmt.Errorf("start worker %d failed: %w", i, err)
		}
		p.workers = append(p.workers, worker)
	}

	p.running = true
	return nil
}

// Stop 停止工作线程池
func (p *WorkerPool) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.running {
		return nil
	}

	// 停止所有工作线程
	for _, worker := range p.workers {
		worker.Stop()
	}

	p.workers = p.workers[:0]
	p.running = false

	return nil
}

// IsRunning 是否正在运行
func (p *WorkerPool) IsRunning() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.running
}

// GetWorkerStats 获取所有工作线程的统计信息
func (p *WorkerPool) GetWorkerStats() []*WorkerStats {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	stats := make([]*WorkerStats, len(p.workers))
	for i, worker := range p.workers {
		stats[i] = worker.GetStats()
	}

	return stats
}

// Scale 调整工作线程数量
func (p *WorkerPool) Scale(ctx context.Context, newSize int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !p.running {
		return errors.New("worker pool is not running")
	}

	currentSize := len(p.workers)

	if newSize > currentSize {
		// 增加工作线程
		for i := currentSize; i < newSize; i++ {
			worker := NewRedisWorker(p.queue, p.handler, p.config)
			if err := worker.Start(ctx); err != nil {
				return fmt.Errorf("start new worker failed: %w", err)
			}
			p.workers = append(p.workers, worker)
		}
	} else if newSize < currentSize {
		// 减少工作线程
		workersToStop := p.workers[newSize:]
		p.workers = p.workers[:newSize]

		// 停止多余的工作线程
		for _, worker := range workersToStop {
			worker.Stop()
		}
	}

	// 更新配置
	p.config.MaxConcurrency = newSize

	return nil
}
