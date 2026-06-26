package main

import (
	"zero-backend/cmd/cli/runner"

	"github.com/241x/zero-kit/logger"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func migrateCmd(db *gorm.DB, log logger.Logger) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "数据迁移命令",
		Long:  `执行 SQL 脚本进行数据库迁移，从上次的进度继续执行`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			r := runner.NewMigrateRunner(log, db)
			return r.Up(ctx, filePath)
		},
	}
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "SQL 文件路径 (默认 data/database.sql)")
	return cmd
}
