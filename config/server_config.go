package config

type ServerConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
	Mode string `yaml:"mode"`
}

func (cfg *ServerConfig) merge(next ServerConfig, section sectionValues) {
	if section.has("name") {
		cfg.Name = next.Name
	}
	if section.has("host") {
		cfg.Host = next.Host
	}
	if section.has("port") {
		cfg.Port = next.Port
	}
	if section.has("env") {
		cfg.Env = next.Env
	}
	if section.has("mode") {
		cfg.Mode = next.Mode
	}
}
