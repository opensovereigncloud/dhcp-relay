package log

import (
	"log/slog"
	"os"
)

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

func Setup() *Logger {
	return &Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}
}
