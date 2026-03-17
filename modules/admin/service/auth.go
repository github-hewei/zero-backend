package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zero-backend/internal/apperror"
	"zero-backend/internal/config"
	"zero-backend/internal/constants"
	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/dto"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"
	"zero-backend/pkg/helper"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// AuthService 管理员认证服务
type AuthService struct {
	userRepo     *repository.RbacUserRepository
	apiRepo      *repository.RbacApiRepository
	roleRepo     *repository.RbacRoleRepository
	menuRepo     *repository.RbacMenuRepository
	roleMenuRepo *repository.RbacRoleMenuRepository
	userRoleRepo *repository.RbacUserRoleRepository
	menuApiRepo  *repository.RbacMenuApiRepository
	cfg          *config.Config
	rdb          *redis.Client
}

// NewAuthService 创建AuthService实例
func NewAuthService(
	userRepo *repository.RbacUserRepository,
	apiRepo *repository.RbacApiRepository,
	roleRepo *repository.RbacRoleRepository,
	menuRepo *repository.RbacMenuRepository,
	roleMenuRepo *repository.RbacRoleMenuRepository,
	userRoleRepo *repository.RbacUserRoleRepository,
	menuApiRepo *repository.RbacMenuApiRepository,
	cfg *config.Config,
	rdb *redis.Client,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		apiRepo:      apiRepo,
		roleRepo:     roleRepo,
		menuRepo:     menuRepo,
		roleMenuRepo: roleMenuRepo,
		userRoleRepo: userRoleRepo,
		menuApiRepo:  menuApiRepo,
		cfg:          cfg,
		rdb:          rdb,
	}
}

// Login 系统登录
func (s *AuthService) Login(ctx context.Context, req *dto.AuthLoginRequest) (*dto.AdminLoginResponse, string, error) {
	filter := &repository.RbacUserUsernameFilterField{Username: req.Username}
	item, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		return nil, "", err
	}

	if item.ID == 0 {
		return nil, "", apperror.NewUserError("用户名或密码错误")
	}

	if ok := helper.CheckPassword(req.Password, item.Password); !ok {
		return nil, "", apperror.NewUserError("用户名或密码错误")
	}

	refreshToken, err := s.getRefreshToken()
	if err != nil {
		return nil, "", err
	}

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return nil, "", apperror.NewSystemError(err, "序列化用户信息失败")
	}

	result := s.rdb.Set(ctx,
		fmt.Sprintf("%s:%s", constants.RedisAdminRefreshTokenKey, refreshToken),
		itemBytes,
		time.Duration(s.cfg.Admin.RefreshTokenTtl)*time.Second)

	if result.Err() != nil {
		return nil, "", apperror.NewSystemError(result.Err(), "保存token失败")
	}

	tokenString, err := s.getAccessToken(item)
	if err != nil {
		return nil, "", apperror.NewSystemError(err, "生成token失败")
	}

	s.WithSU(item)

	return &dto.AdminLoginResponse{
		Token: tokenString,
		Ttl:   s.cfg.Admin.RefreshTokenTtl,
		User:  item,
	}, refreshToken, nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.AdminLoginResponse, error) {
	itemBytes, err := s.rdb.Get(ctx,
		fmt.Sprintf("%s:%s", constants.RedisAdminRefreshTokenKey, refreshToken)).Bytes()

	if err != nil {
		return nil, apperror.NewSystemError(err, "登录已过期，请重新登录")
	}

	item := &model.RbacUser{}
	if err := json.Unmarshal(itemBytes, item); err != nil {
		return nil, apperror.NewSystemError(err, "登录已过期，请重新登录")
	}

	token, err := s.getAccessToken(item)
	if err != nil {
		return nil, apperror.NewSystemError(err, "登录已过期，请重新登录")
	}

	return &dto.AdminLoginResponse{
		Token: token,
		Ttl:   s.cfg.Admin.AccessTokenTtl,
		User:  nil,
	}, nil
}

// getRefreshToken 获取刷新令牌
func (s *AuthService) getRefreshToken() (string, error) {
	token := helper.StringMd5(fmt.Sprintf("%d", time.Now().UnixNano()) + helper.RandomString(16))
	return token, nil
}

// getAccessToken 获取访问令牌
func (s *AuthService) getAccessToken(item *model.RbacUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Duration(s.cfg.Admin.AccessTokenTtl) * time.Second).Unix(),
		"user_id":  item.ID,
		"store_id": item.StoreId,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.Admin.HmacSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, userId uint32) (*model.RbacUser, error) {
	// 1. 先从Redis缓存获取
	cacheKey := fmt.Sprintf("%s:%d", constants.RedisAdminLoginKey, userId)
	result := s.rdb.Get(ctx, cacheKey)

	// 2. 缓存命中
	if result.Err() == nil {
		itemBytes, err := result.Bytes()
		if err == nil {
			item := &model.RbacUser{}
			if err := json.Unmarshal(itemBytes, item); err == nil {
				return item, nil
			}
		}
	}

	// 3. 缓存未命中，查询数据库
	user, err := s.userRepo.FindOne(ctx, userId)
	if err != nil {
		return nil, apperror.NewSystemError(err, "查询用户信息失败")
	}
	if user == nil || user.ID == 0 {
		return nil, apperror.NewUserError("用户不存在")
	}

	s.WithSU(user)
	// 4. 写入缓存
	userBytes, err := json.Marshal(user)
	if err != nil {
		return user, nil // 即使序列化失败也返回用户数据
	}

	// 设置1小时有效期
	s.rdb.Set(ctx, cacheKey, userBytes, time.Hour)

	return user, nil
}

// WithSU 设置是否为超级管理员
func (s *AuthService) WithSU(user *model.RbacUser) {
	if int(user.ID) == s.cfg.Admin.SuperUserId {
		user.SU = true
	}
}

// GetPermissions 获取用户菜单权限
func (s *AuthService) GetPermissions(ctx context.Context, req *dto.AuthGetPermissionsRequest) ([]*model.RbacMenu, error) {
	user := ctx.Value(ctxkeys.UserKey{}).(*model.RbacUser)
	if user == nil {
		return nil, nil
	}

	// 1. 如果是超级管理员，返回全部菜单
	if user.SU {
		menus, err := s.menuRepo.FindAll(ctx, nil, nil, nil, repository.WithScopes(nil))
		if err != nil {
			return nil, apperror.NewSystemError(err, "获取菜单失败")
		}
		if req.IsTree {
			return model.RbacMenuList(menus).Tree(), nil
		}
		return menus, nil
	}

	// 2. 获取用户角色
	roles, err := s.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, apperror.NewSystemError(err, "获取用户角色失败")
	}
	if len(roles) == 0 {
		return nil, nil
	}

	// 3. 获取角色对应的菜单
	var menus []*model.RbacMenu
	for _, role := range roles {
		roleMenus, err := s.GetRoleMenus(ctx, role.ID)
		if err != nil {
			return nil, apperror.NewSystemError(err, "获取角色菜单失败")
		}
		menus = append(menus, roleMenus...)
	}

	// 4. 去重
	uniqueMenus := make(map[uint32]*model.RbacMenu)
	for _, menu := range menus {
		uniqueMenus[menu.ID] = menu
	}

	// 5. 转换为切片返回
	var result []*model.RbacMenu
	for _, menu := range uniqueMenus {
		result = append(result, menu)
	}

	if req.IsTree {
		return model.RbacMenuList(result).Tree(), nil
	}

	return result, nil
}

// GetUserRoles 根据用户ID获取角色列表
func (s *AuthService) GetUserRoles(ctx context.Context, userID uint32) ([]*model.RbacRole, error) {
	filter := &repository.RbacUserRoleFilterField{UserId: userID}
	userRoles, err := s.userRoleRepo.FindAll(ctx, filter, nil, nil)

	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []*model.RbacRole{}, nil
	}

	var roleIDs []uint32
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleId)
	}

	// 获取角色信息
	roles, err := s.roleRepo.FindAll(ctx, &repository.RbacRoleFilterField{IDs: roleIDs}, nil, nil)

	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetRoleMenus 根据角色ID获取角色菜单
func (s *AuthService) GetRoleMenus(ctx context.Context, roleId uint32) ([]*model.RbacMenu, error) {
	roleMenus, err := s.roleMenuRepo.FindAll(ctx, &repository.RbacRoleMenuFilterField{RoleId: roleId}, nil, nil)

	if err != nil {
		return nil, err
	}

	if len(roleMenus) == 0 {
		return []*model.RbacMenu{}, nil
	}

	// 获取菜单ID列表
	var menuIDs []uint32
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuId)
	}

	// 获取菜单详情
	menus, err := s.menuRepo.FindAll(ctx, &repository.RbacMenuFilterField{IDs: menuIDs}, nil, nil, repository.WithScopes(nil))
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error {
	user, _ := ctx.Value(ctxkeys.UserKey{}).(*model.RbacUser)
	if user == nil || user.ID == 0 {
		return apperror.NewUserError("用户不存在")
	}
	// 1. 获取用户信息
	user, err := s.userRepo.FindOne(ctx, user.ID)
	if err != nil {
		return apperror.NewSystemError(err, "查询用户信息失败")
	}
	if user == nil || user.ID == 0 {
		return apperror.NewUserError("用户不存在")
	}

	// 2. 验证旧密码
	if ok := helper.CheckPassword(req.OldPassword, user.Password); !ok {
		return apperror.NewUserError("旧密码不正确")
	}

	// 3. 加密新密码
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.NewSystemError(err, "密码加密失败")
	}

	// 4. 更新密码
	updateData := map[string]any{
		"password": hashedPassword,
	}

	if err := s.userRepo.Updates(ctx, user, updateData); err != nil {
		return apperror.NewSystemError(err, "密码更新失败")
	}

	return nil
}

// CheckAPIPermission 检查API权限
func (s *AuthService) CheckAPIPermission(ctx context.Context, user *model.RbacUser, apiPath string) (bool, error) {
	if user.SU {
		return true, nil
	}

	// 1. 获取用户角色
	roles, err := s.GetUserRoles(ctx, user.ID)
	if err != nil {
		return false, err
	}
	if len(roles) == 0 {
		return false, nil
	}

	// 2. 获取API权限
	api, err := s.apiRepo.GetAPIByPath(ctx, apiPath)
	if err != nil {
		return false, err
	}
	if api == nil || api.ID == 0 {
		return false, nil
	}

	// 3. 检查角色是否有该API权限
	for _, role := range roles {
		hasPerm, err := s.CheckRoleAPIPermission(ctx, role.ID, api.ID)
		if err != nil {
			return false, err
		}
		if hasPerm {
			return true, nil
		}
	}

	return false, nil
}

// CheckRoleAPIPermission 检查角色是否有API权限
func (s *AuthService) CheckRoleAPIPermission(ctx context.Context, roleID uint32, apiID uint32) (bool, error) {
	roleMenus, err := s.roleMenuRepo.FindAll(ctx, &repository.RbacRoleMenuFilterField{RoleId: roleID}, nil, nil)

	if err != nil {
		return false, err
	}

	if len(roleMenus) == 0 {
		return false, nil
	}

	// 获取菜单API关联
	var menuIDs []uint32
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuId)
	}

	filter := &repository.RbacMenuApiFilterField{MenuIDs: menuIDs, ApiId: apiID}
	menuApis, err := s.menuApiRepo.FindAll(ctx, filter, nil, nil, repository.WithScopes(nil))
	if err != nil {
		return false, err
	}

	return len(menuApis) > 0, nil
}
