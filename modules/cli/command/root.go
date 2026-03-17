package command

import (
	"context"
	"fmt"
	"zero-backend/modules/cli"

	"github.com/spf13/cobra"
)

// NewRootCommand 创建根命令
func NewRootCommand(ctx *cli.Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "zero-backend",
		Short: "zero-backend CLI 工具",
		Long:  `zero-backend 命令行管理工具，用于执行数据迁移、用户管理等操作`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			c, ok := cli.FromContext(cmd.Context())
			if !ok || c == nil {
				return fmt.Errorf("CLI 上下文未初始化")
			}
			return nil
		},
	}

	// 添加子命令
	rootCmd.AddCommand(NewUserCommand(ctx))
	rootCmd.AddCommand(NewAdminCommand(ctx))
	rootCmd.AddCommand(NewMigrateCommand(ctx))

	// 全局 flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出")
	rootCmd.PersistentFlags().StringP("config", "c", "config.yaml", "配置文件路径")

	return rootCmd
}

// Execute 执行根命令
func Execute(ctx *cli.Context) error {
	rootCmd := NewRootCommand(ctx)
	// 使用 context.Background() 作为基础，避免 nil context
	rootCmd.SetContext(cli.WithContext(context.Background(), ctx))
	return rootCmd.Execute()
}
