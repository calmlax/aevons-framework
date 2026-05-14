package config

// CORSConfig holds CORS settings.
type CORSConfig struct {
	Enabled        bool     `yaml:"enabled"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

func (cfg *CORSConfig) merge(next CORSConfig, section sectionValues) {
	if section.has("enabled") {
		cfg.Enabled = next.Enabled
	}
	if section.has("allowed_origins") {
		cfg.AllowedOrigins = next.AllowedOrigins
	}
}
