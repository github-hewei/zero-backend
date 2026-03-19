package command

import (
	"zero-backend/pkg/logger"

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

	cmd.AddCommand(user.Command)
	cmd.AddCommand(migrate.Command)
	return cmd
}
