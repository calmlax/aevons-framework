package config

type RedisConfig struct {
	Address  string
	Password string
	PoolSize int
	DB       int
}
