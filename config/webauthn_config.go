package config

// WebAuthnConfig WebAuthn / Passkey 配置
type WebAuthnConfig struct {
	// RPID 依赖方 ID，通常是域名，如 localhost 或 example.com
	RPID string `yaml:"rp_id"`
	// RPOrigins 依赖方来源列表，支持多个，如 ["http://localhost:5173", "https://example.com"]
	RPOrigins []string `yaml:"rp_origins"`
	// RPName 依赖方显示名称
	RPName string `yaml:"rp_name"`
}

// merge 将新配置中显式声明的 WebAuthn 字段覆盖到当前配置。
func (cfg *WebAuthnConfig) merge(next WebAuthnConfig, section sectionValues) {
	if section.has("rp_id") {
		cfg.RPID = next.RPID
	}
	if section.has("rp_origins") {
		cfg.RPOrigins = next.RPOrigins
	}
	if section.has("rp_name") {
		cfg.RPName = next.RPName
	}
}
