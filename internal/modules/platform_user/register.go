package platform_user

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// buildAll 构建所有依赖
func buildAll(db *gorm.DB, binder *bind.Binder, config Config, rdb *redis.Client, captcha CaptchaVerifier) (*handler, *AuthMiddleware) {
	repo := NewPlatformUserRepository(db)
	authServ := NewAuthService(repo, config, rdb, captcha)
	authMid := NewAuthMiddleware(config, authServ)
	userServ := NewPlatformUserService(repo)

	h := newHandler(binder, authServ, config, userServ)
	return h, authMid
}

// Register 注册平台模块路由。public 注册公开路由，protected 注册需认证的路由并挂载 JWT + 角色中间件。
func Register(public, protected *gin.RouterGroup, db *gorm.DB, binder *bind.Binder, config Config, rdb *redis.Client, captcha CaptchaVerifier) *AuthMiddleware {
	h, authMid := buildAll(db, binder, config, rdb, captcha)

	public.POST("/login", h.login)
	public.POST("/refresh-token", h.refreshToken)

	protected.Use(middleware.JWTGuard(config.HmacSecret))
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
