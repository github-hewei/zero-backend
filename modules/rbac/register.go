package rbac

import (
	"zero-backend/config"

	"github.com/241x/zero-kit/bind"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Deps 模块依赖
type Deps struct {
	DB      *gorm.DB
	Binder  *bind.Binder
	AuthCfg config.AdminAuthConfig
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
		deps.AuthCfg, deps.RDB, deps.Captcha,
	)

	authMid := NewAuthMiddleware(deps.AuthCfg, authServ)

	menuServ := NewRbacMenuService(menuRepo, menuApiRepo, deps.DB)
	apiServ := NewRbacApiService(apiRepo)
	roleServ := NewRbacRoleService(roleRepo, roleMenuRepo, deps.DB)
	userServ := NewRbacUserService(deps.DB, userRepo, userRoleRepo)
	storeServ := NewRbacStoreService(storeRepo)

	h := newHandler(deps.Binder, authServ, deps.AuthCfg, menuServ, apiServ, roleServ, userServ, storeServ, authMid)

	return h, authMid
}

// RegisterAdmin 注册 rbac 模块路由（管理端）
func RegisterAdmin(rg *gin.RouterGroup, deps Deps) (*AuthMiddleware, *handler) {
	h, authMid := buildAll(deps)

	rg.POST("/login", h.login)
	rg.POST("/refresh-token", h.refreshToken)

	return authMid, h
}

// RegisterAdminProtected 注册需要 JWT + 权限保护的路由
func RegisterAdminProtected(rg *gin.RouterGroup, h *handler) {
	rg.POST("/logout", h.logout)
	rg.POST("/change-password", h.changePassword)
	rg.POST("/permissions", h.permissions)

	rg.POST("/rbac/menu/list", h.menuList)
	rg.POST("/rbac/menu/create", h.menuCreate)
	rg.POST("/rbac/menu/update", h.menuUpdate)
	rg.POST("/rbac/menu/delete", h.menuDelete)
	rg.POST("/rbac/menu/sync", h.menuSync)
	rg.POST("/rbac/menu/api/list", h.menuApiList)
	rg.POST("/rbac/menu/api/save", h.menuApiSave)

	rg.POST("/rbac/api/list", h.apiList)
	rg.POST("/rbac/api/create", h.apiCreate)
	rg.POST("/rbac/api/update", h.apiUpdate)
	rg.POST("/rbac/api/delete", h.apiDelete)

	rg.POST("/rbac/store/list", h.storeList)
	rg.POST("/rbac/store/create", h.storeCreate)
	rg.POST("/rbac/store/update", h.storeUpdate)
	rg.POST("/rbac/store/delete", h.storeDelete)
	rg.POST("/rbac/store/recycle", h.storeRecycle)
	rg.POST("/rbac/store/restore", h.storeRestore)

	rg.POST("/rbac/role/list", h.roleList)
	rg.POST("/rbac/role/create", h.roleCreate)
	rg.POST("/rbac/role/update", h.roleUpdate)
	rg.POST("/rbac/role/delete", h.roleDelete)
	rg.POST("/rbac/role/set-menus", h.roleSetMenus)

	rg.POST("/rbac/user/list", h.userList)
	rg.POST("/rbac/user/create", h.userCreate)
	rg.POST("/rbac/user/update", h.userUpdate)
	rg.POST("/rbac/user/delete", h.userDelete)
	rg.POST("/rbac/user/set-roles", h.userSetRoles)
	rg.POST("/rbac/user/reset-password", h.userResetPassword)
}
