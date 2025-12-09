package model

type LogLevel string

const (
	LogLevelError   LogLevel = "error"
	LogLevelWarn    LogLevel = "warn"
	LogLevelInfo    LogLevel = "info"
	LogLevelDebug   LogLevel = "debug"
	LogLevelTrace   LogLevel = "trace"
	LogLevelFatal   LogLevel = "fatal"
	LogLevelUnknown LogLevel = "unknown"
)
