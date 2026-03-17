package server

import (
	"net/http"
	"zero-backend/internal/middleware"
	"zero-backend/modules/admin/controller"
	adminMiddleware "zero-backend/modules/admin/middleware"

	"github.com/gin-gonic/gin"
)

func NewGin(
	ctrl *controller.Controllers,
	middlewares *middleware.Middlewares,
	adminMiddlewares *adminMiddleware.Middlewares,
) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors.Handle())
	r.Use(middlewares.Before.Handle())

	apiGroup := r.Group("/api")

	// 系统登录
	apiGroup.POST("/login", ctrl.AuthController.Login)
	apiGroup.POST("/refresh-token", ctrl.AuthController.RefreshToken)

	// 注册中间件验证权限
	apiGroup.Use(adminMiddlewares.Auth.JWTAuth())
	apiGroup.Use(adminMiddlewares.Auth.CheckAPIPermission())

	// 鉴权相关接口
	apiGroup.POST("/logout", ctrl.AuthController.Logout)
	apiGroup.POST("/change-password", ctrl.AuthController.ChangePassword)
	apiGroup.POST("/permissions", ctrl.AuthController.GetPermissions)

	// RBAC权限管理
	apiGroup.POST("/rbac/menu/list", ctrl.RbacMenuController.List)
	apiGroup.POST("/rbac/menu/create", ctrl.RbacMenuController.Create)
	apiGroup.POST("/rbac/menu/update", ctrl.RbacMenuController.Update)
	apiGroup.POST("/rbac/menu/delete", ctrl.RbacMenuController.Delete)
	apiGroup.POST("/rbac/menu/sync", ctrl.RbacMenuController.Sync)
	apiGroup.POST("/rbac/menu/api/list", ctrl.RbacMenuController.ApiList)
	apiGroup.POST("/rbac/menu/api/save", ctrl.RbacMenuController.ApiSave)
	apiGroup.POST("/rbac/api/list", ctrl.RbacApiController.List)
	apiGroup.POST("/rbac/api/create", ctrl.RbacApiController.Create)
	apiGroup.POST("/rbac/api/update", ctrl.RbacApiController.Update)
	apiGroup.POST("/rbac/api/delete", ctrl.RbacApiController.Delete)
	apiGroup.POST("/rbac/store/list", ctrl.RbacStoreController.List)
	apiGroup.POST("/rbac/store/create", ctrl.RbacStoreController.Create)
	apiGroup.POST("/rbac/store/update", ctrl.RbacStoreController.Update)
	apiGroup.POST("/rbac/store/delete", ctrl.RbacStoreController.Delete)
	apiGroup.POST("/rbac/store/recycle", ctrl.RbacStoreController.Recycle)
	apiGroup.POST("/rbac/store/restore", ctrl.RbacStoreController.Restore)
	apiGroup.POST("/rbac/role/list", ctrl.RbacRoleController.List)
	apiGroup.POST("/rbac/role/create", ctrl.RbacRoleController.Create)
	apiGroup.POST("/rbac/role/update", ctrl.RbacRoleController.Update)
	apiGroup.POST("/rbac/role/delete", ctrl.RbacRoleController.Delete)
	apiGroup.POST("/rbac/role/set-menus", ctrl.RbacRoleController.SetMenus)
	apiGroup.POST("/rbac/user/list", ctrl.RbacUserController.List)
	apiGroup.POST("/rbac/user/create", ctrl.RbacUserController.Create)
	apiGroup.POST("/rbac/user/update", ctrl.RbacUserController.Update)
	apiGroup.POST("/rbac/user/delete", ctrl.RbacUserController.Delete)
	apiGroup.POST("/rbac/user/set-roles", ctrl.RbacUserController.SetRoles)
	apiGroup.POST("/rbac/user/reset-password", ctrl.RbacUserController.ResetPassword)

	// 设置管理
	apiGroup.POST("/setting/list", ctrl.SettingController.List)
	apiGroup.POST("/setting/create", ctrl.SettingController.Create)
	apiGroup.POST("/setting/update", ctrl.SettingController.Update)
	apiGroup.POST("/setting/delete", ctrl.SettingController.Delete)
	apiGroup.POST("/setting/form-configs", ctrl.SettingController.FormConfigs)
	apiGroup.POST("/setting/qiniu-token", ctrl.SettingController.QiniuToken)
	apiGroup.POST("/setting/default/list", ctrl.SettingDefaultController.List)
	apiGroup.POST("/setting/default/create", ctrl.SettingDefaultController.Create)
	apiGroup.POST("/setting/default/update", ctrl.SettingDefaultController.Update)
	apiGroup.POST("/setting/default/delete", ctrl.SettingDefaultController.Delete)

	// 文章管理
	apiGroup.POST("/article/category/list", ctrl.ArticleCategoryController.List)
	apiGroup.POST("/article/category/create", ctrl.ArticleCategoryController.Create)
	apiGroup.POST("/article/category/update", ctrl.ArticleCategoryController.Update)
	apiGroup.POST("/article/category/delete", ctrl.ArticleCategoryController.Delete)
	apiGroup.POST("/article/article/list", ctrl.ArticleController.List)
	apiGroup.POST("/article/article/create", ctrl.ArticleController.Create)
	apiGroup.POST("/article/article/update", ctrl.ArticleController.Update)
	apiGroup.POST("/article/article/delete", ctrl.ArticleController.Delete)

	// 用户管理
	apiGroup.POST("/user/user/list", ctrl.UserController.List)
	apiGroup.POST("/user/user/create", ctrl.UserController.Create)
	apiGroup.POST("/user/user/update", ctrl.UserController.Update)
	apiGroup.POST("/user/user/delete", ctrl.UserController.Delete)
	apiGroup.POST("/user/user/detail", ctrl.UserController.Detail)
	apiGroup.POST("/user/points/logs", ctrl.UserController.GetPointsLogs)
	apiGroup.POST("/user/points/change", ctrl.UserController.ChangePoints)

	// 文件上传管理
	apiGroup.POST("/upload/group/list", ctrl.UploadGroupController.List)
	apiGroup.POST("/upload/group/create", ctrl.UploadGroupController.Create)
	apiGroup.POST("/upload/group/update", ctrl.UploadGroupController.Update)
	apiGroup.POST("/upload/group/delete", ctrl.UploadGroupController.Delete)
	apiGroup.POST("/upload/file/list", ctrl.UploadFileController.List)
	apiGroup.POST("/upload/file/upload", ctrl.UploadFileController.Upload)
	apiGroup.POST("/upload/file/delete", ctrl.UploadFileController.Delete)

	// 区划管理
	apiGroup.POST("/region/tree", ctrl.RegionController.Tree)

	// 检测服务健康接口
	r.GET("/health", ctrl.HealthController.Health)

	// 设置视图文件目录
	r.LoadHTMLGlob("./views/*.html")

	// 静态资源文件目录
	r.Static("/assets", "./views/assets")
	r.Static("/uploads", "./uploads")

	// 默认图标
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./views/favicon.ico")
	})
	r.GET("/logo.svg", func(c *gin.Context) {
		c.File("./views/logo.svg")
	})

	// 默认首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return r
}
