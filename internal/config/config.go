package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Log       LogConfig
	Storage   StorageConfig
	Processor ProcessorConfig
	Buffer    BufferConfig
}

type ServerConfig struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	Host     string `env:"HOST" envDefault:"0.0.0.0"`
	Address  string `env:"ADDRESS,expand" envDefault:"$HOST:$PORT"`
	BasePath string `env:"BASE_PATH" envDefault:"/api"`
}

type LogConfig struct {
	Level     string `env:"LOG_LEVEL" envDefault:"info"`
	HTTPLevel string `env:"LOG_HTTP_LEVEL" envDefault:"warn"`
	Format    string `env:"LOG_FORMAT" envDefault:"json"` // json or text
	AddRSBody bool   `env:"LOG_ADD_RESPONSE_BODY" envDefault:"false"`
	AddSource bool   `env:"LOG_ADD_SOURCE" envDefault:"false"`
}

type StorageConfig struct {
	CatalogPath         string `env:"CATALOG_PATH" envDefault:"/data/catalog"`
	ParquetCompression  string `env:"PARQUET_COMPRESSION" envDefault:"snappy"`
	ParquetRowGroupSize int    `env:"PARQUET_ROW_GROUP_SIZE" envDefault:"1000"`
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

	if err := env.Parse(&cfg.Processor); err != nil {
		return nil, fmt.Errorf("failed to parse batch processor config: %w", err)
	}

	if err := env.Parse(&cfg.Buffer); err != nil {
		return nil, fmt.Errorf("failed to parse buffer config: %w", err)
	}

	return cfg, nil
}
