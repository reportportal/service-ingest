package config

import "time"

type BatchProcessorConfig struct {
	FlushInterval string `env:"FLUSH_INTERVAL" envDefault:"30s"`
	BatchWindow   string `env:"BATCH_WINDOW" envDefault:"60s"`
	ReadLimit     int    `env:"READ_LIMIT" envDefault:"1000"`
}

func (s BatchProcessorConfig) FlushIntervalDuration() (time.Duration, error) {
	return time.ParseDuration(s.FlushInterval)
}

func (s BatchProcessorConfig) BatchWindowDuration() (time.Duration, error) {
	return time.ParseDuration(s.BatchWindow)
}
