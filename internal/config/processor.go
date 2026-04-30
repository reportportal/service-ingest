package config

import "time"

type ProcessorConfig struct {
	FlushInterval      string `env:"FLUSH_INTERVAL" envDefault:"10s"`
	FilesFlushInterval string `env:"FILES_FLUSH_INTERVAL" envDefault:"10s"`
	ReadLimit          int    `env:"READ_LIMIT" envDefault:"100000"`
}

func (s ProcessorConfig) FlushIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(s.FlushInterval)
}

func (s ProcessorConfig) FilesFlushIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(s.FilesFlushInterval)
}
