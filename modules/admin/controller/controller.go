package controller

// Controllers 控制器集合
type Controllers struct {
	AuthController            *AuthController
	RbacMenuController        *RbacMenuController
	RbacApiController         *RbacApiController
	RbacRoleController        *RbacRoleController
	RbacUserController        *RbacUserController
	RbacStoreController       *RbacStoreController
	HealthController          *HealthController
}
