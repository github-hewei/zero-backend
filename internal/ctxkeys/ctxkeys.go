package ctxkeys

import (
	"context"
	"zero-backend/modules/rbac"
	"zero-backend/modules/user"

	"github.com/241x/zero-web/ctxkeys"
)

// UserID 从上下文中获取用户 ID，兼容 RbacUser 与 User 两种模型。
func UserID(ctx context.Context) uint32 {
	if u, ok := ctxkeys.User(ctx).(*rbac.RbacUser); ok {
		return u.ID
	}
	if u, ok := ctxkeys.User(ctx).(*user.User); ok {
		return u.ID
	}
	return 0
}

// IsSuperUser 判断当前用户是否为超级管理员。
func IsSuperUser(ctx context.Context) bool {
	if u, ok := ctxkeys.User(ctx).(*rbac.RbacUser); ok {
		return u.SU
	}
	return false
}
