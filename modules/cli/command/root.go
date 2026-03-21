package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"zero-backend/internal/ctxkeys"
	logger2 "zero-backend/internal/logger"
	"zero-backend/pkg/locker"
	"zero-backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// RootCommand 根命令
type RootCommand struct {
	*cobra.Command
	curCmd *cobra.Command
	lock   locker.Lock
}

// NewRootCommand 创建根命令
func NewRootCommand(
	logger logger.Logger,
	locker *locker.RedisLocker,
	user *UserCommand,
	migrate *MigrateCommand,
) *RootCommand {
	cmd := &RootCommand{
		Command: &cobra.Command{
			Use:   "cli",
			Short: "CLI 工具",
			Long:  `命令行管理工具，用于执行数据迁移、用户管理等操作`,
		},
	}

	cmd.Configure(logger, locker)
	cmd.AddCommand(user.Command)
	cmd.AddCommand(migrate.Command)
	return cmd
}

// Configure 配置命令
func (c *RootCommand) Configure(logger logger.Logger, redisLocker *locker.RedisLocker) {
	c.Command.PersistentFlags().IntP("instance-id", "i", 0, "实例编号, 防止并发执行")

	// 命令执行之前
	// 注意： c.Command 与 cmd 不是同一个实例，cmd 是实际执行的子命令的实例
	c.Command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		c.curCmd = cmd

		// 设置 traceID 追踪日志
		traceId := uuid.New().String()
		logger := logger.With("traceId", traceId)
		logger.Info("Start Command")

		// 将日志组件注入上下文
		ctx := logger.WithContext(cmd.Context())
		ctx = context.WithValue(ctx, ctxkeys.TraceIDKey{}, traceId)
		ctx = context.WithValue(ctx, ctxkeys.BeginTimeKey{}, time.Now())

		// 为进程加锁逻辑
		id, _ := cmd.Flags().GetInt("instance-id")
		key := fmt.Sprintf("%s:%d", strings.ReplaceAll(cmd.CommandPath(), " ", "-"), id)
		lock, err := redisLocker.Lock(ctx, key, locker.WithTTL(time.Minute), locker.WithWatchDog())
		if err != nil {
			return err
		}
		c.lock = lock
		logger.Info("Locked Command", "key", key)

		cmd.SetContext(ctx)
		return nil
	}

	// 命令执行结束
	c.Command.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		var cost time.Duration
		if val := cmd.Context().Value(ctxkeys.BeginTimeKey{}); val != nil {
			cost = time.Since(val.(time.Time))
		}

		logger2.Ctx(cmd.Context()).Info("End Command",
			"cost", fmt.Sprintf("%.4f", cost.Seconds()))
	}
}

// Cmd 获取当前执行的命令
func (c *RootCommand) Cmd() *cobra.Command {
	if c.curCmd != nil {
		return c.curCmd
	}

	return c.Command
}

// Run 执行命令
func (c *RootCommand) Run() {
	sigChan := make(chan os.Signal, 1)
	// 监听退出信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// 通过上下文控制取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听命令执行完成信号
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer c.unlock()

		err := c.Command.ExecuteContext(ctx)
		if err != nil {
			c.handleError(err)
		}
	}()

	select {
	case <-sigChan:
		logger := logger2.Ctx(c.Cmd().Context())
		logger.Info("Exiting manually...")
		cancel()
		timeout := time.NewTimer(10 * time.Second)
		defer timeout.Stop()
		select {
		case <-done:
		case <-timeout.C:
			logger.Warn("Timeout exit")
		}
	case <-done:
	}
}

// handleError 处理错误
func (c *RootCommand) handleError(err error) {
	ctx := c.Cmd().Context()

	var cost time.Duration
	if val := ctx.Value(ctxkeys.BeginTimeKey{}); val != nil {
		cost = time.Since(val.(time.Time))
	}

	logger := logger2.Ctx(ctx)
	logger.Err(err, "Error Command", "cost", cost)
}

// unlock 解锁
func (c *RootCommand) unlock() {
	if c.lock == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	c.lock.Unlock(ctx)
	logger := logger2.Ctx(c.Cmd().Context())
	logger.Info("Unlocked Command")
}
