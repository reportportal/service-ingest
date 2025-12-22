package config

import (
	"time"
)

type StorageConfig struct {
	Type string `env:"STORAGE_TYPE" envDefault:"local"`

	CatalogPath string `env:"CATALOG_PATH" envDefault:"/data/catalog"`

	BufferPath  string `env:"BUFFER_PATH" envDefault:"/data/buffer"`
	BufferLease string `env:"BUFFER_LEASE_DURATION" envDefault:"5m"`

	FlushInterval string `env:"FLUSH_INTERVAL" envDefault:"30s"`

	ParquetCompression  string `env:"PARQUET_COMPRESSION" envDefault:"snappy"`
	ParquetRowGroupSize int    `env:"PARQUET_ROW_GROUP_SIZE" envDefault:"1000"`
}

func (s StorageConfig) BufferLeaseDuration() (time.Duration, error) {
	return time.ParseDuration(s.BufferLease)
}

func (s StorageConfig) FlushIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(s.FlushInterval)
}
