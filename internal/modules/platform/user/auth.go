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

const (
	redisPlatformRefreshTokenKey = "ZAG:PLATFORM:REFRESH:TOKEN"
	redisPlatformLoginKey        = "ZAG:PLATFORM:LOGIN"
)

// CaptchaVerifier 验证码验证接口（由宿主注入）
type CaptchaVerifier interface {
	Verify(ctx context.Context, captchaID, captchaCode string) error
}

// AuthService 平台认证服务
type AuthService struct {
	repo    *PlatformUserRepository
	cfg     Config
	rdb     *redis.Client
	captcha CaptchaVerifier
}

// NewAuthService 创建平台认证服务
func NewAuthService(repo *PlatformUserRepository, cfg Config, rdb *redis.Client, captcha CaptchaVerifier) *AuthService {
	return &AuthService{repo: repo, cfg: cfg, rdb: rdb, captcha: captcha}
}

// Login 平台登录
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, string, error) {
	if err := s.captcha.Verify(ctx, req.CaptchaID, req.CaptchaCode); err != nil {
		return nil, "", err
	}

	filter := &PlatformUserUsernameFilterField{Username: req.Username}
	item, err := s.repo.FindOne(ctx, filter, baserepo.WithScopes(nil))
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return nil, "", apperror.New(errcode.InvalidInput, apperror.WithMsg("用户名或密码错误"))
		}
		return nil, "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("登录失败"))
	}

	if item.Status != 1 {
		return nil, "", apperror.New(errcode.Forbidden, apperror.WithMsg("账号已被禁用"))
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
		fmt.Sprintf("%s:%s", redisPlatformRefreshTokenKey, refreshToken),
		itemBytes,
		time.Duration(s.cfg.RefreshTokenTtl)*time.Second)
	if result.Err() != nil {
		return nil, "", apperror.Wrap(errcode.Internal, result.Err(), apperror.WithMsg("登录失败"))
	}

	tokenString, err := s.getAccessToken(item)
	if err != nil {
		return nil, "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("登录失败"))
	}

	return &LoginResponse{
		Token: tokenString,
		Ttl:   s.cfg.RefreshTokenTtl,
		User:  item,
	}, refreshToken, nil
}

// RefreshToken 刷新Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	itemBytes, err := s.rdb.Get(ctx,
		fmt.Sprintf("%s:%s", redisPlatformRefreshTokenKey, refreshToken)).Bytes()
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	item := &PlatformUser{}
	if err := json.Unmarshal(itemBytes, item); err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	token, err := s.getAccessToken(item)
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("刷新令牌失败"))
	}

	return &LoginResponse{
		Token: token,
		Ttl:   s.cfg.AccessTokenTtl,
		User:  nil,
	}, nil
}

// getRefreshToken 获取刷新Token
func (s *AuthService) getRefreshToken() (string, error) {
	token := helper.StringMd5(fmt.Sprintf("%d", time.Now().UnixNano()) + helper.RandomString(16))
	return token, nil
}

// getAccessToken 获取访问Token
func (s *AuthService) getAccessToken(item *PlatformUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Duration(s.cfg.AccessTokenTtl) * time.Second).Unix(),
		"user_id": item.ID,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.HmacSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, userId uint32) (*PlatformUser, error) {
	cacheKey := fmt.Sprintf("%s:%d", redisPlatformLoginKey, userId)
	result := s.rdb.Get(ctx, cacheKey)

	if result.Err() == nil {
		itemBytes, err := result.Bytes()
		if err == nil {
			item := &PlatformUser{}
			if err := json.Unmarshal(itemBytes, item); err == nil {
				return item, nil
			}
		}
	}

	user, err := s.repo.FindOne(ctx, userId, baserepo.WithScopes(nil))
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

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	user, _ := ctxkeys.User(ctx).(*PlatformUser)
	if user == nil || user.ID == 0 {
		return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
	}

	user, err := s.repo.FindOne(ctx, user.ID, baserepo.WithScopes(nil))
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
	if err := s.repo.Updates(ctx, user, updateData); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("修改密码失败"))
	}
	return nil
}
