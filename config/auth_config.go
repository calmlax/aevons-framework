package config

// AuthConfig 认证相关配置。
type AuthConfig struct {
	AccessTokenTTL int64 `yaml:"access_token_ttl"`
	// AccessTokenTTL 访问令牌有效期，单位：秒。

	RefreshTokenTTL int64 `yaml:"refresh_token_ttl"`
	// RefreshTokenTTL 刷新令牌有效期，单位：秒。

	EmailCodeTTL int64 `yaml:"email_ttl"`
	// EmailCodeTTL 邮箱验证码有效期，单位：秒。

	Excludes []string `yaml:"excludes"`
	// Excludes 免认证路径规则，支持精确路径和前缀通配符（/prefix/*）。
}

// merge 将新配置中显式声明的认证字段覆盖到当前配置。
func (cfg *AuthConfig) merge(next AuthConfig, section sectionValues) {
	if section.has("access_token_ttl") {
		cfg.AccessTokenTTL = next.AccessTokenTTL
	}
	if section.has("refresh_token_ttl") {
		cfg.RefreshTokenTTL = next.RefreshTokenTTL
	}
	if section.has("email_ttl") {
		cfg.EmailCodeTTL = next.EmailCodeTTL
	}
	if section.has("excludes") {
		cfg.Excludes = next.Excludes
	}
}
