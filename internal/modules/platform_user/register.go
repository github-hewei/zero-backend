package platform_user

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Deps 模块依赖
type Deps struct {
	DB     *gorm.DB
	Binder *bind.Binder
	Config Config
	RDB    *redis.Client
}

func buildAll(deps Deps) (*handler, *AuthMiddleware) {
	repo := NewPlatformUserRepository(deps.DB)
	authServ := NewAuthService(repo, deps.Config, deps.RDB)
	authMid := NewAuthMiddleware(deps.Config, authServ)
	userServ := NewPlatformUserService(repo)

	h := newHandler(deps.Binder, authServ, deps.Config, userServ)
	return h, authMid
}

// Register 注册平台模块路由。public 注册公开路由，protected 注册需认证的路由并挂载 JWT + 角色中间件。
func Register(public, protected *gin.RouterGroup, deps Deps) *AuthMiddleware {
	h, authMid := buildAll(deps)

	public.POST("/login", h.login)
	public.POST("/refresh-token", h.refreshToken)

	protected.Use(middleware.JWTGuard(deps.Config.HmacSecret))
	protected.Use(authMid.LoadUser())

	protected.POST("/logout", h.logout)
	protected.POST("/change-password", h.changePassword)

	protected.POST("/platform/user/list", h.userList)
	protected.POST("/platform/user/create", h.userCreate)
	protected.POST("/platform/user/update", h.userUpdate)
	protected.POST("/platform/user/delete", h.userDelete)
	protected.POST("/platform/user/reset-password", h.userResetPassword)

	return authMid
}
