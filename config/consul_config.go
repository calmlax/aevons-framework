package config

type ConsulConfig struct {
	Address string `yaml:"address"`
	Enabled bool   `yaml:"enabled"`
}

func (cfg *ConsulConfig) merge(next ConsulConfig, section sectionValues) {
	if section.has("address") {
		cfg.Address = next.Address
	}
	if section.has("enabled") {
		cfg.Enabled = next.Enabled
	}
}
