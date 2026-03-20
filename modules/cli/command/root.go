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
	logger logger.Logger
}

// NewRootCommand 创建根命令
func NewRootCommand(
	l logger.Logger,
	user *UserCommand,
	migrate *MigrateCommand,
) *RootCommand {
	cmd := &RootCommand{
		Command: &cobra.Command{
			Use:   "cli",
			Short: "CLI 工具",
			Long:  `命令行管理工具，用于执行数据迁移、用户管理等操作`,
		},
		logger: l,
	}

	cmd.Configure()
	cmd.AddCommand(user.Command)
	cmd.AddCommand(migrate.Command)
	return cmd
}

// Configure 配置命令
func (c *RootCommand) Configure() {
	// 命令执行之前
	c.Command.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		traceId := uuid.New().String()
		logger := c.logger.With("traceId", traceId)
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
