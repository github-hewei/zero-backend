package providers

import (
	"github.com/241x/zero-kit/apperror"
	"github.com/241x/zero-kit/bind"
	"github.com/241x/zero-web/errcode"
	"github.com/google/wire"
)

// BindProviderSet 提供请求相关依赖集合。
var BindProviderSet = wire.NewSet(
	bind.New,
	bind.NewValidate,
	bind.NewTrans,
	ProvideBindErrCode,
)

func ProvideBindErrCode() apperror.Code {
	return errcode.InvalidInput
}
