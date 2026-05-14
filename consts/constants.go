package consts

// 权限相关常量
const (
	SuperAdminId = 1 // 超级管理员

	AllPermission     = "*"                // 超级权限（拥有所有权限）
	UserPermissionKey = "user_permissions" // Context中存储用户权限的key

	SuperAdminRoleKey = "admin"      // 超级管理员角色标识
	UserRoleKey       = "user_roles" // Context中存储用户角色的key
	UserDeptKey       = "user_depts" // Context中存储用户部门的key
)

// 错误码/错误信息常量（建议与业务错误码统一管理）
const (
	ErrUnauthorized     = "authorization.unauthorized"      // 未授权
	ErrPermissionDenied = "authorization.permission.denied" // 权限不足
	ErrRoleDenied       = "authorization.role.denied"       // 角色拒绝
	ErrRoleInvalid      = "authorization.role.invalid"      // 角色格式无效
)

// 认证相关 Context key
const (
	UserIdKey    = "user_id"    // gin.Context 中存储当前用户 Id 的 key
	LoginUserKey = "login_user" // gin.Context 中存储 LoginUser 对象的 key
)

// Redis key 前缀常量
const (
	RedisKeyAccessToken   = "aevons:auth:access_token:"  // access token → LoginUser JSON
	RedisKeyRefreshToken  = "aevons:auth:refresh_token:" // refresh token → access token
	RedisKeyAuthCode      = "aevons:auth:auth_code:"     // auth code → AuthCodeInfo JSON
	RedisKeyEmailCode     = "aevons:captcha:email:"      // email → 6-digit code
	RedisKeyUserSessions  = "aevons:auth:user_sessions:" // userId → Hash(clientId → accessToken)
	RedisKeyOAuthState    = "aevons:auth:oauth2_state:"  // OAuth2 state → clientId (CSRF protection)
	RedisKeyRSAPrivateKey = "aevons:auth:rsa_priv:"      // keyId → RSAPrivateKey

	ConfCacheKeyPrefix = "aevons:sys:conf:"

	DictCacheKeyPrefix = "aevons:sys:dict_data:"
)

// 认证错误码常量
const (
	ErrInvalidCredentials  = "auth.invalid_credentials"   // 用户名或密码错误
	ErrAccountDisabled     = "auth.account_disabled"      // 账号已禁用
	ErrTokenMissing        = "auth.token_missing"         // 令牌缺失或格式错误
	ErrTokenExpired        = "auth.token_expired"         // 令牌已过期或不存在
	ErrInvalidCode         = "auth.invalid_code"          // 验证码无效或已过期
	ErrInvalidAuthCode     = "auth.invalid_auth_code"     // 授权码无效或已过期
	ErrUserNotFound        = "auth.user_not_found"        // 用户不存在
	ErrInvalidRefreshToken = "auth.invalid_refresh_token" // Refresh Token 无效
)

// OAuth2 错误码常量
const (
	ErrOAuthInvalidClient       = "oauth2.invalid_client"         // client_id 不存在或 client_secret 不匹配
	ErrOAuthUnsupportedGrant    = "oauth2.unsupported_grant_type" // grant_type 不被客户端支持
	ErrOAuthRedirectURIMismatch = "oauth2.redirect_uri_mismatch"  // redirect_uri 与客户端配置不匹配
	ErrOAuthInvalidState        = "oauth2.invalid_state"          // state 参数校验失败（CSRF 防护）
)

// BizType 业务操作日志类型
type BizType string

const (
	OTHER   BizType = "OTHER"
	INSERT  BizType = "INSERT"
	UPDATE  BizType = "UPDATE"
	DELETE  BizType = "DELETE"
	AUTH    BizType = "AUTH"
	EXPORT  BizType = "EXPORT"
	IMPORT  BizType = "IMPORT"
	KICKED  BizType = "KICKED"
	CLEAN   BizType = "CLEAN"
	SETUP   BizType = "SETUP"
	SAVE    BizType = "SAVE"
	RELEASE BizType = "RELEASE"
	COPY    BizType = "COPY"
	SYNCH   BizType = "SYNCH"
)

const (
	AcceptLanguage = "Accept-Language"
)
