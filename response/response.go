package response

import (
	"bytes"
	"encoding/json"
	"net/http"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/gin-gonic/gin"
)

// Response 是统一定义的 API 响应载体结构。
type Response struct {
	Code    int            `json:"code"`
	Message string         `json:"message,omitempty"`
	Args    map[string]any `json:"args,omitempty"`
	Data    any            `json:"data"`
}

// Success 输出一个成功的 JSON 响应，默认 HTTP 状态码为 200，业务错误码 Code 为 0。
// args 变长参数用于支持替换多语言里的动态占位符。
func Success(c *gin.Context, data any, args ...map[string]any) {
	var params map[string]any
	if len(args) > 0 {
		params = args[0]
	}
	writeJSON(c, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Args:    params,
		Data:    data,
	})
}

// Fail 根据传入的 HTTP 状态码、业务错误码 Code 和多语言标识 Key，输出一个错误 JSON 响应。
func Fail(c *gin.Context, httpStatus int, code int, message string, args ...map[string]any) {
	var params map[string]any
	if len(args) > 0 {
		params = args[0]
	}
	writeJSON(c, httpStatus, Response{
		Code:    code,
		Message: message,
		Args:    params,
		Data:    nil,
	})
}

// FailByErr 使用统一封装的 ErrorDef 常量结构体快速输出标准的错误 JSON 响应。
func FailByErr(c *gin.Context, httpStatus int, errDef apperr.ErrorDef, args ...map[string]any) {
	Fail(c, httpStatus, errDef.Code, errDef.Key, args...)
}

func FailBy(c *gin.Context, errDef apperr.ErrorDef, args ...map[string]any) {
	Fail(c, errDef.HttpStatus, errDef.Code, errDef.Key, args...)
}

// FailServerError 包装系统内部报错，返回 HTTP 500 以及对应的多语言标识 Key。
func FailServerError(c *gin.Context, message string, args ...map[string]any) {
	Fail(c, http.StatusInternalServerError, 500, message, args...)
}

// writeJSON 使用关闭 HTML 转义的 encoder 写出 JSON，避免 & < > 被转义为 \u0026 等。
func writeJSON(c *gin.Context, status int, v any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(status, "application/json; charset=utf-8", buf.Bytes())
}
