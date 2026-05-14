package config

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	PoolSize int    `yaml:"pool_size"`
	DB       int    `yaml:"db"`
}

func (cfg *RedisConfig) merge(next RedisConfig, section sectionValues) {
	if section.has("address") {
		cfg.Address = next.Address
	}
	if section.has("password") {
		cfg.Password = next.Password
	}
	if section.has("pool_size") {
		cfg.PoolSize = next.PoolSize
	}
	if section.has("db") {
		cfg.DB = next.DB
	}
}
