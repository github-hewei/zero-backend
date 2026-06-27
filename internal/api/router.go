package api

import (
	_ "zero-backend/docs/api"
	"zero-backend/internal/app"
	"zero-backend/internal/modules/region"
	"zero-backend/internal/modules/setting"
	"zero-backend/internal/modules/upload"
	"zero-backend/internal/modules/user"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/logger"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// NewGin 创建一个 gin.Engine 实例
func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	rdb *redis.Client,
	authCfg user.Config,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(app.LoadApiCorsConfig()))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")
	protected := apiGroup.Group("")

	user.RegisterApi(apiGroup, protected, user.ApiDeps{
		DB:      db,
		Binder:  binder,
		AuthCfg: authCfg,
		RDB:     rdb,
	})

	settingSvc := app.NewSettingService(db)

	upload.RegisterApi(protected, upload.Deps{DB: db, Binder: binder, Settings: settingSvc})
	setting.RegisterApi(protected, setting.Deps{DB: db, Binder: binder})
	region.Register(protected, region.Deps{DB: db, Binder: binder})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
