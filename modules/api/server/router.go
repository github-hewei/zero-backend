package server

import (
	"zero-backend/internal/middleware"
	"zero-backend/modules/api/controller"
	apiMiddleware "zero-backend/modules/api/middleware"

	"github.com/gin-gonic/gin"
)

func NewGin(
	ctrl *controller.Controllers,
	middlewares *middleware.Middlewares,
	apiMiddleware *apiMiddleware.Middlewares,
) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Before.Handle())

	apiGroup := r.Group("/api")

	// 系统登录
	apiGroup.POST("/login", ctrl.AuthController.Login)
	apiGroup.POST("/refresh-token", ctrl.AuthController.RefreshToken)

	// 注册中间件验证权限
	apiGroup.Use(apiMiddleware.Auth.JWTAuth())

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
