package command

import (
	"context"
	"fmt"
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

// Run 执行命令
func (c *RootCommand) Run() {
	err := c.Command.Execute()
	if err != nil {
		c.handleError(err)
	}
}

// handleError 处理错误
func (c *RootCommand) handleError(err error) {
	cmd := c.Command
	if c.curCmd != nil {
		cmd = c.curCmd
	}

	var cost time.Duration
	if val := cmd.Context().Value(ctxkeys.BeginTimeKey{}); val != nil {
		cost = time.Since(val.(time.Time))
	}

	logger := logger2.Ctx(cmd.Context())
	logger.Err(err, "Error Command", "cost", cost)
}
