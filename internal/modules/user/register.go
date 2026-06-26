package user

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
}

// ApiDeps API 端模块依赖
type ApiDeps struct {
	DB      *gorm.DB
	Binder  *bind.Binder
	RDB     *redis.Client
	AuthCfg Config
}

// Register 注册用户模块路由（管理端）
func Register(rg *gin.RouterGroup, deps Deps) {
	repo := NewRepository(deps.DB)
	pointsLogRepo := NewPointsLogRepo(deps.DB)
	svc := NewService(deps.DB, repo, pointsLogRepo)
	h := newHandler(deps.Binder, svc)

	rg.POST("/user/user/list", h.List)
	rg.POST("/user/user/create", h.Create)
	rg.POST("/user/user/update", h.Update)
	rg.POST("/user/user/delete", h.Delete)
	rg.POST("/user/user/detail", h.Detail)
	rg.POST("/user/points/logs", h.GetPointsLogs)
	rg.POST("/user/points/change", h.ChangePoints)
}

// RegisterApi 注册用户模块 API 端路由。public 注册公开路由，protected 注册需认证的路由并挂载 JWT 中间件。
func RegisterApi(public, protected *gin.RouterGroup, deps ApiDeps) {
	repo := NewRepository(deps.DB)
	authServ := NewAuthService(repo, deps.AuthCfg, deps.RDB)
	authMid := NewAuthMiddleware(deps.AuthCfg, authServ)
	pointsLogRepo := NewPointsLogRepo(deps.DB)
	svc := NewService(deps.DB, repo, pointsLogRepo)
	h := newAuthHandler(deps.Binder, svc, authServ, deps.AuthCfg)

	public.POST("/login", h.login)
	public.POST("/refresh-token", h.refreshToken)

	protected.Use(middleware.JWTGuard(deps.AuthCfg.HmacSecret))
	protected.Use(authMid.LoadUser())

	protected.POST("/logout", h.logout)
	protected.POST("/change-password", h.changePassword)
}
