package request

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
)

// NewValidate 创建验证器
func NewValidate() *validator.Validate {
	return validator.New()
}

// NewTrans 创建翻译器
func NewTrans(v *validator.Validate) ut.Translator {
	zt := zh.New()
	uni := ut.New(zt, zt)
	trans, _ := uni.GetTranslator("zh")
	_ = translations.RegisterDefaultTranslations(v, trans)
	return trans
}
