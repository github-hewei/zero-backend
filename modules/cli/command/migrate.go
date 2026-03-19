package command

import (
	"zero-backend/modules/cli/runner"
	"zero-backend/pkg/logger"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// MigrateCommand 数据迁移命令
type MigrateCommand struct {
	*cobra.Command
	logger logger.Logger
}

// NewMigrateCommand 创建数据迁移命令
func NewMigrateCommand(l logger.Logger, up *MigrateUpCommand) *MigrateCommand {
	cmd := &MigrateCommand{
		Command: &cobra.Command{
			Use:   "migrate",
			Short: "数据迁移命令",
			Long:  `数据库迁移操作，执行 SQL 脚本进行数据库初始化`,
		},
		logger: l,
	}

	cmd.AddCommand(up.Command)
	return cmd
}

// MigrateUpCommand 创建向上迁移命令
type MigrateUpCommand struct {
	*cobra.Command
	logger   logger.Logger
	db       *gorm.DB
	filePath string
}

// NewMigrateUpCommand 创建向上迁移命令
func NewMigrateUpCommand(l logger.Logger, db *gorm.DB) *MigrateUpCommand {
	cmd := &MigrateUpCommand{
		Command: &cobra.Command{
			Use:   "up",
			Short: "执行数据库迁移",
			Long:  `执行 SQL 脚本进行数据库迁移，从上次的进度继续执行`,
		},
		logger: l,
		db:     db,
	}

	cmd.Command.RunE = cmd.RunE
	cmd.Flags().StringVarP(&cmd.filePath, "file", "f", "", "SQL 文件路径")
	return cmd
}

// RunE 运行命令
func (c *MigrateUpCommand) RunE(cmd *cobra.Command, args []string) error {
	migrateRunner := runner.NewMigrateRunner(c.logger, c.db)
	return migrateRunner.Up(c.filePath)
}
