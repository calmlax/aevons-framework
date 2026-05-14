package config

// AuthConfig 认证相关配置。
type AuthConfig struct {
	AccessTokenTTL  int64    `mapstructure:"access_token_ttl"`  // 单位：秒，默认 7200（2小时），oauth_client 未配置时兜底
	RefreshTokenTTL int64    `mapstructure:"refresh_token_ttl"` // 单位：秒，默认 604800（7天），oauth_client 未配置时兜底
	EmailCodeTTL    int64    `mapstructure:"email_ttl"`         // 单位：秒，默认 300（5分钟）
	Excludes        []string `mapstructure:"excludes"`          // 免认证路径规则，支持精确后缀和通配符前缀（/prefix/*）
}
