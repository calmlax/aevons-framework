package config

import (
	"net/url"
	"strings"
)

// DBConfig 定义数据库连接与连接池配置。
type DBConfig struct {
	Driver string `yaml:"driver"`
	// Driver 数据库驱动名称，如 mysql、postgres。

	DSN string `yaml:"dsn"`
	// DSN 数据库连接串。

	MaxOpenConns int `yaml:"max_open_conns"`
	// MaxOpenConns 连接池允许的最大打开连接数。

	MaxIdleConns int `yaml:"max_idle_conns"`
	// MaxIdleConns 连接池允许保留的最大空闲连接数。

	ConnMaxLifetimeSeconds int `yaml:"conn_max_lifetime_seconds"`
	// ConnMaxLifetimeSeconds 单个连接的最大生命周期，单位为秒。

	ConnMaxIdleTimeSeconds int `yaml:"conn_max_idle_time_seconds"`
	// ConnMaxIdleTimeSeconds 空闲连接的最大保留时间，单位为秒。
}

// merge 将新配置中显式声明的数据库字段覆盖到当前配置。
func (cfg *DBConfig) merge(next DBConfig, section sectionValues) {
	if section.has("driver") {
		cfg.Driver = next.Driver
	}
	if section.has("dsn") {
		cfg.DSN = next.DSN
	}
	if section.has("max_open_conns") {
		cfg.MaxOpenConns = next.MaxOpenConns
	}
	if section.has("max_idle_conns") {
		cfg.MaxIdleConns = next.MaxIdleConns
	}
	if section.has("conn_max_lifetime_seconds") {
		cfg.ConnMaxLifetimeSeconds = next.ConnMaxLifetimeSeconds
	}
	if section.has("conn_max_idle_time_seconds") {
		cfg.ConnMaxIdleTimeSeconds = next.ConnMaxIdleTimeSeconds
	}
}

// address 解析当前 DSN 对应的数据库地址，主要用于探活、展示或运维辅助场景。
func (cfg DBConfig) address() string {
	driver := strings.TrimSpace(strings.ToLower(cfg.Driver))
	dsn := strings.TrimSpace(cfg.DSN)
	if driver == "" || dsn == "" {
		return ""
	}

	switch driver {
	case "mysql":
		return mysqlAddressFromDSN(dsn)
	case "postgres", "postgresql":
		return postgresAddressFromDSN(dsn)
	default:
		return ""
	}
}

// mysqlAddressFromDSN 从 MySQL DSN 中提取 host:port。
func mysqlAddressFromDSN(dsn string) string {
	const marker = "@tcp("
	start := strings.Index(dsn, marker)
	if start < 0 {
		return ""
	}
	start += len(marker)
	end := strings.Index(dsn[start:], ")")
	if end < 0 {
		return ""
	}
	return strings.TrimSpace(dsn[start : start+end])
}

// postgresAddressFromDSN 从 PostgreSQL DSN 中提取 host:port。
func postgresAddressFromDSN(dsn string) string {
	if strings.Contains(dsn, "://") {
		u, err := url.Parse(dsn)
		if err == nil {
			host := u.Hostname()
			port := u.Port()
			if host != "" && port != "" {
				return host + ":" + port
			}
			if host != "" {
				return host + ":5432"
			}
		}
	}

	var host string
	port := "5432"
	for _, field := range strings.Fields(dsn) {
		key, value, ok := strings.Cut(field, "=")
		if !ok {
			continue
		}
		switch key {
		case "host":
			host = value
		case "port":
			port = value
		}
	}
	if host == "" {
		return ""
	}
	return host + ":" + port
}
