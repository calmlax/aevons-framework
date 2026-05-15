package auth

// AuthError 表示带有 HTTP 状态码的认证领域错误。
type AuthError struct {
	Code       string
	HTTPStatus int
}

func (e *AuthError) Error() string {
	return e.Code
}
