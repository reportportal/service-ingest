package config

import (
	"fmt"
	"strings"
	"time"
)

type StorageConfig struct {
	Type string `env:"STORAGE_TYPE" envDefault:"local"`

	CatalogPath string `env:"CATALOG_PATH" envDefault:"/data/catalog"`
	BufferPath  string `env:"BUFFER_PATH" envDefault:"/data/buffer"`
	BufferLease string `env:"BUFFER_LEASE_DURATION" envDefault:"5m"`

	FlushInterval string `env:"FLUSH_INTERVAL" envDefault:"30s"`
	FlushMaxSize  string `env:"FLUSH_MAX_SIZE" envDefault:"10Mb"`
	FlushMaxItems int    `env:"FLUSH_MAX_ITEMS" envDefault:"1000"`

	ParquetCompression  string `env:"PARQUET_COMPRESSION" envDefault:"snappy"`
	ParquetRowGroupSize int    `env:"PARQUET_ROW_GROUP_SIZE" envDefault:"1000"`
}

func (s StorageConfig) FlushIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(s.FlushInterval)
}

func (s StorageConfig) BufferLeaseDuration() (time.Duration, error) {
	return time.ParseDuration(s.BufferLease)
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
