package user

import (
	"github.com/241x/zero-kit/bind"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Deps 模块依赖
type Deps struct {
	DB     *gorm.DB
	Binder *bind.Binder
}

// Register 注册用户模块路由
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
