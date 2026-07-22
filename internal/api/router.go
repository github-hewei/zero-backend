package api

import (
	"zero-backend/internal/modules/region"
	"zero-backend/internal/modules/setting"
	"zero-backend/internal/modules/upload"
	"zero-backend/internal/modules/user"
	"zero-backend/internal/provider"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewGin 创建一个 gin.Engine 实例
func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	rdb *redis.Client,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(provider.LoadApiCorsConfig()))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")
	protected := apiGroup.Group("")

	cfg := user.MustLoadConfig()
	user.RegisterApi(apiGroup, protected, db, binder, rdb, cfg)
	settingSvc := provider.NewSettingService(db)
	upload.RegisterApi(protected, db, binder, settingSvc)
	setting.RegisterApi(protected, db, binder)
	region.Register(protected, db, binder)

	return r
}
