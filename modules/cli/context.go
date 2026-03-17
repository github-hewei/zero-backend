package cli

import (
	"context"

	"zero-backend/internal/config"
	"zero-backend/internal/ctxkeys"
	"zero-backend/pkg/logger"

	"gorm.io/gorm"
)

// Context CLI 上下文，封装依赖供命令使用
type Context struct {
	Config *config.Config
	Logger logger.Logger
	DB     *gorm.DB
}

// NewContext 创建 CLI 上下文
func NewContext(cfg *config.Config, logger logger.Logger, db *gorm.DB) *Context {
	return &Context{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}
}

// WithContext 将 CLI 上下文添加到 context.Context
func WithContext(ctx context.Context, c *Context) context.Context {
	return context.WithValue(ctx, ctxkeys.CLIContextKey{}, c)
}

// FromContext 从 context.Context 获取 CLI 上下文
func FromContext(ctx context.Context) (*Context, bool) {
	c, ok := ctx.Value(ctxkeys.CLIContextKey{}).(*Context)
	return c, ok
}
