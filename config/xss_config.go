package config

// XSSConfig XSS 过滤相关配置
type XSSConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Excludes []string `yaml:"excludes"` // 免过滤路径规则，如：/api/v1/content/*
}

func (cfg *XSSConfig) merge(next XSSConfig, section sectionValues) {
	if section.has("enabled") {
		cfg.Enabled = next.Enabled
	}
	if section.has("excludes") {
		cfg.Excludes = next.Excludes
	}
}
