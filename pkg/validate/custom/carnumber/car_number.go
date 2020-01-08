package carnumber

import (
	"regexp"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func ValidateCarNumber(fl validator.FieldLevel) bool {
	b := ValidCarNumber(fl.Field().String())
	return b
}

func ValidateCarNumberTranslator(ut ut.Translator) (err error) {
	return ut.Add("carNumber", "车牌号格式错误", true)
}

func ValidCarNumber(CarNumber string) bool {
	carNumberReg := "^([京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领A-Z]{1}[A-Z]{1}(([0-9]{5}[DF])|([DF]([A-HJ-NP-Z0-9])[0-9]{4})))|([京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领A-Z]{1}[A-Z]{1}[A-HJ-NP-Z0-9]{4}[A-HJ-NP-Z0-9挂学警港澳]{1})$"
	reg := regexp.MustCompile(carNumberReg)
	return reg.MatchString(CarNumber)
}
