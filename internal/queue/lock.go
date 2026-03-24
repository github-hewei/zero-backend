package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisLock 基于Redis的分布式锁
type RedisLock struct {
	client    *redis.Client
	key       string
	token     string
	timeout   time.Duration
	renewTick *time.Ticker
	renewStop chan struct{}
}

// NewRedisLock 创建新的Redis分布式锁
func NewRedisLock(client *redis.Client, key string, timeout time.Duration) *RedisLock {
	return &RedisLock{
		client:  client,
		key:     key,
		timeout: timeout,
		token:   generateLockToken(),
	}
}

// TryLock 尝试获取锁
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
	// 使用SET NX EX命令原子性地获取锁
	result, err := l.client.SetNX(ctx, l.key, l.token, l.timeout).Result()
	if err != nil {
		return false, fmt.Errorf("try lock failed: %w", err)
	}

	if result {
		// 启动锁续期
		l.startRenewal(ctx)
	}

	return result, nil
}

// Lock 阻塞获取锁
func (l *RedisLock) Lock(ctx context.Context) error {
	timeout := time.NewTimer(l.timeout)
	defer timeout.Stop()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout.C:
			return errors.New("lock timeout")
		case <-ticker.C:
			acquired, err := l.TryLock(ctx)
			if err != nil {
				return err
			}
			if acquired {
				return nil
			}
		}
	}
}

// Unlock 释放锁
func (l *RedisLock) Unlock(ctx context.Context) error {
	// 停止锁续期
	l.stopRenewal()

	// 使用Lua脚本确保只有锁的持有者才能释放锁
	script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`

	result, err := l.client.Eval(ctx, script, []string{l.key}, l.token).Result()
	if err != nil {
		return fmt.Errorf("unlock failed: %w", err)
	}

	if result.(int64) == 0 {
		return errors.New("lock token mismatch or lock already released")
	}

	return nil
}

// IsLocked 检查锁是否被持有
func (l *RedisLock) IsLocked(ctx context.Context) (bool, error) {
	val, err := l.client.Get(ctx, l.key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return val == l.token, nil
}

// startRenewal 启动锁续期
func (l *RedisLock) startRenewal(ctx context.Context) {
	l.renewTick = time.NewTicker(l.timeout / 2)
	l.renewStop = make(chan struct{})

	go func() {
		for {
			select {
			case <-l.renewTick.C:
				l.renew(ctx)
			case <-l.renewStop:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// stopRenewal 停止锁续期
func (l *RedisLock) stopRenewal() {
	if l.renewTick != nil {
		l.renewTick.Stop()
	}
	if l.renewStop != nil {
		close(l.renewStop)
	}
}

// renew 续期锁
func (l *RedisLock) renew(ctx context.Context) {
	// 使用Lua脚本续期锁
	script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("expire", KEYS[1], ARGV[2])
else
    return 0
end
`

	expireSeconds := int(l.timeout / time.Second)
	_, err := l.client.Eval(ctx, script, []string{l.key}, l.token, expireSeconds).Result()
	if err != nil {
		fmt.Printf("Renew lock failed: %v\n", err)
	}
}

// generateLockToken 生成锁令牌
func generateLockToken() string {
	return fmt.Sprintf("lock-%d-%s", time.Now().UnixNano(), randomString(8))
}

// LockManager 锁管理器
type LockManager struct {
	client *redis.Client
	locks  map[string]*RedisLock
	mutex  sync.RWMutex
}

// NewLockManager 创建锁管理器
func NewLockManager(client *redis.Client) *LockManager {
	return &LockManager{
		client: client,
		locks:  make(map[string]*RedisLock),
	}
}

// AcquireLock 获取锁
func (m *LockManager) AcquireLock(ctx context.Context, lockName string, timeout time.Duration) (*RedisLock, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 检查是否已经持有该锁
	if lock, exists := m.locks[lockName]; exists {
		locked, err := lock.IsLocked(ctx)
		if err != nil {
			return nil, err
		}
		if locked {
			return lock, nil
		}
	}

	// 创建新锁
	lock := NewRedisLock(m.client, "ZAG:LOCK:"+lockName, timeout)

	// 尝试获取锁
	acquired, err := lock.TryLock(ctx)
	if err != nil {
		return nil, err
	}

	if !acquired {
		// 阻塞获取锁
		if err := lock.Lock(ctx); err != nil {
			return nil, err
		}
	}

	// 保存锁引用
	m.locks[lockName] = lock

	return lock, nil
}

// ReleaseLock 释放锁
func (m *LockManager) ReleaseLock(ctx context.Context, lockName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	lock, exists := m.locks[lockName]
	if !exists {
		return errors.New("lock not found")
	}

	if err := lock.Unlock(ctx); err != nil {
		return err
	}

	delete(m.locks, lockName)
	return nil
}

// ReleaseAllLocks 释放所有锁
func (m *LockManager) ReleaseAllLocks(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var errs []error
	for name, lock := range m.locks {
		if err := lock.Unlock(ctx); err != nil {
			errs = append(errs, fmt.Errorf("release lock %s failed: %w", name, err))
		}
	}

	m.locks = make(map[string]*RedisLock)

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// WithLock 使用锁执行函数
func WithLock(ctx context.Context, lockManager *LockManager, lockName string, timeout time.Duration, fn func() error) error {
	_, err := lockManager.AcquireLock(ctx, lockName, timeout)
	if err != nil {
		return fmt.Errorf("acquire lock failed: %w", err)
	}

	defer func() {
		if err := lockManager.ReleaseLock(ctx, lockName); err != nil {
			fmt.Printf("Release lock failed: %v\n", err)
		}
	}()

	return fn()
}
