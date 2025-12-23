package config

import (
	"time"
)

type StorageConfig struct {
	CatalogPath   string `env:"CATALOG_PATH" envDefault:"/data/catalog"`
	BufferPath    string `env:"BUFFER_PATH" envDefault:"/data/buffer"`
	FlushInterval string `env:"FLUSH_INTERVAL" envDefault:"30s"`
	BatchWindow   string `env:"BATCH_WINDOW" envDefault:"60s"`
	ReadLimit     int    `env:"READ_LIMIT" envDefault:"1000"`

	ParquetCompression  string `env:"PARQUET_COMPRESSION" envDefault:"snappy"`
	ParquetRowGroupSize int    `env:"PARQUET_ROW_GROUP_SIZE" envDefault:"1000"`
}

func (s StorageConfig) FlushIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(s.FlushInterval)
}

func (s StorageConfig) BatchWindowDuration() (time.Duration, error) {
	return time.ParseDuration(s.BatchWindow)
}
