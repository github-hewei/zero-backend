package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/baserepo"
	"github.com/241x/zero-kit/helper"
	"github.com/241x/zero-web/ctxkeys"
	"github.com/241x/zero-web/errcode"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// AuthServiceInterface 用户认证服务接口
type AuthServiceInterface interface {
	Login(ctx context.Context, req *AuthLoginRequest) (*UserLoginResponse, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (*UserLoginResponse, error)
	ChangePassword(ctx context.Context, req *ChangePasswordRequest) error
	GetUserInfo(ctx context.Context, userId uint32) (*User, error)
}

const (
	redisUserLoginKey        = "ZAG:USER:LOGIN"
	redisUserRefreshTokenKey = "ZAG:USER:REFRESH:TOKEN"
)

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token string `json:"token"`
	Ttl   int    `json:"ttl"`
	User  *User  `json:"user,omitempty"`
}

// AuthService 用户认证服务
type AuthService struct {
	userRepo *Repository
	cfg      Config
	rdb      *redis.Client
}

// NewAuthService 创建AuthService实例
func NewAuthService(
	userRepo *Repository,
	cfg Config,
	rdb *redis.Client,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
		rdb:      rdb,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *AuthLoginRequest) (*UserLoginResponse, string, error) {
	filter := Filter{Username: req.Username}
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
		fmt.Sprintf("%s:%s", redisUserRefreshTokenKey, refreshToken),
		itemBytes,
		time.Duration(s.cfg.RefreshTokenTtl)*time.Second)
	if result.Err() != nil {
		return nil, "", apperror.Wrap(errcode.Internal, result.Err(), apperror.WithMsg("登录失败"))
	}

	tokenString, err := s.getAccessToken(item)
	if err != nil {
		return nil, "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("登录失败"))
	}

	return &UserLoginResponse{
		Token: tokenString,
		Ttl:   s.cfg.AccessTokenTtl,
		User:  item,
	}, refreshToken, nil
}

// RefreshToken 刷新用户Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*UserLoginResponse, error) {
	itemBytes, err := s.rdb.Get(ctx,
		fmt.Sprintf("%s:%s", redisUserRefreshTokenKey, refreshToken)).Bytes()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	item := &User{}
	if err := json.Unmarshal(itemBytes, item); err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	token, err := s.getAccessToken(item)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	return &UserLoginResponse{
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

// getAccessToken 获取访问令牌
func (s *AuthService) getAccessToken(item *User) (string, error) {
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
func (s *AuthService) GetUserInfo(ctx context.Context, userId uint32) (*User, error) {
	cacheKey := fmt.Sprintf("%s:%d", redisUserLoginKey, userId)
	result := s.rdb.Get(ctx, cacheKey)

	if result.Err() == nil {
		itemBytes, err := result.Bytes()
		if err == nil {
			item := &User{}
			if err := json.Unmarshal(itemBytes, item); err == nil {
				return item, nil
			}
		}
	}

	user, err := s.userRepo.FindOne(ctx, userId)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return nil, apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取用户信息失败"))
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return user, nil
	}
	s.rdb.Set(ctx, cacheKey, userBytes, time.Hour)

	return user, nil
}

// ChangePassword 修改用户密码
func (s *AuthService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	user := ctxkeys.User(ctx).(*User)
	if user == nil || user.ID == 0 {
		return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
	}

	user, err := s.userRepo.FindOne(ctx, user.ID)
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}

	ok, err := helper.CheckPassword(req.OldPassword, user.Password)
	if err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("验证密码失败"))
	}
	if !ok {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("旧密码不正确"))
	}

	hashedPassword, err := helper.HashPassword(req.NewPassword)
	if err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}

	updateData := map[string]any{"password": hashedPassword}
	if err := s.userRepo.Updates(ctx, user, updateData); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}
	return nil
}
