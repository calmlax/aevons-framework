package config

type RedisConfig struct {
	Mode                string   `yaml:"mode"`
	Address             string   `yaml:"address"`
	Addresses           []string `yaml:"addresses"`
	MasterName          string   `yaml:"master_name"`
	Password            string   `yaml:"password"`
	DB                  int      `yaml:"db"`
	PoolSize            int      `yaml:"pool_size"`
	MinIdleConns        int      `yaml:"min_idle_conns"`
	MaxRetries          int      `yaml:"max_retries"`
	DialTimeoutSeconds  int      `yaml:"dial_timeout_seconds"`
	ReadTimeoutSeconds  int      `yaml:"read_timeout_seconds"`
	WriteTimeoutSeconds int      `yaml:"write_timeout_seconds"`
}

func (cfg *RedisConfig) merge(next RedisConfig, section sectionValues) {
	if section.has("mode") {
		cfg.Mode = next.Mode
	}
	if section.has("address") {
		cfg.Address = next.Address
	}
	if section.has("addresses") {
		cfg.Addresses = next.Addresses
	}
	if section.has("master_name") {
		cfg.MasterName = next.MasterName
	}
	if section.has("password") {
		cfg.Password = next.Password
	}
	if section.has("db") {
		cfg.DB = next.DB
	}
	if section.has("pool_size") {
		cfg.PoolSize = next.PoolSize
	}
	if section.has("min_idle_conns") {
		cfg.MinIdleConns = next.MinIdleConns
	}
	if section.has("max_retries") {
		cfg.MaxRetries = next.MaxRetries
	}
	if section.has("dial_timeout_seconds") {
		cfg.DialTimeoutSeconds = next.DialTimeoutSeconds
	}
	if section.has("read_timeout_seconds") {
		cfg.ReadTimeoutSeconds = next.ReadTimeoutSeconds
	}
	if section.has("write_timeout_seconds") {
		cfg.WriteTimeoutSeconds = next.WriteTimeoutSeconds
	}
}
