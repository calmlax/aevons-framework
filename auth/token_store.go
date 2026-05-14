package auth

import (
	"context"
	"time"
)

// TokenStore 定义在 Redis 中存取令牌和验证码的接口。
type TokenStore interface {
	// SaveAccessToken 将 LoginUser 以 Access Token 为键存入 Redis。
	SaveAccessToken(ctx context.Context, token string, user *LoginUser, ttl time.Duration) error

	// SaveRefreshToken 存储 Refresh Token → Access Token 的映射。
	SaveRefreshToken(ctx context.Context, refreshToken, accessToken string, ttl time.Duration) error

	// GetLoginUser 根据 Access Token 读取 LoginUser。
	GetLoginUser(ctx context.Context, accessToken string) (*LoginUser, error)

	// UpdateLoginUser 动态更新缓存在 Redis 中的 LoginUser，保持原 TTL 不变。
	UpdateLoginUser(ctx context.Context, accessToken string, updateFn func(*LoginUser)) error

	// GetAccessTokenByRefresh 根据 Refresh Token 读取对应的 Access Token。
	GetAccessTokenByRefresh(ctx context.Context, refreshToken string) (string, error)

	// GetRefreshTokenTTL 返回 Refresh Token 在 Redis 中的剩余存活时间。
	// 用于刷新时保持绝对过期时间不变，防止无限续期。
	GetRefreshTokenTTL(ctx context.Context, refreshToken string) (time.Duration, error)

	// DeleteAccessToken 删除 Access Token 条目。
	DeleteAccessToken(ctx context.Context, accessToken string) error

	// DeleteRefreshToken 删除 Refresh Token 条目。
	DeleteRefreshToken(ctx context.Context, refreshToken string) error

	// SaveEmailCode 存储针对特定业务用途的一次性邮箱验证码。
	SaveEmailCode(ctx context.Context, email, purpose, code string, ttl time.Duration) error

	// GetEmailCode 读取特定邮箱及业务用途的验证码。
	GetEmailCode(ctx context.Context, email, purpose string) (string, error)

	// DeleteEmailCode 删除特定邮箱及业务用途的验证码条目。
	DeleteEmailCode(ctx context.Context, email, purpose string) error

	// SaveAuthCode 存储授权码，关联 userId、clientId 和授权 scopes。
	SaveAuthCode(ctx context.Context, code string, userId int64, clientId string, scopes []string, ttl time.Duration) error

	// GetAuthCodeInfo 根据授权码读取对应的 userId、clientId 和 scopes。
	GetAuthCodeInfo(ctx context.Context, code string) (userId int64, clientId string, scopes []string, err error)

	// DeleteAuthCode 删除授权码条目。
	DeleteAuthCode(ctx context.Context, code string) error

	// AddUserSession 记录用户的客户端会话（SLO 用）。
	// Key: auth:user_sessions:{userId}  Field: clientId  Value: accessToken
	AddUserSession(ctx context.Context, userId int64, clientId, accessToken string) error

	// GetUserSessions 获取用户所有客户端会话，返回 clientId → accessToken 映射。
	GetUserSessions(ctx context.Context, userId int64) (map[string]string, error)

	// RemoveUserSession 删除用户指定客户端的会话记录。
	RemoveUserSession(ctx context.Context, userId int64, clientId string) error

	// SaveOAuthState 存储授权码模式的 state 参数，防止 CSRF 攻击。
	// Key: auth:oauth2_state:{state}  Value: clientId
	SaveOAuthState(ctx context.Context, state, clientId string, ttl time.Duration) error

	// GetOAuthState 根据 state 读取关联的 clientId。
	GetOAuthState(ctx context.Context, state string) (clientId string, err error)

	// DeleteOAuthState 删除 state 条目。
	DeleteOAuthState(ctx context.Context, state string) error

	// SaveRSAPrivateKey 存储用于登录密码解密的临时私钥。
	SaveRSAPrivateKey(ctx context.Context, keyId, privateKey string, ttl time.Duration) error

	// GetRSAPrivateKey 获取私钥。
	GetRSAPrivateKey(ctx context.Context, keyId string) (string, error)

	// DeleteRSAPrivateKey 删除用过的私钥，保证一次性。
	DeleteRSAPrivateKey(ctx context.Context, keyId string) error
}
