package request

import (
	"errors"
	"zero-backend/internal/apperror"
	"zero-backend/internal/ctxkeys"
	"zero-backend/internal/model"

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
		return apperror.NewSystemError(err, "传入参数错误")
	}

	// 对传入参数进行验证
	if err := r.validate.Struct(data); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return apperror.NewUserError(validationErrors[0].Translate(r.trans))
		}

		return apperror.NewSystemError(err, "参数验证错误")
	}

	return nil
}

// ShouldBindJSONArray 验证并绑定数据
func (r *Request) ShouldBindJSONArray(ctx *gin.Context, data any) error {
	// 将Json数据绑定到变量
	if err := ctx.ShouldBindJSON(data); err != nil {
		return apperror.NewSystemError(err, "传入参数错误")
	}

	// 对传入参数进行验证
	list, _ := data.([]any)
	for _, item := range list {
		if err := r.validate.Struct(item); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				return apperror.NewUserError(validationErrors[0].Translate(r.trans))
			}

			return apperror.NewSystemError(err, "参数验证错误")
		}
	}

	return nil
}

// GetStoreId 获取企业ID
func GetStoreId(ctx *gin.Context) uint32 {
	value := ctx.Request.Context().Value(ctxkeys.StoreIdKey{})
	if value == nil {
		return 0
	}

	if storeId, ok := value.(uint32); ok {
		return storeId
	}

	return 0
}

// IsSuperUser 是否是超级管理员
func IsSuperUser(ctx *gin.Context) bool {
	value := ctx.Request.Context().Value(ctxkeys.UserKey{})
	if value == nil {
		return false
	}

	if user, ok := value.(*model.RbacUser); ok {
		return user.SU
	}

	return false
}

// GetUserID 获取用户ID
func GetUserID(ctx *gin.Context) uint32 {
	value := ctx.Request.Context().Value(ctxkeys.UserKey{})
	if value == nil {
		return 0
	}

	if user, ok := value.(*model.RbacUser); ok {
		return user.ID
	}

	if user, ok := value.(*model.User); ok {
		return user.ID
	}

	return 0
}
