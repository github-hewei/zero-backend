package ctxkeys

import (
	"context"

	"zero-backend/internal/model"

	"github.com/241x/zero-web/ctxkeys"
)

// UserID 从上下文中获取用户 ID，兼容 RbacUser 与 User 两种模型。
func UserID(ctx context.Context) uint32 {
	if user, ok := ctxkeys.User(ctx).(*model.RbacUser); ok {
		return user.ID
	}
	if user, ok := ctxkeys.User(ctx).(*model.User); ok {
		return user.ID
	}
	return 0
}

// IsSuperUser 判断当前用户是否为超级管理员。
func IsSuperUser(ctx context.Context) bool {
	if user, ok := ctxkeys.User(ctx).(*model.RbacUser); ok {
		return user.SU
	}
	return false
}

// storeIdKey 上下文传递企业ID
type storeIdKey struct{}

// WithStoreID 注入企业ID
func WithStoreID(ctx context.Context, id uint32) context.Context {
	return context.WithValue(ctx, storeIdKey{}, id)
}

// StoreID 读取企业ID
func StoreID(ctx context.Context) uint32 {
	v, ok := ctx.Value(storeIdKey{}).(uint32)
	if !ok {
		return 0
	}
	return v
}
