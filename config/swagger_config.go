package config

// SwaggerConfig Swagger / Apifox 同步配置
type SwaggerConfig struct {
	// Enabled 控制是否暴露 /apifox/openapi.json 接口
	// 生产环境务必设为 false 或不配置（默认 false）
	Enabled bool `mapstructure:"enabled"`
}
