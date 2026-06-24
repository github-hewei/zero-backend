package router

import (
	"zero-backend/internal/config"
	"zero-backend/modules/api/controller"
	apiMiddleware "zero-backend/modules/api/middleware"

	"github.com/241x/zero-kit/logger"
	basecfg "github.com/241x/zero-web/config"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
)

func NewGin(
	log logger.Logger,
	ctrl *controller.Controllers,
	apiMiddlewares *apiMiddleware.Middlewares,
	corsConfig basecfg.CorsConfig,
	authConfig config.ApiAuthConfig,
) *gin.Engine {
	r := gin.Default()
	// CORS 跨域配置
	r.Use(middleware.CORS(corsConfig))
	// Trace 跟踪配置
	r.Use(middleware.Trace(log))
	// Request 日志配置
	r.Use(middleware.RequestLog())

	apiGroup := r.Group("/api")

	// 系统登录
	apiGroup.POST("/login", ctrl.AuthController.Login)
	apiGroup.POST("/refresh-token", ctrl.AuthController.RefreshToken)

	// 注册中间件验证权限
	apiGroup.Use(middleware.JWTGuard(authConfig.HmacSecret))
	apiGroup.Use(apiMiddlewares.Auth.LoadUser())

	// 鉴权相关接口
	apiGroup.POST("/logout", ctrl.AuthController.Logout)
	apiGroup.POST("/change-password", ctrl.AuthController.ChangePassword)

	// 文件上传管理
	apiGroup.POST("/upload/file/list", ctrl.UploadFileController.List)
	apiGroup.POST("/upload/file/upload", ctrl.UploadFileController.Upload)
	apiGroup.POST("/upload/file/delete", ctrl.UploadFileController.Delete)

	// 设置相关接口
	apiGroup.POST("/setting/qiniu-token", ctrl.SettingController.QiniuToken)

	// 区域相关接口
	apiGroup.POST("/region/tree", ctrl.RegionController.Regions)

	return r
}
