package idcard

import (
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

//
func ValidateIdCardNo(fl validator.FieldLevel) bool {
	citizenNo := fl.Field().Bytes()
	b := IsValidCitizenNo18(&citizenNo)
	return b
}

// ValidateTimeTranslator 翻译
func ValidateIdCardNoTranslator(ut ut.Translator) (err error) {
	return ut.Add("idCard", "{0} 身份证格式错误", true)
}
