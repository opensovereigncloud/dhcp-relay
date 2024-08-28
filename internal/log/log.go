package log

import (
	"log/slog"
	"os"
)

type Params struct {
	Level  string
	Format string
}

type Logger struct {
	*slog.Logger
}

func (l *Logger) Debug(msg string, keysAndValues ...any) {
	l.With(keysAndValues...).Debug(msg)
}

func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.With(keysAndValues...).Info(msg)
}

func (l *Logger) Warn(err error, msg string, keysAndValues ...any) {
	l.With(keysAndValues...).Warn(msg, slog.Any("error", err))
}

func (l *Logger) Error(err error, msg string, keysAndValues ...any) {
	l.With(keysAndValues...).Error(msg, slog.Any("error", err))
}

func Setup(params Params) *Logger {
	var l *slog.Logger
	level := mapLogLevel(params.Level)
	switch params.Format {
	case "json":
		l = setupJSONHandler(level)
	case "text":
		l = setupTextHandler(level)
	default:
		l = slog.Default()
	}
	return &Logger{l}
}

func setupJSONHandler(level slog.Leveler) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

func setupTextHandler(level slog.Leveler) *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}

func mapLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
