package errors

import (
	stderrors "errors"
	"fmt"
	"net/http"
)

var (
	// ErrNotFound 表示目标资源不存在。
	ErrNotFound = stderrors.New("resource not found")
	// ErrAlreadyExists 表示目标资源已存在。
	ErrAlreadyExists = stderrors.New("resource already exists")
	// ErrUnauthorized 表示请求未认证。
	ErrUnauthorized = stderrors.New("unauthorized")
	// ErrForbidden 表示请求无权限执行。
	ErrForbidden = stderrors.New("forbidden")
	// ErrConflict 表示请求与当前资源状态冲突。
	ErrConflict = stderrors.New("resource conflict")
	// ErrInternal 表示系统内部错误。
	ErrInternal = stderrors.New("internal server error")
)

// 兼容旧命名，避免现有调用点一次性全改。
var (
	ErrorNotFound      = ErrNotFound
	ErrorNoUpdateField = stderrors.New("no update fields")
	ErrorBadRequest    = stderrors.New("bad request")
	ErrorExisting      = ErrAlreadyExists
)

// Error 将 ErrorDef 与底层错误链绑定在一起，方便同时保留业务码和原始错误。
type Error struct {
	Def   ErrorDef
	Err   error
	Extra map[string]any
}

// Error 实现 error 接口。
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Def.Key, e.Err)
	}
	return e.Def.Key
}

// Unwrap 返回底层错误，支持 errors.Is / errors.As。
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// Wrap 使用指定错误定义包装底层错误。
func Wrap(def ErrorDef, err error, extra ...map[string]any) error {
	if err == nil {
		return nil
	}
	appErr := &Error{
		Def: def,
		Err: err,
	}
	if len(extra) > 0 {
		appErr.Extra = extra[0]
	}
	return appErr
}

// DefOf 提取错误对应的 ErrorDef；无法识别时回退为 ErrServerInternal。
func DefOf(err error) ErrorDef {
	if err == nil {
		return ErrorDef{}
	}

	var appErr *Error
	if stderrors.As(err, &appErr) && appErr != nil {
		return appErr.Def
	}

	switch {
	case stderrors.Is(err, ErrNotFound):
		return ErrDataNotFound
	case stderrors.Is(err, ErrorNoUpdateField):
		return ErrNoUpdateField
	case stderrors.Is(err, ErrAlreadyExists):
		return ErrorDef{
			HttpStatus: http.StatusConflict,
			Code:       4090,
			Key:        "err.sys.already_exists",
		}
	case stderrors.Is(err, ErrorBadRequest):
		return ErrBadRequest
	case stderrors.Is(err, ErrUnauthorized):
		return ErrorDef{
			HttpStatus: http.StatusUnauthorized,
			Code:       4010,
			Key:        "err.sys.unauthorized",
		}
	case stderrors.Is(err, ErrForbidden):
		return ErrorDef{
			HttpStatus: http.StatusForbidden,
			Code:       4030,
			Key:        "err.sys.forbidden",
		}
	case stderrors.Is(err, ErrConflict):
		return ErrorDef{
			HttpStatus: http.StatusConflict,
			Code:       4091,
			Key:        "err.sys.conflict",
		}
	default:
		return ErrServerInternal
	}
}

// Code 返回业务错误码；无法识别时回退为服务端内部错误码。
func Code(err error) int {
	if err == nil {
		return 0
	}
	return DefOf(err).Code
}

// HTTPStatusOf 返回错误对应的 HTTP 状态码。
func HTTPStatusOf(err error) int {
	if err == nil {
		return http.StatusOK
	}
	def := DefOf(err)
	if def.HttpStatus == 0 {
		return http.StatusInternalServerError
	}
	return def.HttpStatus
}

// KeyOf 返回错误对应的国际化消息键。
func KeyOf(err error) string {
	if err == nil {
		return ""
	}
	return DefOf(err).Key
}

// IsAny 判断 err 是否匹配任意一个目标错误。
func IsAny(err error, targets ...error) bool {
	for _, target := range targets {
		if stderrors.Is(err, target) {
			return true
		}
	}
	return false
}
