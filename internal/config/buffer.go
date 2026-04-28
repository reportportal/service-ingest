package config

import (
	"log/slog"

	"github.com/dustin/go-humanize"
)

type BufferConfig struct {
	BufferPath      string `env:"BUFFER_PATH" envDefault:""`
	BufferCacheSize string `env:"BUFFER_CACHE_SIZE" envDefault:"256MiB"`
	BufferIndexSize string `env:"BUFFER_INDEX_SIZE"`

	FileBufferPath string `env:"FILE_BUFFER_PATH,expand" envDefault:"/data/staging"`
}

func (b *BufferConfig) GetBufferCacheSize() int64 {
	size, err := humanize.ParseBytes(b.BufferCacheSize)
	if err != nil {
		slog.Warn("invalid buffer cache size, using default 256MiB", "value", b.BufferCacheSize, "error", err)
		return 256 << 20
	}
	return int64(size)
}
