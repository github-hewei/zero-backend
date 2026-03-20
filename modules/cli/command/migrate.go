package command

import (
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// MigrateCommand 数据迁移命令
type MigrateCommand struct {
	*cobra.Command
}

// NewMigrateCommand 创建数据迁移命令
func NewMigrateCommand(up *MigrateUpCommand) *MigrateCommand {
	cmd := &MigrateCommand{
		Command: &cobra.Command{
			Use:   "migrate",
			Short: "数据迁移命令",
			Long:  `数据库迁移操作，执行 SQL 脚本进行数据库初始化`,
		},
	}

	cmd.AddCommand(up.Command)
	return cmd
}

// MigrateUpCommand 创建向上迁移命令
type MigrateUpCommand struct {
	*cobra.Command
	db *gorm.DB
}

// NewMigrateUpCommand 创建向上迁移命令
func NewMigrateUpCommand(db *gorm.DB) *MigrateUpCommand {
	cmd := &MigrateUpCommand{
		Command: &cobra.Command{
			Use:   "up",
			Short: "执行数据库迁移",
			Long:  `执行 SQL 脚本进行数据库迁移，从上次的进度继续执行`,
		},
		db: db,
	}

	cmd.Configure()
	return cmd
}

// Configure 配置命令
func (c *MigrateUpCommand) Configure() {
	var filePath string
	c.Flags().StringVarP(&filePath, "file", "f", "", "SQL 文件路径")

	c.Command.RunE = func(cmd *cobra.Command, args []string) error {
		cmd.PrintErrln("测试错误")
		cmd.Println("测试成功")
		// return apperror.NewUserError("测试错误")
		return nil
		// migrateRunner := runner.NewMigrateRunner(logger.Ctx(cmd.Context()), c.db)
		// return migrateRunner.Up(filePath)
	}
}
