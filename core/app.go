package core

import (
	"errors"

	"github.com/calmlax/aevons-framework/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var errNilApp = errors.New("app is not initialized")

// App 作为应用装配层上下文，统一承载运行期共享基础依赖。
// 该对象应尽量只在 main、bootstrap、module wiring 等装配层使用。
type App struct {
	config     *config.Config
	redis      *redis.Client
	dataSource *gorm.DB
	// mq     *rocketmq.Client
}

// NewApp 创建应用运行时上下文，统一承载配置和基础组件。
func NewApp(cfg *config.Config, redisClient *redis.Client, dataSource *gorm.DB) *App {
	return &App{
		config:     cfg,
		redis:      redisClient,
		dataSource: dataSource,
	}
}

// RawRedis 返回底层 go-redis 客户端。
func (app *App) RawRedis() (*redis.Client, error) {
	if app == nil {
		return nil, errNilApp
	}
	if app.redis == nil {
		return nil, errors.New("redis client is not initialized")
	}
	return app.redis, nil
}

// RawConfig 返回底层配置对象。
func (app *App) RawConfig() (*config.Config, error) {
	if app == nil {
		return nil, errNilApp
	}
	if app.config == nil {
		return nil, errors.New("config is not initialized")
	}
	return app.config, nil
}

// RawDatabase 返回底层数据库连接。
func (app *App) RawDatabase() (*gorm.DB, error) {
	if app == nil {
		return nil, errNilApp
	}
	if app.dataSource == nil {
		return nil, errors.New("database is not initialized")
	}
	return app.dataSource, nil
}
