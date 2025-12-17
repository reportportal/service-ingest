package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

type Options struct {
	Level     string
	Format    string
	AddSource bool
	Output    io.Writer
}

func New(opts Options) *slog.Logger {
	output := opts.Output
	if output == nil {
		output = os.Stdout
	}

	options := &slog.HandlerOptions{
		Level:     parseLevel(opts.Level),
		AddSource: opts.AddSource,
	}

	handler := getHandler(opts.Format, output, options)

	return slog.New(handler)
}

func parseLevel(levelStr string) slog.Level {
	var level slog.Level
	byteLevel := []byte(strings.ToLower(strings.TrimSpace(levelStr)))
	if err := level.UnmarshalText(byteLevel); err != nil {
		level = slog.LevelInfo
	}
	return level
}

func getHandler(format string, writer io.Writer, opts *slog.HandlerOptions) slog.Handler {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		return slog.NewJSONHandler(writer, opts)
	case "text":
		return slog.NewTextHandler(writer, opts)
	default:
		return slog.NewJSONHandler(writer, opts)
	}
}
