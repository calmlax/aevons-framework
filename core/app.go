package core

import (
	"errors"

	"github.com/calmlax/aevons-framework/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var errNilApp = errors.New("app is not initialized")

// App 封装服务运行期共享的基础依赖。
type App struct {
	Config     *config.Config
	Redis      *redis.Client
	DataSource *gorm.DB
	// mq     *rocketmq.Client
}

// NewApp 创建应用运行时上下文，统一承载配置和基础组件。
func NewApp(cfg *config.Config, redisClient *redis.Client, dataSource *gorm.DB) *App {
	return &App{
		Config:     cfg,
		Redis:      redisClient,
		DataSource: dataSource,
	}
}

// RawRedis 返回底层 go-redis 客户端。
func (app *App) RawRedis() (*redis.Client, error) {
	if app == nil {
		return nil, errNilApp
	}
	if app.Redis == nil {
		return nil, errors.New("redis client is not initialized")
	}
	return app.Redis, nil
}

// RawConfig 返回底层配置对象。
func (app *App) RawConfig() (*config.Config, error) {
	if app == nil {
		return nil, errNilApp
	}
	if app.Config == nil {
		return nil, errors.New("config is not initialized")
	}
	return app.Config, nil
}

// RawDatabase 返回底层数据库连接。
func (app *App) RawDatabase() (*gorm.DB, error) {
	if app == nil {
		return nil, errNilApp
	}
	if app.DataSource == nil {
		return nil, errors.New("database is not initialized")
	}
	return app.DataSource, nil
}
