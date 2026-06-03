package request

import (
	"errors"
	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/errcode"
	"zero-backend/internal/model"
	"zero-backend/pkg/apperror"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewRequest(validate *validator.Validate, trans ut.Translator) *Request {
	return &Request{
		validate: validate,
		trans:    trans,
	}
}

// ShouldBindJSON 验证并绑定数据
func (r *Request) ShouldBindJSON(ctx *gin.Context, data any) error {
	// 将Json数据绑定到变量
	if err := ctx.ShouldBindJSON(data); err != nil {
		return apperror.New(errcode.InvalidInput, apperror.WithCause(err))
	}

	// 对传入参数进行验证
	if err := r.validate.Struct(data); err != nil {
		if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
			return apperror.New(errcode.InvalidInput, apperror.WithMsg(validationErrors[0].Translate(r.trans)))
		}

		return apperror.Wrap(errcode.Internal, err)
	}

	return nil
}

// ShouldBindJSONArray 验证并绑定数据
func (r *Request) ShouldBindJSONArray(ctx *gin.Context, data any) error {
	// 将Json数据绑定到变量
	if err := ctx.ShouldBindJSON(data); err != nil {
		return apperror.New(errcode.InvalidInput, apperror.WithCause(err))
	}

	// 对传入参数进行验证
	list, _ := data.([]any)
	for _, item := range list {
		if err := r.validate.Struct(item); err != nil {
			if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
				return apperror.New(errcode.InvalidInput, apperror.WithMsg(validationErrors[0].Translate(r.trans)))
			}

			return apperror.Wrap(errcode.Internal, err)
		}
	}

	return nil
}

// GetStoreId 获取企业ID
func GetStoreId(ctx *gin.Context) uint32 {
	return ctxkeys.StoreID(ctx.Request.Context())
}

// IsSuperUser 是否是超级管理员
func IsSuperUser(ctx *gin.Context) bool {
	if user, ok := ctxkeys.User(ctx.Request.Context()).(*model.RbacUser); ok {
		return user.SU
	}

	return false
}

// GetUserID 获取用户ID
func GetUserID(ctx *gin.Context) uint32 {
	if user, ok := ctxkeys.User(ctx.Request.Context()).(*model.RbacUser); ok {
		return user.ID
	}

	if user, ok := ctxkeys.User(ctx.Request.Context()).(*model.User); ok {
		return user.ID
	}

	return 0
}
