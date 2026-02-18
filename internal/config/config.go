package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type (
	Config struct {
		InstanceConfig
		DBConfig
	}

	InstanceConfig struct {
		ServiceName string `env:"SERVICE_NAME,required"`
		HTTPPort    int    `env:"HTTP_PORT,required"`
	}

	DBConfig struct {
		PGURL          string `env:"PG_URL,required"`
		MigrationsPath string `env:"MIGRATIONS_PATH, default=migrations/"`
	}
)

func Get() (Config, error) {
	return parseConfig()
}

func parseConfig() (cfg Config, err error) {
	godotenv.Load()

	err = envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("fill config: %w", err)
	}

	return cfg, nil
}
