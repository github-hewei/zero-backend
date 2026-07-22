package user

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RegisterAdmin 注册用户模块路由（管理端）
func RegisterAdmin(rg *gin.RouterGroup, db *gorm.DB, binder *bind.Binder) {
	repo := NewRepository(db)
	pointsLogRepo := NewPointsLogRepository(db)
	svc := NewService(db, repo, pointsLogRepo)
	h := newHandler(binder, svc)

	rg.POST("/user/user/list", h.List)
	rg.POST("/user/user/create", h.Create)
	rg.POST("/user/user/update", h.Update)
	rg.POST("/user/user/delete", h.Delete)
	rg.POST("/user/user/detail", h.Detail)
	rg.POST("/user/points/logs", h.GetPointsLogs)
	rg.POST("/user/points/change", h.ChangePoints)
}

// RegisterApi 注册用户模块 API 端路由。public 注册公开路由，protected 注册需认证的路由并挂载 JWT 中间件。
func RegisterApi(public, protected *gin.RouterGroup, db *gorm.DB, binder *bind.Binder, rdb *redis.Client, cfg Config) {
	repo := NewRepository(db)
	authServ := NewAuthService(repo, cfg, rdb)
	authMid := NewAuthMiddleware(cfg, authServ)
	pointsLogRepo := NewPointsLogRepository(db)
	svc := NewService(db, repo, pointsLogRepo)
	h := newAuthHandler(binder, svc, authServ, cfg)

	public.POST("/login", h.login)
	public.POST("/refresh-token", h.refreshToken)

	protected.Use(middleware.JWTGuard(cfg.HmacSecret))
	protected.Use(authMid.LoadUser())

	protected.POST("/logout", h.logout)
	protected.POST("/change-password", h.changePassword)
}
