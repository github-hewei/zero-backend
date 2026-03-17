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

// AuthService 用户认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
	rdb      *redis.Client
}

// NewAuthService 创建AuthService实例
func NewAuthService(
	userRepo *repository.UserRepository,
	cfg *config.Config,
	rdb *redis.Client,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		rdb:      rdb,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *dto.AuthLoginRequest) (*dto.UserLoginResponse, string, error) {
	filter := repository.UserFilterField{Username: req.Username}
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
		fmt.Sprintf("%s:%s", constants.RedisUserRefreshTokenKey, refreshToken),
		itemBytes,
		time.Duration(s.cfg.Api.RefreshTokenTtl)*time.Second)

	if result.Err() != nil {
		return nil, "", apperror.NewSystemError(result.Err(), "保存token失败")
	}

	tokenString, err := s.getAccessToken(item)
	if err != nil {
		return nil, "", apperror.NewSystemError(err, "生成token失败")
	}

	return &dto.UserLoginResponse{
		Token: tokenString,
		Ttl:   s.cfg.Api.AccessTokenTtl,
		User:  item,
	}, refreshToken, nil
}

// RefreshToken 刷新用户Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.UserLoginResponse, error) {
	itemBytes, err := s.rdb.Get(ctx,
		fmt.Sprintf("%s:%s", constants.RedisUserRefreshTokenKey, refreshToken)).Bytes()

	if err != nil {
		return nil, apperror.NewSystemError(err, "登录已过期，请重新登录")
	}

	item := &model.User{}
	if err := json.Unmarshal(itemBytes, item); err != nil {
		return nil, apperror.NewSystemError(err, "登录已过期，请重新登录")
	}

	token, err := s.getAccessToken(item)
	if err != nil {
		return nil, apperror.NewSystemError(err, "登录已过期，请重新登录")
	}

	return &dto.UserLoginResponse{
		Token: token,
		Ttl:   s.cfg.Api.AccessTokenTtl,
		User:  nil,
	}, nil
}

// getRefreshToken 获取刷新令牌
func (s *AuthService) getRefreshToken() (string, error) {
	token := helper.StringMd5(fmt.Sprintf("%d", time.Now().UnixNano()) + helper.RandomString(16))
	return token, nil
}

// getAccessToken 获取用户访问令牌
func (s *AuthService) getAccessToken(item *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Duration(s.cfg.Api.AccessTokenTtl) * time.Second).Unix(),
		"user_id":  item.ID,
		"store_id": item.StoreId,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.Api.HmacSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, userId uint32) (*model.User, error) {
	// 1. 先从Redis缓存获取
	cacheKey := fmt.Sprintf("%s:%d", constants.RedisUserLoginKey, userId)
	result := s.rdb.Get(ctx, cacheKey)

	// 2. 缓存命中
	if result.Err() == nil {
		itemBytes, err := result.Bytes()
		if err == nil {
			item := &model.User{}
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

	// 4. 写入缓存
	userBytes, err := json.Marshal(user)
	if err != nil {
		return user, nil // 即使序列化失败也返回用户数据
	}

	// 设置1小时有效期
	s.rdb.Set(ctx, cacheKey, userBytes, time.Hour)

	return user, nil
}

// ChangePassword 修改用户密码
func (s *AuthService) ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error {
	// 1. 获取当前用户
	user := ctx.Value(ctxkeys.UserKey{}).(*model.User)
	if user == nil || user.ID == 0 {
		return apperror.NewUserError("用户不存在")
	}

	// 2. 获取用户信息
	user, err := s.userRepo.FindOne(ctx, user.ID)
	if err != nil {
		return apperror.NewSystemError(err, "查询用户信息失败")
	}
	if user == nil || user.ID == 0 {
		return apperror.NewUserError("用户不存在")
	}

	// 3. 验证旧密码
	if ok := helper.CheckPassword(req.OldPassword, user.Password); !ok {
		return apperror.NewUserError("旧密码不正确")
	}

	// 4. 加密新密码
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.NewSystemError(err, "密码加密失败")
	}

	// 5. 更新密码
	updateData := map[string]any{
		"password": hashedPassword,
	}

	if err := s.userRepo.Updates(ctx, user, updateData); err != nil {
		return apperror.NewSystemError(err, "密码更新失败")
	}

	return nil
}
