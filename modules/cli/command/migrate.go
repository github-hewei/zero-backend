package command

import (
	"zero-backend/internal/logger"
	"zero-backend/modules/cli/runner"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// MigrateCommand 数据迁移命令
type MigrateCommand struct {
	*cobra.Command
	db *gorm.DB
}

// NewMigrateCommand 创建数据迁移命令
func NewMigrateCommand(db *gorm.DB) *MigrateCommand {
	cmd := &MigrateCommand{
		Command: &cobra.Command{
			Use:   "migrate",
			Short: "数据迁移命令",
			Long:  `执行 SQL 脚本进行数据库迁移，从上次的进度继续执行`,
		},
		db: db,
	}

	cmd.Configure()
	return cmd
}

// Configure 配置命令
func (c *MigrateCommand) Configure() {
	var filePath string
	c.Flags().StringVarP(&filePath, "file", "f", "", "SQL 文件路径 (默认 data/database.sql)")

	c.Command.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		migrateRunner := runner.NewMigrateRunner(logger.Ctx(ctx), c.db)
		return migrateRunner.Up(ctx, filePath)
	}
}
