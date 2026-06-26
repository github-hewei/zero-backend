package region

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

// Register 注册区划模块路由
func Register(rg *gin.RouterGroup, deps Deps) {
	repo := NewRepository(deps.DB)
	svc := NewService(repo)
	h := newHandler(deps.Binder, svc)

	rg.POST("/region/tree", h.Tree)
}
