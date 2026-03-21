package command

import (
	"time"
	"zero-backend/internal/dto"
	"zero-backend/internal/logger"
	"zero-backend/internal/model"
	"zero-backend/internal/service"

	"github.com/spf13/cobra"
)

// UserCommand 用户命令
type UserCommand struct {
	*cobra.Command
}

// NewUserCommand 创建用户命令
func NewUserCommand(list *UserListCommand) *UserCommand {
	cmd := &UserCommand{
		Command: &cobra.Command{
			Use:   "user",
			Short: "用户管理",
			Long:  `用户管理命令行工具，用于执行用户管理操作`,
		},
	}

	cmd.AddCommand(list.Command)
	return cmd
}

// UserListCommand 获取用户列表命令
type UserListCommand struct {
	*cobra.Command
	svc *service.UserService
}

// NewUserListCommand 创建获取用户列表命令
func NewUserListCommand(svc *service.UserService) *UserListCommand {
	cmd := &UserListCommand{
		Command: &cobra.Command{
			Use:   "list",
			Short: "列出所有用户",
			Long:  `列出所有用户`,
		},
		svc: svc,
	}

	cmd.Configure()
	return cmd
}

// Configure 配置命令
func (c *UserListCommand) Configure() {
	c.Command.RunE = func(cmd *cobra.Command, args []string) error {
		logger := logger.Ctx(cmd.Context())
		logger.Info("列出所有用户")

		time.Sleep(time.Second * 10)
		result, err := c.svc.List(cmd.Context(), &dto.UserListRequest{Page: 1, Limit: 10})
		if err != nil {
			return err
		}

		users, _ := result.List.([]*model.User)
		for _, item := range users {
			cmd.Printf("%d: %s\n", item.ID, item.NickName)
		}
		return nil
	}
}
