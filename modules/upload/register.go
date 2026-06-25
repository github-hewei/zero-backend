package upload

import (
	"github.com/241x/zero-kit/bind"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Deps 模块依赖
type Deps struct {
	DB       *gorm.DB
	Binder   *bind.Binder
	Settings SettingProvider
}

func buildHandler(deps Deps) *Handler {
	groupRepo := NewGroupRepository(deps.DB)
	groupSvc := NewGroupService(groupRepo)
	fileRepo := NewFileRepository(deps.DB)
	fileSvc := NewFileService(fileRepo, deps.Settings)
	return newHandler(deps.Binder, groupSvc, fileSvc)
}

// RegisterAdmin 注册上传模块路由（管理端）
func RegisterAdmin(rg *gin.RouterGroup, deps Deps) {
	h := buildHandler(deps)
	rg.POST("/upload/group/list", h.ListGroup)
	rg.POST("/upload/group/create", h.CreateGroup)
	rg.POST("/upload/group/update", h.UpdateGroup)
	rg.POST("/upload/group/delete", h.DeleteGroup)
	rg.POST("/upload/file/list", h.ListFile)
	rg.POST("/upload/file/upload", h.UploadFileAdmin)
	rg.POST("/upload/file/delete", h.DeleteFile)
}

// RegisterApi 注册上传模块路由（API端）
func RegisterApi(rg *gin.RouterGroup, deps Deps) {
	h := buildHandler(deps)
	rg.POST("/upload/file/list", h.ListFile)
	rg.POST("/upload/file/upload", h.UploadFile)
	rg.POST("/upload/file/delete", h.DeleteFile)
}
