package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/calmlax/aevons-framework/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

var errClientNotInitialized = errors.New("redis client is not initialized")

// Init 根据配置初始化 Redis 客户端，并通过 Ping 验证连接。
func Init(cfg *config.Config) error {
	if cfg == nil {
		return errors.New("redis init: nil config")
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		PoolSize: cfg.Redis.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return err
	}
	return nil
}

// Raw returns the underlying go-redis client.
func Raw() (*redis.Client, error) {
	if Client == nil {
		return nil, errClientNotInitialized
	}
	return Client, nil
}

// Set 设置键值对，ttl 为过期时间，0 表示永不过期。
func Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.Set(ctx, key, value, ttl).Err()
}

// Get 获取键对应的字符串值。
// 键不存在时返回 ("", redis.Nil)。
func Get(ctx context.Context, key string) (string, error) {
	client, err := Raw()
	if err != nil {
		return "", err
	}
	return client.Get(ctx, key).Result()
}

// Del 删除一个或多个键。
func Del(ctx context.Context, keys ...string) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.Del(ctx, keys...).Err()
}

// Exists 判断键是否存在，返回存在的键数量。
func Exists(ctx context.Context, keys ...string) (int64, error) {
	client, err := Raw()
	if err != nil {
		return 0, err
	}
	return client.Exists(ctx, keys...).Result()
}

// Expire 为已存在的键设置过期时间。
func Expire(ctx context.Context, key string, ttl time.Duration) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.Expire(ctx, key, ttl).Err()
}

// TTL 返回键的剩余过期时间。
func TTL(ctx context.Context, key string) (time.Duration, error) {
	client, err := Raw()
	if err != nil {
		return 0, err
	}
	return client.TTL(ctx, key).Result()
}

// SetJSON 将 value 序列化为 JSON 后存入 Redis，支持设置过期时间。
func SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, b, ttl).Err()
}

// GetJSON 从 Redis 获取 JSON 数据并反序列化到 dest。
func GetJSON(ctx context.Context, key string, dest any) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	b, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

// Incr 将键的整数值原子性地加 1。
func Incr(ctx context.Context, key string) (int64, error) {
	client, err := Raw()
	if err != nil {
		return 0, err
	}
	return client.Incr(ctx, key).Result()
}

// IncrBy 将键的整数值原子性地加 n。
func IncrBy(ctx context.Context, key string, n int64) (int64, error) {
	client, err := Raw()
	if err != nil {
		return 0, err
	}
	return client.IncrBy(ctx, key, n).Result()
}

// HSet 设置哈希表中的字段值，支持多个 field-value 对。
func HSet(ctx context.Context, key string, values ...any) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希表中指定字段的值。
func HGet(ctx context.Context, key, field string) (string, error) {
	client, err := Raw()
	if err != nil {
		return "", err
	}
	return client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希表中所有字段和值。
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	client, err := Raw()
	if err != nil {
		return nil, err
	}
	return client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表中的一个或多个字段。
func HDel(ctx context.Context, key string, fields ...string) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.HDel(ctx, key, fields...).Err()
}

// LPush 将一个或多个值插入列表头部。
func LPush(ctx context.Context, key string, values ...any) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.LPush(ctx, key, values...).Err()
}

// RPush 将一个或多个值追加到列表尾部。
func RPush(ctx context.Context, key string, values ...any) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.RPush(ctx, key, values...).Err()
}

// LPop 移除并返回列表的第一个元素。
func LPop(ctx context.Context, key string) (string, error) {
	client, err := Raw()
	if err != nil {
		return "", err
	}
	return client.LPop(ctx, key).Result()
}

// RPop 移除并返回列表的最后一个元素。
func RPop(ctx context.Context, key string) (string, error) {
	client, err := Raw()
	if err != nil {
		return "", err
	}
	return client.RPop(ctx, key).Result()
}

// SAdd 向集合中添加一个或多个成员。
func SAdd(ctx context.Context, key string, members ...any) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.SAdd(ctx, key, members...).Err()
}

// SMembers 返回集合中的所有成员。
func SMembers(ctx context.Context, key string) ([]string, error) {
	client, err := Raw()
	if err != nil {
		return nil, err
	}
	return client.SMembers(ctx, key).Result()
}

// SIsMember 判断 member 是否是集合的成员。
func SIsMember(ctx context.Context, key string, member any) (bool, error) {
	client, err := Raw()
	if err != nil {
		return false, err
	}
	return client.SIsMember(ctx, key, member).Result()
}

// ZAdd 向有序集合中添加带分数的成员。
func ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	client, err := Raw()
	if err != nil {
		return err
	}
	return client.ZAdd(ctx, key, members...).Err()
}

// ZRange 按排名范围返回有序集合中的成员（升序）。
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	client, err := Raw()
	if err != nil {
		return nil, err
	}
	return client.ZRange(ctx, key, start, stop).Result()
}

// ZScore 返回有序集合中指定成员的分数。
func ZScore(ctx context.Context, key string, member string) (float64, error) {
	client, err := Raw()
	if err != nil {
		return 0, err
	}
	return client.ZScore(ctx, key, member).Result()
}

// Close 关闭 Redis 客户端连接。
func Close() error {
	if Client != nil {
		err := Client.Close()
		Client = nil
		return err
	}
	return nil
}
