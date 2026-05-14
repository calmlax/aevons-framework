package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Consul     ConsulConfig     `yaml:"consul"`
	DB         DBConfig         `yaml:"db"`
	Redis      RedisConfig      `yaml:"redis"`
	Downstream DownstreamConfig `yaml:"downstream"`
	XSS        XSSConfig        `yaml:"xss"`
	CORS       CORSConfig       `yaml:"cors"`
	Swagger    SwaggerConfig    `yaml:"swagger"`
	Auth       AuthConfig       `yaml:"auth"`
}

type sectionValues map[string]any

func Load(configDir, env string) (Config, error) {
	cfg := defaultConfig(env)

	basePath := filepath.Join(configDir, "config.yaml")
	if err := mergeConfigFromFile(basePath, &cfg); err != nil {
		return Config{}, err
	}

	if env != "" {
		overlayPath := filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
		if _, err := os.Stat(overlayPath); err == nil {
			if err := mergeConfigFromFile(overlayPath, &cfg); err != nil {
				return Config{}, err
			}
		}
	}

	if cfg.Server.Env == "" {
		cfg.Server.Env = env
	}

	return cfg, nil
}

func defaultConfig(env string) Config {
	return Config{
		Server: ServerConfig{
			Env:  env,
			Host: "0.0.0.0",
			Mode: "debug",
			Port: 8080,
		},
		DB: DBConfig{
			Driver:                 "mysql",
			MaxOpenConns:           20,
			MaxIdleConns:           10,
			ConnMaxLifetimeSeconds: 3600,
			ConnMaxIdleTimeSeconds: 1800,
		},
		Redis: RedisConfig{
			Mode:                "standalone",
			Address:             "127.0.0.1:6379",
			PoolSize:            10,
			DialTimeoutSeconds:  5,
			ReadTimeoutSeconds:  3,
			WriteTimeoutSeconds: 3,
		},
		XSS: XSSConfig{
			Enabled: true,
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"*"},
		},
		Swagger: SwaggerConfig{
			Enabled: false,
		},
	}
}

func mergeConfigFromFile(path string, cfg *Config) error {
	next, sections, err := loadConfigFile(path)
	if err != nil {
		return err
	}
	cfg.merge(next, sections)
	return nil
}

func loadConfigFile(path string) (Config, map[string]sectionValues, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	root := make(map[string]any)
	if err := yaml.Unmarshal(data, &root); err != nil {
		return Config{}, nil, fmt.Errorf("parse config keys %s: %w", path, err)
	}

	return cfg, extractSections(root), nil
}

func (cfg *Config) merge(next Config, sections map[string]sectionValues) {
	cfg.Server.merge(next.Server, sections["server"])
	cfg.Consul.merge(next.Consul, sections["consul"])
	cfg.DB.merge(next.DB, sections["db"])
	cfg.Redis.merge(next.Redis, sections["redis"])
	cfg.Downstream.merge(next.Downstream, sections["downstream"])
	cfg.XSS.merge(next.XSS, sections["xss"])
	cfg.CORS.merge(next.CORS, sections["cors"])
	cfg.Swagger.merge(next.Swagger, sections["swagger"])
}

func extractSections(root map[string]any) map[string]sectionValues {
	sections := make(map[string]sectionValues, len(root))
	for key := range root {
		if section, ok := lookupSection(root, key); ok {
			sections[key] = section
		}
	}
	return sections
}

func lookupSection(root map[string]any, key string) (sectionValues, bool) {
	sectionValue, ok := root[key]
	if !ok {
		return nil, false
	}
	section, ok := sectionValue.(map[string]any)
	return section, ok
}

func (section sectionValues) has(key string) bool {
	if section == nil {
		return false
	}
	_, ok := section[key]
	return ok
}

func (c Config) DBAddress() string {
	return c.DB.address()
}
