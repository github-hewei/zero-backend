package setting

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

// buildHandler 创建设置模块处理器
func buildHandler(deps Deps) *Handler {
	repo := NewRepository(deps.DB)
	defaultRepo := NewDefaultRepository(deps.DB)
	svc := NewService(repo, defaultRepo)
	defaultSvc := NewDefaultService(defaultRepo)
	return newHandler(deps.Binder, svc, defaultSvc)
}

// RegisterAdmin 注册设置模块路由（管理端）
func RegisterAdmin(rg *gin.RouterGroup, deps Deps) {
	h := buildHandler(deps)
	rg.POST("/setting/list", h.List)
	rg.POST("/setting/create", h.Create)
	rg.POST("/setting/update", h.Update)
	rg.POST("/setting/delete", h.Delete)
	rg.POST("/setting/form-configs", h.FormConfigs)
	rg.POST("/setting/qiniu-token", h.QiniuToken)
	rg.POST("/setting/default/list", h.DefaultList)
	rg.POST("/setting/default/create", h.DefaultCreate)
	rg.POST("/setting/default/update", h.DefaultUpdate)
	rg.POST("/setting/default/delete", h.DefaultDelete)
}

// RegisterApi 注册设置模块路由（API端）
func RegisterApi(rg *gin.RouterGroup, deps Deps) {
	h := buildHandler(deps)
	rg.POST("/setting/qiniu-token", h.QiniuToken)
}
