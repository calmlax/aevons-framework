package err

import "net/http"

// ErrorDef 全局业务错误定义结构体
// 统一封装 HTTP 状态码、业务错误码、国际化错误键，实现错误标准化处理
type ErrorDef struct {
	// HttpStatus HTTP 响应状态码（如 400、404、500），用于 HTTP 响应头
	HttpStatus int

	// Code 业务自定义错误码（如 1001、1002），用于前端业务逻辑判断
	Code int

	// Key 多语言国际化键，前端通过该 key 匹配对应语言的错误提示文本
	Key string
}

// ========================
// 全局错误常量定义（带完整注释）
// ========================
var (
	// ========================
	// 通用请求校验类错误（1001 - 1004）
	// ========================

	// ErrInvalidQuery
	// @HttpStatus 400
	// @业务含义 URL 查询参数格式错误、缺失或不合法
	ErrInvalidQuery = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       1001,
		Key:        "err.sys.invalid_query",
	}

	// ErrInvalidId
	// @HttpStatus 400
	// @业务含义 传入的 ID 格式非法、非数字或为空
	ErrInvalidId = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       1002,
		Key:        "err.sys.invalid_id",
	}

	// ErrInvalidBody
	// @HttpStatus 400
	// @业务含义 请求体 JSON 解析失败、字段类型不匹配
	ErrInvalidBody = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       1003,
		Key:        "err.sys.invalid_body",
	}

	// ErrNoUpdateField
	// @HttpStatus 400
	// @业务含义 更新操作未传入任何可修改字段
	ErrNoUpdateField = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       1004,
		Key:        "err.sys.no_update_field",
	}

	// ========================
	// 通用业务逻辑类错误（1005 - 1010）
	// ========================

	// ErrDataNotFound
	// @HttpStatus 404
	// @业务含义 操作/查询的数据不存在
	ErrDataNotFound = ErrorDef{
		HttpStatus: http.StatusNotFound,
		Code:       1005,
		Key:        "err.sys.data_not_found",
	}

	// ErrQueryFailed
	// @HttpStatus 500
	// @业务含义 数据查询失败（数据库异常、SQL 错误等）
	ErrQueryFailed = ErrorDef{
		HttpStatus: http.StatusInternalServerError,
		Code:       1006,
		Key:        "err.sys.query_failed",
	}

	// ErrCreateFailed
	// @HttpStatus 500
	// @业务含义 数据创建/保存失败
	ErrCreateFailed = ErrorDef{
		HttpStatus: http.StatusInternalServerError,
		Code:       1007,
		Key:        "err.sys.create_failed",
	}

	// ErrUpdateFailed
	// @HttpStatus 500
	// @业务含义 数据更新/修改失败
	ErrUpdateFailed = ErrorDef{
		HttpStatus: http.StatusInternalServerError,
		Code:       1008,
		Key:        "err.sys.update_failed",
	}

	// ErrDeleteFailed
	// @HttpStatus 500
	// @业务含义 数据删除失败（可能存在关联数据、状态限制）
	ErrDeleteFailed = ErrorDef{
		HttpStatus: http.StatusInternalServerError,
		Code:       1009,
		Key:        "err.sys.delete_failed",
	}

	// ErrCacheRefresh
	// @HttpStatus 500
	// @业务含义 缓存刷新/同步失败（Redis 等中间件异常）
	ErrCacheRefresh = ErrorDef{
		HttpStatus: http.StatusInternalServerError,
		Code:       1010,
		Key:        "err.sys.refresh_failed",
	}
	// ==================== 参数校验类错误 4000 ====================
	// ErrValidationRequired
	// @HttpStatus 400
	// @业务含义 字段未填写，为必填项
	ErrValidationRequired = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       4000,
		Key:        "err.sys.validation.required",
	}

	// ErrValidationMax
	// @HttpStatus 400
	// @业务含义 字段长度/值超出最大限制
	ErrValidationMax = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       4000,
		Key:        "err.sys.validation.max",
	}

	// ErrValidationMin
	// @HttpStatus 400
	// @业务含义 字段长度/值小于最小限制
	ErrValidationMin = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       4000,
		Key:        "err.sys.validation.min",
	}

	// ErrValidationPhone
	// @HttpStatus 400
	// @业务含义 手机号格式不合法
	ErrValidationPhone = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       4000,
		Key:        "err.sys.validation.phone",
	}

	// ErrValidationEmail
	// @HttpStatus 400
	// @业务含义 邮箱格式不合法
	ErrValidationEmail = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       4000,
		Key:        "err.sys.validation.email",
	}

	// ErrValidationOneof
	// @HttpStatus 400
	// @业务含义 字段值不在允许的枚举范围内
	ErrValidationOneof = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       4000,
		Key:        "err.sys.validation.oneof",
	}

	// ErrBadRequest
	// @HttpStatus 400
	// @业务含义 通用非法请求
	ErrBadRequest = ErrorDef{
		HttpStatus: http.StatusBadRequest,
		Code:       4000,
		Key:        "err.sys.bad_request",
	}

	// ErrServerInternal
	// @HttpStatus 500
	// @业务含义 服务器内部错误（未捕获异常、panic 等）
	ErrServerInternal = ErrorDef{
		HttpStatus: http.StatusInternalServerError,
		Code:       5000,
		Key:        "err.sys.server_error",
	}
)
