package utils

import (
	"errors"
	"strings"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/go-playground/validator/v10"
)

// GetValidateErrorKey 解析校验错误 → 返回多语言key
func GetValidateErrorKey(err error) (*apperr.ErrorDef, string) {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return &apperr.ErrInvalidBody, ""
	}

	var errDef apperr.ErrorDef
	// 取第一个错误
	e := validationErrs[0]
	field := strings.ToLower(e.Field()) // 转小驼峰给前端用
	switch e.Tag() {
	case "required":
		errDef = apperr.ErrValidationRequired
	case "max":
		errDef = apperr.ErrValidationMax
	case "min":
		errDef = apperr.ErrValidationMin
	case "oneof":
		errDef = apperr.ErrValidationOneof
	case "email":
		errDef = apperr.ErrValidationEmail
	case "phone":
		errDef = apperr.ErrValidationPhone
	default:
		errDef = apperr.ErrInvalidBody
	}

	return &errDef, field
}
