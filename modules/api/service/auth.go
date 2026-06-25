package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"zero-backend/internal/config"
	"zero-backend/internal/constants"
	"zero-backend/internal/dto"
	"zero-backend/modules/rbac"
	"zero-backend/internal/model"
	"zero-backend/internal/repository"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/baserepo"
	"github.com/241x/zero-kit/helper"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// AuthService 用户认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	cfg      config.ApiAuthConfig
	rdb      *redis.Client
}

// NewAuthService 创建AuthService实例
func NewAuthService(
	userRepo *repository.UserRepository,
	cfg config.ApiAuthConfig,
	rdb *redis.Client,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		rdb:      rdb,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *rbac.AuthLoginRequest) (*dto.UserLoginResponse, string, error) {
	filter := repository.UserFilterField{Username: req.Username}
	item, err := s.userRepo.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return nil, "", apperror.New(errcode.InvalidInput, apperror.WithMsg("用户名或密码错误"))
		}
		return nil, "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("登录失败"))
	}

	ok, err := helper.CheckPassword(req.Password, item.Password)
	if err != nil {
		return nil, "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("验证密码失败"))
	}
	if !ok {
		return nil, "", apperror.New(errcode.InvalidInput, apperror.WithMsg("用户名或密码错误"))
	}

	refreshToken, err := s.getRefreshToken()
	if err != nil {
		return nil, "", err
	}

	itemBytes, err := json.Marshal(item)
	if err != nil {
		return nil, "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("登录失败"))
	}

	result := s.rdb.Set(ctx,
		fmt.Sprintf("%s:%s", constants.RedisUserRefreshTokenKey, refreshToken),
		itemBytes,
		time.Duration(s.cfg.RefreshTokenTtl)*time.Second)

	if result.Err() != nil {
		return nil, "", apperror.Wrap(errcode.Internal, result.Err(), apperror.WithMsg("登录失败"))
	}

	tokenString, err := s.getAccessToken(item)
	if err != nil {
		return nil, "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("登录失败"))
	}

	return &dto.UserLoginResponse{
		Token: tokenString,
		Ttl:   s.cfg.AccessTokenTtl,
		User:  item,
	}, refreshToken, nil
}

// RefreshToken 刷新用户Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.UserLoginResponse, error) {
	itemBytes, err := s.rdb.Get(ctx,
		fmt.Sprintf("%s:%s", constants.RedisUserRefreshTokenKey, refreshToken)).Bytes()

	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	item := &model.User{}
	if err := json.Unmarshal(itemBytes, item); err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	token, err := s.getAccessToken(item)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	return &dto.UserLoginResponse{
		Token: token,
		Ttl:   s.cfg.AccessTokenTtl,
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
		"exp":      time.Now().Add(time.Duration(s.cfg.AccessTokenTtl) * time.Second).Unix(),
		"user_id":  item.ID,
		"store_id": item.StoreId,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.HmacSecret))
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
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return nil, apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取用户信息失败"))
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
func (s *AuthService) ChangePassword(ctx context.Context, req *rbac.ChangePasswordRequest) error {
	// 1. 获取当前用户
	user := ctxkeys.User(ctx).(*model.User)
	if user == nil || user.ID == 0 {
		return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
	}

	// 2. 获取用户信息
	user, err := s.userRepo.FindOne(ctx, user.ID)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}

	// 3. 验证旧密码
	ok, err := helper.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}
	if !ok {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("旧密码不正确"))
	}

	// 4. 加密新密码
	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}

	// 5. 更新密码
	updateData := map[string]any{
		"password": hashedPassword,
	}

	if err := s.userRepo.Updates(ctx, user, updateData); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}

	return nil
}
