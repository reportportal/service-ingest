package config

import (
	"fmt"
	"strings"
	"time"

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

type StorageConfig struct {
	Type string `env:"STORAGE_TYPE" envDefault:"local"`

	CatalogPath string `env:"CATALOG_PATH" envDefault:"/data/catalog"`
	BufferPath  string `env:"BUFFER_PATH" envDefault:"/data/buffer"`

	FlushInterval string `env:"FLUSH_INTERVAL" envDefault:"30s"`
	FlushMaxSize  string `env:"FLUSH_MAX_SIZE" envDefault:"10Mb"`
	FlushMaxItems int    `env:"FLUSH_MAX_ITEMS" envDefault:"1000"`

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

	return cfg, nil
}

func (s ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s StorageConfig) FlushIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(s.FlushInterval)
}

func (s StorageConfig) FlushMaxSizeBytes() (int64, error) {
	var size int64
	var unit string

	_, err := fmt.Sscanf(s.FlushMaxSize, "%d%s", &size, &unit)
	if err != nil {
		return 0, fmt.Errorf("invalid flush max size: %w", err)
	}

	switch strings.ToLower(unit) {
	case "b", "":
		return size, nil
	case "kb":
		return size * 1024, nil
	case "mb":
		return size * 1024 * 1024, nil
	case "gb":
		return size * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("unknown size unit: %s", unit)
	}
}
