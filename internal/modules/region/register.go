package region

import (
	"github.com/241x/zero-kit/bind"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register 注册区划模块路由
func Register(rg *gin.RouterGroup, db *gorm.DB, binder *bind.Binder) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := newHandler(binder, svc)

	rg.POST("/region/tree", h.Tree)
}
