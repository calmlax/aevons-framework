package config

// XSSConfig XSS 过滤相关配置
type XSSConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	Excludes []string `mapstructure:"excludes"` // 免过滤路径规则，如：/api/v1/content/*
}
