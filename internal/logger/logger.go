package logger

import (
	"io"
	"log/slog"
)

type LogFormat string

const (
	JSON LogFormat = "json"
	Text LogFormat = "text"
)

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

func (l LogLevel) ToSLog() slog.Level {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func CreateLogger(out io.Writer, format LogFormat, level LogLevel) *slog.Logger {
	options := &slog.HandlerOptions{
		AddSource: true,
		Level:     level.ToSLog(),
	}

	var handler slog.Handler
	switch format {
	case JSON:
		handler = slog.NewJSONHandler(out, options)
	case Text:
		handler = slog.NewTextHandler(out, options)
	default:
		handler = slog.NewTextHandler(out, options)
	}

  return slog.New(handler)
}
