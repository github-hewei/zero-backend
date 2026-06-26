package api

import (
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
	"gorm.io/gorm"
)

func NewGin(
	log logger.Logger,
	db *gorm.DB,
	binder *bind.Binder,
	rdb *redis.Client,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(app.LoadApiCorsConfig()))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")

	authMid, h := user.RegisterApi(apiGroup, user.ApiDeps{
		DB:      db,
		Binder:  binder,
		AuthCfg: user.LoadConfig(),
		RDB:     rdb,
	})

	apiGroup.Use(middleware.JWTGuard(user.LoadConfig().HmacSecret))
	apiGroup.Use(authMid.LoadUser())

	user.RegisterApiProtected(apiGroup, h)

	settingSvc := app.NewSettingService(db)

	upload.RegisterApi(apiGroup, upload.Deps{DB: db, Binder: binder, Settings: settingSvc})
	setting.RegisterApi(apiGroup, setting.Deps{DB: db, Binder: binder})
	region.Register(apiGroup, region.Deps{DB: db, Binder: binder})

	return r
}
