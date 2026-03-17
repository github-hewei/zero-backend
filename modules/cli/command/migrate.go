package command

import (
	"zero-backend/modules/cli"

	"github.com/spf13/cobra"
)

// NewMigrateCommand 创建数据迁移命令
func NewMigrateCommand(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "数据迁移命令",
		Long:  `数据库迁移操作，包括初始化、迁移、回滚等`,
	}

	cmd.AddCommand(NewMigrateUpCommand(ctx))
	cmd.AddCommand(NewMigrateDownCommand(ctx))
	cmd.AddCommand(NewMigrateInitCommand(ctx))

	return cmd
}

// NewMigrateUpCommand 创建向上迁移命令
func NewMigrateUpCommand(ctx *cli.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "执行数据库迁移",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Logger.Info("执行数据库迁移")
			return nil
		},
	}
}

// NewMigrateDownCommand 创建向下迁移命令
func NewMigrateDownCommand(ctx *cli.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "回滚数据库迁移",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Logger.Info("回滚数据库迁移")
			return nil
		},
	}
}

// NewMigrateInitCommand 创建初始化迁移命令
func NewMigrateInitCommand(ctx *cli.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "初始化数据库",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Logger.Info("初始化数据库")
			return nil
		},
	}
}
