package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/calmlax/aevons-framework/consts"

	"github.com/redis/go-redis/v9"
)

// RedisTokenStore 基于 Redis 实现 TokenStore 接口。
type RedisTokenStore struct {
	client *redis.Client
}

// NewRedisTokenStore 创建 RedisTokenStore 实例。
func NewRedisTokenStore(client *redis.Client) *RedisTokenStore {
	return &RedisTokenStore{client: client}
}

func (s *RedisTokenStore) SaveAccessToken(ctx context.Context, token string, user *LoginUser, ttl time.Duration) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("token_store: 序列化 LoginUser 失败: %w", err)
	}
	key := consts.RedisKeyAccessToken + token
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *RedisTokenStore) SaveRefreshToken(ctx context.Context, refreshToken, accessToken string, ttl time.Duration) error {
	key := consts.RedisKeyRefreshToken + refreshToken
	return s.client.Set(ctx, key, accessToken, ttl).Err()
}

func (s *RedisTokenStore) GetLoginUser(ctx context.Context, accessToken string) (*LoginUser, error) {
	key := consts.RedisKeyAccessToken + accessToken
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("token_store: 读取 Access Token 失败: %w", err)
	}
	var user LoginUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("token_store: 反序列化 LoginUser 失败: %w", err)
	}
	return &user, nil
}

func (s *RedisTokenStore) UpdateLoginUser(ctx context.Context, token string, updateFn func(*LoginUser)) error {
	user, err := s.GetLoginUser(ctx, token)
	if err != nil {
		return err
	}
	updateFn(user)
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("token_store: 序列化 LoginUser 失败: %w", err)
	}
	key := consts.RedisKeyAccessToken + token
	return s.client.Set(ctx, key, data, redis.KeepTTL).Err()
}

func (s *RedisTokenStore) GetAccessTokenByRefresh(ctx context.Context, refreshToken string) (string, error) {
	key := consts.RedisKeyRefreshToken + refreshToken
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("token_store: 读取 Refresh Token 失败: %w", err)
	}
	return val, nil
}

// GetRefreshTokenTTL 返回 Refresh Token 在 Redis 中的剩余存活时间。
// 若 key 不存在或已无 TTL，返回 -1。
func (s *RedisTokenStore) GetRefreshTokenTTL(ctx context.Context, refreshToken string) (time.Duration, error) {
	key := consts.RedisKeyRefreshToken + refreshToken
	ttl, err := s.client.TTL(ctx, key).Result()
	if err != nil {
		return -1, fmt.Errorf("token_store: 读取 Refresh Token TTL 失败: %w", err)
	}
	return ttl, nil
}

func (s *RedisTokenStore) DeleteAccessToken(ctx context.Context, accessToken string) error {
	key := consts.RedisKeyAccessToken + accessToken
	return s.client.Del(ctx, key).Err()
}

func (s *RedisTokenStore) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := consts.RedisKeyRefreshToken + refreshToken
	return s.client.Del(ctx, key).Err()
}

func (s *RedisTokenStore) SaveEmailCode(ctx context.Context, email, purpose, code string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s:%s", consts.RedisKeyEmailCode, purpose, email)
	return s.client.Set(ctx, key, code, ttl).Err()
}

func (s *RedisTokenStore) GetEmailCode(ctx context.Context, email, purpose string) (string, error) {
	key := fmt.Sprintf("%s%s:%s", consts.RedisKeyEmailCode, purpose, email)
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("token_store: 读取邮箱验证码失败: %w", err)
	}
	return val, nil
}

func (s *RedisTokenStore) DeleteEmailCode(ctx context.Context, email, purpose string) error {
	key := fmt.Sprintf("%s%s:%s", consts.RedisKeyEmailCode, purpose, email)
	return s.client.Del(ctx, key).Err()
}

func (s *RedisTokenStore) SaveAuthCode(ctx context.Context, code string, userId int64, clientId string, scopes []string, ttl time.Duration) error {
	info := AuthCodeInfo{UserId: userId, ClientId: clientId, Scopes: scopes}
	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("token_store: 序列化 AuthCodeInfo 失败: %w", err)
	}
	key := consts.RedisKeyAuthCode + code
	return s.client.Set(ctx, key, data, ttl).Err()
}

func (s *RedisTokenStore) GetAuthCodeInfo(ctx context.Context, code string) (int64, string, []string, error) {
	key := consts.RedisKeyAuthCode + code
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return 0, "", nil, fmt.Errorf("token_store: 读取授权码失败: %w", err)
	}
	var info AuthCodeInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return 0, "", nil, fmt.Errorf("token_store: 反序列化 AuthCodeInfo 失败: %w", err)
	}
	return info.UserId, info.ClientId, info.Scopes, nil
}

func (s *RedisTokenStore) DeleteAuthCode(ctx context.Context, code string) error {
	key := consts.RedisKeyAuthCode + code
	return s.client.Del(ctx, key).Err()
}

// AddUserSession 记录用户的客户端会话，用于 SLO 全局登出。
func (s *RedisTokenStore) AddUserSession(ctx context.Context, userId int64, clientId, accessToken string) error {
	key := fmt.Sprintf("%s%d", consts.RedisKeyUserSessions, userId)
	return s.client.HSet(ctx, key, clientId, accessToken).Err()
}

// GetUserSessions 获取用户所有客户端会话，返回 clientId → accessToken 映射。
func (s *RedisTokenStore) GetUserSessions(ctx context.Context, userId int64) (map[string]string, error) {
	key := fmt.Sprintf("%s%d", consts.RedisKeyUserSessions, userId)
	result, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("token_store: 读取用户会话失败: %w", err)
	}
	return result, nil
}

// RemoveUserSession 删除用户指定客户端的会话记录。
func (s *RedisTokenStore) RemoveUserSession(ctx context.Context, userId int64, clientId string) error {
	key := fmt.Sprintf("%s%d", consts.RedisKeyUserSessions, userId)
	return s.client.HDel(ctx, key, clientId).Err()
}

// SaveOAuthState 存储授权码模式的 state 参数，防止 CSRF 攻击。
func (s *RedisTokenStore) SaveOAuthState(ctx context.Context, state, clientId string, ttl time.Duration) error {
	key := consts.RedisKeyOAuthState + state
	return s.client.Set(ctx, key, clientId, ttl).Err()
}

// GetOAuthState 根据 state 读取关联的 clientId。
func (s *RedisTokenStore) GetOAuthState(ctx context.Context, state string) (string, error) {
	key := consts.RedisKeyOAuthState + state
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("token_store: 读取 OAuth state 失败: %w", err)
	}
	return val, nil
}

// DeleteOAuthState 删除 state 条目。
func (s *RedisTokenStore) DeleteOAuthState(ctx context.Context, state string) error {
	key := consts.RedisKeyOAuthState + state
	return s.client.Del(ctx, key).Err()
}

// SaveRSAPrivateKey 存储用于登录密码解密的临时私钥。
func (s *RedisTokenStore) SaveRSAPrivateKey(ctx context.Context, keyId, privateKey string, ttl time.Duration) error {
	key := consts.RedisKeyRSAPrivateKey + keyId
	return s.client.Set(ctx, key, privateKey, ttl).Err()
}

// GetRSAPrivateKey 获取私钥。
func (s *RedisTokenStore) GetRSAPrivateKey(ctx context.Context, keyId string) (string, error) {
	key := consts.RedisKeyRSAPrivateKey + keyId
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("token_store: 读取 RSA 私钥失败: %w", err)
	}
	return val, nil
}

// DeleteRSAPrivateKey 删除用过的私钥，保证一次性。
func (s *RedisTokenStore) DeleteRSAPrivateKey(ctx context.Context, keyId string) error {
	key := consts.RedisKeyRSAPrivateKey + keyId
	return s.client.Del(ctx, key).Err()
}
