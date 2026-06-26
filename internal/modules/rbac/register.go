package rbac

import (
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Deps 模块依赖
type Deps struct {
	DB      *gorm.DB
	Binder  *bind.Binder
	Config  Config
	RDB     *redis.Client
	Captcha CaptchaVerifier
}

func buildAll(deps Deps) (*handler, *AuthMiddleware) {
	userRepo := NewRbacUserRepository(deps.DB)
	userRoleRepo := NewRbacUserRoleRepository(deps.DB)
	menuRepo := NewRbacMenuRepository(deps.DB)
	menuApiRepo := NewRbacMenuApiRepository(deps.DB)
	apiRepo := NewRbacApiRepository(deps.DB)
	roleRepo := NewRbacRoleRepository(deps.DB)
	roleMenuRepo := NewRbacRoleMenuRepository(deps.DB)
	storeRepo := NewRbacStoreRepository(deps.DB)

	authServ := NewAuthService(
		userRepo, apiRepo, roleRepo, menuRepo,
		roleMenuRepo, userRoleRepo, menuApiRepo,
		deps.Config, deps.RDB, deps.Captcha,
	)

	authMid := NewAuthMiddleware(deps.Config, authServ)

	menuServ := NewRbacMenuService(menuRepo, menuApiRepo, deps.DB)
	apiServ := NewRbacApiService(apiRepo)
	roleServ := NewRbacRoleService(roleRepo, roleMenuRepo, deps.DB)
	userServ := NewRbacUserService(deps.DB, userRepo, userRoleRepo)
	storeServ := NewRbacStoreService(storeRepo)

	h := newHandler(deps.Binder, authServ, deps.Config, menuServ, apiServ, roleServ, userServ, storeServ, authMid)

	return h, authMid
}

// Register 注册 rbac 模块路由。public 注册公开路由，protected 注册需认证的路由并挂载 JWT 中间件。
func Register(public, protected *gin.RouterGroup, deps Deps) {
	h, authMid := buildAll(deps)

	public.POST("/login", h.login)
	public.POST("/refresh-token", h.refreshToken)

	protected.Use(middleware.JWTGuard(deps.Config.HmacSecret))
	protected.Use(authMid.LoadUser())
	protected.Use(authMid.CheckAPIPermission())

	protected.POST("/logout", h.logout)
	protected.POST("/change-password", h.changePassword)
	protected.POST("/permissions", h.permissions)

	protected.POST("/rbac/menu/list", h.menuList)
	protected.POST("/rbac/menu/create", h.menuCreate)
	protected.POST("/rbac/menu/update", h.menuUpdate)
	protected.POST("/rbac/menu/delete", h.menuDelete)
	protected.POST("/rbac/menu/sync", h.menuSync)
	protected.POST("/rbac/menu/api/list", h.menuApiList)
	protected.POST("/rbac/menu/api/save", h.menuApiSave)

	protected.POST("/rbac/api/list", h.apiList)
	protected.POST("/rbac/api/create", h.apiCreate)
	protected.POST("/rbac/api/update", h.apiUpdate)
	protected.POST("/rbac/api/delete", h.apiDelete)

	protected.POST("/rbac/store/list", h.storeList)
	protected.POST("/rbac/store/create", h.storeCreate)
	protected.POST("/rbac/store/update", h.storeUpdate)
	protected.POST("/rbac/store/delete", h.storeDelete)
	protected.POST("/rbac/store/recycle", h.storeRecycle)
	protected.POST("/rbac/store/restore", h.storeRestore)

	protected.POST("/rbac/role/list", h.roleList)
	protected.POST("/rbac/role/create", h.roleCreate)
	protected.POST("/rbac/role/update", h.roleUpdate)
	protected.POST("/rbac/role/delete", h.roleDelete)
	protected.POST("/rbac/role/set-menus", h.roleSetMenus)

	protected.POST("/rbac/user/list", h.userList)
	protected.POST("/rbac/user/create", h.userCreate)
	protected.POST("/rbac/user/update", h.userUpdate)
	protected.POST("/rbac/user/delete", h.userDelete)
	protected.POST("/rbac/user/set-roles", h.userSetRoles)
	protected.POST("/rbac/user/reset-password", h.userResetPassword)
}
