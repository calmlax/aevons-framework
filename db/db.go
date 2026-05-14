package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/calmlax/aevons-framework/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB                  *gorm.DB
	errDBNotInitialized = errors.New("db is not initialized")
	errNilConfig        = errors.New("db init: nil config")
	defaultPingTimeout  = 5 * time.Second
)

// Init opens the database connection, configures the connection pool and verifies connectivity.
func Init(cfg *config.Config) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errNilConfig
	}

	gdb, err := open(cfg.DB)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		return nil, err
	}

	applyPoolConfig(sqlDB, cfg.DB)

	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	DB = gdb
	return DB, nil
}

// Raw returns the initialized gorm DB instance.
func Raw() (*gorm.DB, error) {
	if DB == nil {
		return nil, errDBNotInitialized
	}
	return DB, nil
}

// SQLDB returns the underlying database/sql DB for advanced usage.
func SQLDB() (*sql.DB, error) {
	gdb, err := Raw()
	if err != nil {
		return nil, err
	}
	return gdb.DB()
}

// Close closes the underlying sql.DB and resets the global DB handle.
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		DB = nil
		return err
	}

	DB = nil
	return sqlDB.Close()
}

func open(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := strings.TrimSpace(cfg.DSN)
	driver := strings.TrimSpace(strings.ToLower(cfg.Driver))

	switch driver {
	case "postgres", "pg", "postgresql":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mysql", "":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		log.Printf("unsupported db driver %q, defaulting to mysql", cfg.Driver)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}
}

func applyPoolConfig(sqlDB *sql.DB, cfg config.DBConfig) {
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetimeSeconds > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)
	}
	if cfg.ConnMaxIdleTimeSeconds > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTimeSeconds) * time.Second)
	}
}
