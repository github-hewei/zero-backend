package command

import (
	"fmt"
	"zero-backend/internal/dto"
	"zero-backend/internal/model"
	"zero-backend/internal/service"
	"zero-backend/pkg/logger"

	"github.com/spf13/cobra"
)

// UserCommand 用户命令
type UserCommand struct {
	*cobra.Command
	logger logger.Logger
}

// NewUserCommand 创建用户命令
func NewUserCommand(l logger.Logger, list *UserListCommand) *UserCommand {
	cmd := &UserCommand{
		Command: &cobra.Command{
			Use:   "user",
			Short: "用户管理",
			Long:  `用户管理命令行工具，用于执行用户管理操作`,
		},
		logger: l,
	}

	cmd.AddCommand(list.Command)
	return cmd
}

// UserListCommand 获取用户列表命令
type UserListCommand struct {
	*cobra.Command
	logger logger.Logger
	svc    *service.UserService
}

// NewUserListCommand 创建获取用户列表命令
func NewUserListCommand(l logger.Logger, svc *service.UserService) *UserListCommand {
	cmd := &UserListCommand{
		Command: &cobra.Command{
			Use:   "list",
			Short: "列出所有用户",
			Long:  `列出所有用户`,
		},
		logger: l,
		svc:    svc,
	}

	cmd.Command.RunE = cmd.RunE
	return cmd
}

// Run 运行命令
func (c *UserListCommand) RunE(cmd *cobra.Command, args []string) error {
	c.logger.Info("列出所有用户")
	result, err := c.svc.List(cmd.Context(), &dto.UserListRequest{Page: 1, Limit: 10})
	if err != nil {
		return err
	}

	users, _ := result.List.([]*model.User)
	for _, item := range users {
		fmt.Printf("%d: %s\n", item.ID, item.NickName)
	}
	return nil
}
