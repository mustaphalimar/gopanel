package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mustaphalimar/gopanel/pkg/config"
)

type Logger struct {
	*slog.Logger
}

func New(cfg *config.LoggerConfig) *Logger {
	var handler slog.Handler
	var level slog.Level

	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// create handler based on format
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.Level == "debug",
	}

	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	case "color", "colored":
		handler = tint.NewTextHandler(os.Stdout, &tint.Options{
			Level:      level,
			AddSource:  cfg.Level == "debug",
			TimeFormat: time.Kitchen,
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

func NewWithWriter(w io.Writer, format string, level slog.Level) *Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch format {
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	case "color", "colored":
		handler = tint.NewTextHandler(w, &tint.Options{
			Level:      level,
			TimeFormat: time.Kitchen,
		})
	default:
		handler = slog.NewTextHandler(w, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args),
	}
}

func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		Logger: l.Logger.WithGroup(name),
	}
}

func (l *Logger) Request(method, path string, statusCode int, duration string, attrs ...any) {
	l.Info(
		"HTTP Request",
		append([]any{
			"method", method,
			"path", path,
			"status", statusCode,
			"duration", duration,
		}, attrs...)...,
	)
}

func (l *Logger) ErrorWithContext(msg string, err error, attrs ...any) {
	l.Error(
		msg,
		append(
			[]any{
				"error", err,
			}, attrs...,
		)...,
	)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}
