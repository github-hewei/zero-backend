package article

import (
	"github.com/241x/zero-kit/bind"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAdmin 注册文章模块路由
func RegisterAdmin(rg *gin.RouterGroup, db *gorm.DB, binder *bind.Binder) {
	categoryRepo := NewCategoryRepository(db)
	categorySvc := NewCategoryService(categoryRepo)

	articleRepo := NewRepository(db)
	articleSvc := NewService(articleRepo)

	h := NewHandler(binder, categorySvc, articleSvc)

	rg.POST("/article/category/list", h.ListCategory)
	rg.POST("/article/category/create", h.CreateCategory)
	rg.POST("/article/category/update", h.UpdateCategory)
	rg.POST("/article/category/delete", h.DeleteCategory)
	rg.POST("/article/article/list", h.List)
	rg.POST("/article/article/create", h.Create)
	rg.POST("/article/article/update", h.Update)
	rg.POST("/article/article/delete", h.Delete)
}
