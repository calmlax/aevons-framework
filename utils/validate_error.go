package utils

import (
	"errors"
	"strings"

	errdef "github.com/calmlax/aevons-framework/err"

	"github.com/go-playground/validator/v10"
)

// GetValidateErrorKey 解析校验错误 → 返回多语言key
func GetValidateErrorKey(err error) (*errdef.ErrorDef, string) {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return &errdef.ErrInvalidBody, ""
	}

	var errDef errdef.ErrorDef
	// 取第一个错误
	e := validationErrs[0]
	field := strings.ToLower(e.Field()) // 转小驼峰给前端用
	switch e.Tag() {
	case "required":
		errDef = errdef.ErrValidationRequired
	case "max":
		errDef = errdef.ErrValidationMax
	case "min":
		errDef = errdef.ErrValidationMin
	case "oneof":
		errDef = errdef.ErrValidationOneof
	case "email":
		errDef = errdef.ErrValidationEmail
	case "phone":
		errDef = errdef.ErrValidationPhone
	default:
		errDef = errdef.ErrInvalidBody
	}

	return &errDef, field
}
