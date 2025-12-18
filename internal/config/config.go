package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	Log     LogConfig
	Storage StorageConfig
}

type ServerConfig struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	Host     string `env:"HOST" envDefault:"0.0.0.0"`
	Env      string `env:"ENVIRONMENT" envDefault:"development"`
	BasePath string `env:"BASE_PATH" envDefault:"/api"`
}

type LogConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"info"`
	Format string `env:"LOG_FORMAT" envDefault:"json"` // json or text
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	if err := env.Parse(&cfg.Server); err != nil {
		return nil, fmt.Errorf("failed to parse server config: %w", err)
	}

	if err := env.Parse(&cfg.Log); err != nil {
		return nil, fmt.Errorf("failed to parse log config: %w", err)
	}

	if err := env.Parse(&cfg.Storage); err != nil {
		return nil, fmt.Errorf("failed to parse storage config: %w", err)
	}

	return cfg, nil
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
