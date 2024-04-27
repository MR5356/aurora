package validate

import (
	"errors"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"strings"
)

var validate *validator.Validate

var zh = zh_Hans.New()
var uni = ut.New(zh, zh)
var trans, _ = uni.GetTranslator("zh")

func init() {
	validate = validator.New()
	_ = zh_translations.RegisterDefaultTranslations(validate, trans)
}

// Validate 参数校验
func Validate(obj any) error {
	err := validate.Struct(obj)
	if err != nil {
		var msg []string
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			msg = append(msg, e.Translate(trans))
		}
		return errors.New(strings.Join(msg, ", "))
	}
	return nil
}
