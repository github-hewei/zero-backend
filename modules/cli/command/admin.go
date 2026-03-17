package command

import (
	"zero-backend/modules/cli"

	"github.com/spf13/cobra"
)

// NewAdminCommand 创建管理员命令
func NewAdminCommand(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admin",
		Short: "管理员命令",
		Long:  `管理员相关操作`,
	}

	cmd.AddCommand(NewAdminListCommand(ctx))

	return cmd
}

// NewAdminListCommand 创建管理员列表命令
func NewAdminListCommand(ctx *cli.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有管理员",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Logger.Info("列出所有管理员")
			return nil
		},
	}
}
