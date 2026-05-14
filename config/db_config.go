package config

import (
	"net/url"
	"strings"
)

type DBConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

func (cfg *DBConfig) merge(next DBConfig, section sectionValues) {
	if section.has("driver") {
		cfg.Driver = next.Driver
	}
	if section.has("dsn") {
		cfg.DSN = next.DSN
	}
}

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
