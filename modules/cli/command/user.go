package command

import (
	"zero-backend/modules/cli"

	"github.com/spf13/cobra"
)

// NewUserCommand 创建用户命令
func NewUserCommand(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "用户管理命令",
		Long:  `用户相关操作，包括创建、查询、修改等`,
	}

	cmd.AddCommand(NewUserListCommand(ctx))
	cmd.AddCommand(NewUserCreateCommand(ctx))

	return cmd
}

// NewUserListCommand 创建用户列表命令
func NewUserListCommand(ctx *cli.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "列出所有用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Logger.Info("列出所有用户")
			return nil
		},
	}
}

// NewUserCreateCommand 创建用户命令
func NewUserCreateCommand(ctx *cli.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "创建新用户",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Logger.Info("创建新用户")
			return nil
		},
	}
}
