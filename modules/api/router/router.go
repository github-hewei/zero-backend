package router

import (
	"zero-backend/internal/config"
	"zero-backend/internal/service"
	"zero-backend/modules/api/controller"
	apiMiddleware "zero-backend/modules/api/middleware"
	"zero-backend/modules/upload"

	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-kit/logger"
	basecfg "github.com/241x/zero-web/config"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewGin(
	log logger.Logger,
	ctrl *controller.Controllers,
	apiMiddlewares *apiMiddleware.Middlewares,
	corsConfig basecfg.CorsConfig,
	authConfig config.ApiAuthConfig,
	db *gorm.DB,
	binder *bind.Binder,
	settingSvc *service.SettingService,
) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS(corsConfig))
	r.Use(middleware.Trace(log))
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")

	apiGroup.POST("/login", ctrl.AuthController.Login)
	apiGroup.POST("/refresh-token", ctrl.AuthController.RefreshToken)

	apiGroup.Use(middleware.JWTGuard(authConfig.HmacSecret))
	apiGroup.Use(apiMiddlewares.Auth.LoadUser())

	apiGroup.POST("/logout", ctrl.AuthController.Logout)
	apiGroup.POST("/change-password", ctrl.AuthController.ChangePassword)

	upload.RegisterApi(apiGroup, upload.Deps{DB: db, Binder: binder, Settings: settingSvc})

	apiGroup.POST("/setting/qiniu-token", ctrl.SettingController.QiniuToken)
	apiGroup.POST("/region/tree", ctrl.RegionController.Regions)

	return r
}
