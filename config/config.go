package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Server     ServerConfig
	Consul     ConsulConfig
	DB         DBConfig
	Redis      RedisConfig
	Downstream DownstreamConfig
	XSS        XSSConfig
	CORS       CORSConfig
	Swagger    SwaggerConfig
}

type DownstreamConfig struct {
	TaskServiceURL     string
	ResourceServiceURL string
}

func Load(configDir, env string) (Config, error) {
	cfg := Config{
		Server: ServerConfig{
			Env:  env,
			Host: "0.0.0.0",
		},
	}

	basePath := filepath.Join(configDir, "config.yaml")
	if err := loadFile(basePath, &cfg); err != nil {
		return Config{}, err
	}

	overlayPath := filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
	if _, err := os.Stat(overlayPath); err == nil {
		if err := loadFile(overlayPath, &cfg); err != nil {
			return Config{}, err
		}
	}

	if cfg.Server.Env == "" {
		cfg.Server.Env = env
	}

	return cfg, nil
}

func loadFile(path string, cfg *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	section := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasSuffix(line, ":") && !strings.Contains(line, " ") {
			section = strings.TrimSuffix(line, ":")
			continue
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return fmt.Errorf("invalid config line %q in %s", line, path)
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"'")

		switch section {
		case "server":
			if err := applyServerField(&cfg.Server, key, value); err != nil {
				return fmt.Errorf("parse server config %s: %w", path, err)
			}
		case "consul":
			if err := applyConsulField(&cfg.Consul, key, value); err != nil {
				return fmt.Errorf("parse consul config %s: %w", path, err)
			}
		case "db":
			applyDBField(&cfg.DB, key, value)
		case "redis":
			applyRedisField(&cfg.Redis, key, value)
		case "downstream":
			applyDownstreamField(&cfg.Downstream, key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan config %s: %w", path, err)
	}

	return nil
}

func applyServerField(server *ServerConfig, key, value string) error {
	switch key {
	case "name":
		server.Name = value
	case "host":
		server.Host = value
	case "port":
		port, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid port %q: %w", value, err)
		}
		server.Port = port
	case "env":
		server.Env = value
	}

	return nil
}

func applyConsulField(consul *ConsulConfig, key, value string) error {
	switch key {
	case "address":
		consul.Address = value
	case "enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid enabled value %q: %w", value, err)
		}
		consul.Enabled = enabled
	}

	return nil
}

func applyDBField(db *DBConfig, key, value string) {
	switch key {
	case "driver":
		db.Driver = value
	case "dsn":
		db.DSN = value
	}
}

func applyRedisField(redis *RedisConfig, key, value string) {
	switch key {
	case "address":
		redis.Address = value
	case "host":
		redis.Address = value
	}
}

func applyDownstreamField(downstream *DownstreamConfig, key, value string) {
	switch key {
	case "task_service_url":
		downstream.TaskServiceURL = value
	case "resource_service_url":
		downstream.ResourceServiceURL = value
	}
}

func (c Config) DBAddress() string {
	driver := strings.TrimSpace(strings.ToLower(c.DB.Driver))
	dsn := strings.TrimSpace(c.DB.DSN)
	if driver == "" || dsn == "" {
		return ""
	}

	switch driver {
	case "mysql":
		return parseMySQLAddress(dsn)
	case "postgres", "postgresql":
		return parsePostgresAddress(dsn)
	default:
		return ""
	}
}

func parseMySQLAddress(dsn string) string {
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

func parsePostgresAddress(dsn string) string {
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
