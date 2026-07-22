package platform_user

import (
	"context"
	"errors"

	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/baserepo"
	"github.com/241x/zero-kit/helper"
	"github.com/241x/zero-web/errcode"
)

// PlatformUserService 平台用户服务
type PlatformUserService struct {
	repo *PlatformUserRepository
}

// NewPlatformUserService 创建平台用户服务
func NewPlatformUserService(repo *PlatformUserRepository) *PlatformUserService {
	return &PlatformUserService{repo: repo}
}

// FindList 获取平台用户列表
func (s *PlatformUserService) FindList(ctx context.Context, req *UserListRequest) (*ListResult, error) {
	result := &ListResult{List: []*PlatformUser{}, Total: 0}

	filter := &PlatformUserFilterField{
		Username: req.Username,
		RealName: req.RealName,
	}

	pagination := baserepo.NewPagination(req.Page, req.Limit)
	orders := baserepo.Orders{
		{Field: "id", Sort: "desc"},
	}

	total, err := s.repo.Count(ctx, filter, baserepo.WithScopes(nil))
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取用户列表失败"))
	}
	if total == 0 {
		return result, nil
	}
	result.Total = total

	list, err := s.repo.FindAll(ctx, filter, pagination, orders, baserepo.WithScopes(nil))
	if err != nil {
		return nil, apperror.Wrap(errcode.Internal, err, apperror.WithMsg("获取用户列表失败"))
	}
	result.List = list
	return result, nil
}

// Create 创建平台用户
func (s *PlatformUserService) Create(ctx context.Context, req *UserCreateRequest) error {
	if err := s.checkUsername(ctx, req.Username); err != nil {
		return err
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("密码加密失败"))
	}

	item := &PlatformUser{
		Username: req.Username,
		Password: hashedPassword,
		RealName: req.RealName,
		Role:     req.Role,
		AvatarID: req.AvatarID,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("创建用户失败"))
	}
	return nil
}

// Update 更新平台用户
func (s *PlatformUserService) Update(ctx context.Context, req *UserUpdateRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, baserepo.WithScopes(nil))
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新用户失败"))
	}

	if item.Username != req.Username {
		if err := s.checkUsername(ctx, req.Username); err != nil {
			return err
		}
	}

	updateData := map[string]any{
		"username":  req.Username,
		"real_name": req.RealName,
		"role":      req.Role,
		"status":    req.Status,
		"avatar_id": req.AvatarID,
	}

	if req.Password != "" {
		hashedPassword, err := helper.HashPassword(req.Password)
		if err != nil {
			return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("密码加密失败"))
		}
		updateData["password"] = hashedPassword
	}

	if err := s.repo.Updates(ctx, item, updateData); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("更新用户失败"))
	}
	return nil
}

// Delete 删除平台用户
func (s *PlatformUserService) Delete(ctx context.Context, req *UserDeleteRequest) error {
	item, err := s.repo.FindOne(ctx, req.ID, baserepo.WithScopes(nil))
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除用户失败"))
	}

	if err := s.repo.Delete(ctx, item.ID); err != nil {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("删除用户失败"))
	}
	return nil
}

// ResetPassword 重置平台用户密码
func (s *PlatformUserService) ResetPassword(ctx context.Context, req *UserResetPasswordRequest) (string, error) {
	item, err := s.repo.FindOne(ctx, req.ID, baserepo.WithScopes(nil))
	if err != nil {
		if errors.Is(err, baserepo.ErrRecordNotFound) {
			return "", apperror.New(errcode.NotFound, apperror.WithMsg("用户不存在"))
		}
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("重置密码失败"))
	}

	newPassword := helper.RandomStringWithSymbols(12)
	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("密码加密失败"))
	}

	if err := s.repo.Updates(ctx, item, map[string]any{"password": hashedPassword}); err != nil {
		return "", apperror.Wrap(errcode.Internal, err, apperror.WithMsg("重置密码失败"))
	}
	return newPassword, nil
}

// checkUsername 检查用户名是否已存在
func (s *PlatformUserService) checkUsername(ctx context.Context, username string) error {
	filter := &PlatformUserUsernameFilterField{Username: username}
	_, err := s.repo.FindOne(ctx, filter, baserepo.WithScopes(nil))
	if err == nil {
		return apperror.New(errcode.InvalidInput, apperror.WithMsg("用户名已存在"))
	}
	if !errors.Is(err, baserepo.ErrRecordNotFound) {
		return apperror.Wrap(errcode.Internal, err, apperror.WithMsg("检查用户名失败"))
	}
	return nil
}
