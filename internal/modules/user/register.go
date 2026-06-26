package user

import (
	"zero-backend/internal/config"

	"github.com/241x/zero-kit/bind"
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
	AuthCfg config.ApiAuthConfig
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

// RegisterApi 注册用户模块路由（API端），返回 AuthMiddleware 和 Handler
func RegisterApi(rg *gin.RouterGroup, deps ApiDeps) (*AuthMiddleware, *Handler) {
	repo := NewRepository(deps.DB)
	authServ := NewAuthService(repo, deps.AuthCfg, deps.RDB)
	authMid := NewAuthMiddleware(deps.AuthCfg, authServ)
	pointsLogRepo := NewPointsLogRepo(deps.DB)
	svc := NewService(deps.DB, repo, pointsLogRepo)
	h := newAuthHandler(deps.Binder, svc, authServ, deps.AuthCfg)

	rg.POST("/login", h.login)
	rg.POST("/refresh-token", h.refreshToken)

	return authMid, h
}

// RegisterApiProtected 注册需要 JWT 保护的路由
func RegisterApiProtected(rg *gin.RouterGroup, h *Handler) {
	rg.POST("/logout", h.logout)
	rg.POST("/change-password", h.changePassword)
}
