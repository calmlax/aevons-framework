package config

// RedisConfig 定义 Redis 连接与连接池配置。
// 支持单机、哨兵和集群三种模式。
type RedisConfig struct {
	Mode string `yaml:"mode"`
	// Mode 指定 Redis 模式：standalone、sentinel、cluster。

	Address string `yaml:"address"`
	// Address 单机模式下的 Redis 地址。

	Addresses []string `yaml:"addresses"`
	// Addresses 哨兵或集群模式下的节点地址列表。

	MasterName string `yaml:"master_name"`
	// MasterName 哨兵模式下的主节点名称。

	Password string `yaml:"password"`
	// Password Redis 认证密码。

	DB int `yaml:"db"`
	// DB 单机或哨兵模式下使用的逻辑库编号。

	PoolSize int `yaml:"pool_size"`
	// PoolSize 连接池最大连接数。

	MinIdleConns int `yaml:"min_idle_conns"`
	// MinIdleConns 连接池最小空闲连接数。

	MaxRetries int `yaml:"max_retries"`
	// MaxRetries 命令失败后的最大重试次数。

	DialTimeoutSeconds int `yaml:"dial_timeout_seconds"`
	// DialTimeoutSeconds 建立连接超时时间，单位为秒。

	ReadTimeoutSeconds int `yaml:"read_timeout_seconds"`
	// ReadTimeoutSeconds 读取超时时间，单位为秒。

	WriteTimeoutSeconds int `yaml:"write_timeout_seconds"`
	// WriteTimeoutSeconds 写入超时时间，单位为秒。
}

// merge 将新配置中显式声明的 Redis 字段覆盖到当前配置。
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
