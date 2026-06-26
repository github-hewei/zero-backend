package article

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

// Register 注册文章模块路由
func Register(rg *gin.RouterGroup, deps Deps) {
	categoryRepo := NewCategoryRepository(deps.DB)
	categorySvc := NewCategoryService(categoryRepo)

	articleRepo := NewRepository(deps.DB)
	articleSvc := NewService(articleRepo)

	h := NewHandler(deps.Binder, categorySvc, articleSvc)

	rg.POST("/article/category/list", h.ListCategory)
	rg.POST("/article/category/create", h.CreateCategory)
	rg.POST("/article/category/update", h.UpdateCategory)
	rg.POST("/article/category/delete", h.DeleteCategory)
	rg.POST("/article/article/list", h.List)
	rg.POST("/article/article/create", h.Create)
	rg.POST("/article/article/update", h.Update)
	rg.POST("/article/article/delete", h.Delete)
}
