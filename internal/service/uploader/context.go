package uploader

import (
	"context"

	"zero-backend/internal/dto"
)

// qiniuConfigKey 上下文传递七牛云配置
type qiniuConfigKey struct{}

// WithQiniuConfig 注入七牛云配置
func WithQiniuConfig(ctx context.Context, cfg *dto.QiniuConfig) context.Context {
	return context.WithValue(ctx, qiniuConfigKey{}, cfg)
}

// QiniuConfig 读取七牛云配置
func QiniuConfig(ctx context.Context) *dto.QiniuConfig {
	v, ok := ctx.Value(qiniuConfigKey{}).(*dto.QiniuConfig)
	if !ok {
		return nil
	}
	return v
}
