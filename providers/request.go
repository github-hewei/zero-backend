package providers

import (
	"zero-backend/internal/errcode"
	"zero-backend/pkg/bind"

	"github.com/241x/zero-kit/apperror"
	"github.com/google/wire"
)

// RequestProviderSet 提供请求相关依赖集合。
var RequestProviderSet = wire.NewSet(
	bind.New,
	bind.NewValidate,
	bind.NewTrans,
	ProvideBindErrCode,
)

func ProvideBindErrCode() apperror.Code {
	return errcode.InvalidInput
}
