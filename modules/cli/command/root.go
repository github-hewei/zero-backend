package command

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zero-backend/internal/ctxkeys"
	logger2 "zero-backend/internal/logger"
	"zero-backend/pkg/logger"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// RootCommand 根命令
type RootCommand struct {
	*cobra.Command
	curCmd *cobra.Command
}

// NewRootCommand 创建根命令
func NewRootCommand(
	logger logger.Logger,
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

	cmd.Configure(logger)
	cmd.AddCommand(user.Command)
	cmd.AddCommand(migrate.Command)
	return cmd
}

// Configure 配置命令
func (c *RootCommand) Configure(logger logger.Logger) {
	// 命令执行之前
	// 注意： c.Command 与 cmd 不是同一个实例，cmd 是实际执行的子命令的实例
	c.Command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		c.curCmd = cmd
		traceId := uuid.New().String()
		logger := logger.With("traceId", traceId)
		logger.Info("Start Command")

		ctx := logger.WithContext(cmd.Context())
		ctx = context.WithValue(ctx, ctxkeys.TraceIDKey{}, traceId)
		ctx = context.WithValue(ctx, ctxkeys.BeginTimeKey{}, time.Now())
		cmd.SetContext(ctx)
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
