package runner

import (
	"zero-backend/modules/cli"
)

// UserRunner 用户命令执行器
type UserRunner struct {
	ctx *cli.Context
}

// NewUserRunner 创建用户命令执行器
func NewUserRunner(ctx *cli.Context) *UserRunner {
	return &UserRunner{ctx: ctx}
}

// List 执行用户列表
func (r *UserRunner) List() error {
	r.ctx.Logger.Info("执行用户列表命令")
	return nil
}

// Create 执行创建用户
func (r *UserRunner) Create(username, email string) error {
	r.ctx.Logger.Info("创建用户", "username", username, "email", email)
	return nil
}
