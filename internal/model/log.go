package model

import "time"

const (
	LogLevelError   LogLevel = "error"
	LogLevelWarn    LogLevel = "warn"
	LogLevelInfo    LogLevel = "info"
	LogLevelDebug   LogLevel = "debug"
	LogLevelTrace   LogLevel = "trace"
	LogLevelFatal   LogLevel = "fatal"
	LogLevelUnknown LogLevel = "unknown"
)

type Log struct {
	UUID       string
	ItemUUID   string
	LaunchUUID string
	Timestamp  time.Time
	Level      LogLevel
	Message    string
	File       LogFile
}

type LogFile struct {
	Name string
}

type LogLevel string
