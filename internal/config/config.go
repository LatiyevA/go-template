package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	InstanceConfig
	DBConfig
}

type InstanceConfig struct {
	ServiceName     string        `env:"SERVICE_NAME,required"`
	HTTPPort        int           `env:"HTTP_PORT,default=8080"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,default=10s"`
}

type DBConfig struct {
	PGURL           string        `env:"PG_URL,required"`
	MigrationsPath  string        `env:"MIGRATIONS_PATH,default=migrations"`
	MaxOpenConns    int           `env:"PG_MAX_OPEN_CONNS,default=10"`
	MaxIdleConns    int           `env:"PG_MAX_IDLE_CONNS,default=2"`
	ConnMaxLifetime time.Duration `env:"PG_CONN_MAX_LIFETIME,default=30m"`
	ConnMaxIdleTime time.Duration `env:"PG_CONN_MAX_IDLE_TIME,default=5m"`
	ConnectTimeout  time.Duration `env:"PG_CONNECT_TIMEOUT,default=5s"`
}

// Get loads optional local dotenv values without overriding the process environment.
func Get() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("load dotenv: %w", err)
	}

	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return Config{}, fmt.Errorf("fill config: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) validate() error {
	switch {
	case strings.TrimSpace(c.ServiceName) == "":
		return errors.New("SERVICE_NAME must not be empty")
	case strings.TrimSpace(c.PGURL) == "":
		return errors.New("PG_URL must not be empty")
	case c.HTTPPort < 1 || c.HTTPPort > 65535:
		return errors.New("HTTP_PORT must be between 1 and 65535")
	case c.ShutdownTimeout <= 0:
		return errors.New("HTTP_SHUTDOWN_TIMEOUT must be positive")
	case c.MaxOpenConns <= 0:
		return errors.New("PG_MAX_OPEN_CONNS must be positive")
	case c.MaxIdleConns < 0 || c.MaxIdleConns > c.MaxOpenConns:
		return errors.New("PG_MAX_IDLE_CONNS must be between 0 and PG_MAX_OPEN_CONNS")
	case c.ConnectTimeout <= 0 || c.ConnMaxLifetime <= 0 || c.ConnMaxIdleTime <= 0:
		return errors.New("database connection timeouts must be positive")
	}
	return nil
}
