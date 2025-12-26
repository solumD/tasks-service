package logger

import (
	"log/slog"
	"os"
)

const (
	levelDebug string = "debug"
	levelInfo  string = "info"
	levelWarn  string = "warn"
	levelError string = "error"
)

// NewLogger возвращает новый логгер согласно уровню логирования
func NewLogger(loggerLevel string) *slog.Logger {
	var log *slog.Logger

	switch loggerLevel {
	case levelDebug:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case levelInfo:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	case levelWarn:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelWarn},
			),
		)
	case levelError:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelError},
			),
		)
	}

	return log
}

// Error возвращает аттрибут с ошибкой (обертка)
func Error(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// String возвращает аттрибут со строкой (обертка)
func String(key string, value string) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.StringValue(value),
	}
}

// Int возвращает аттрибут с числом (обертка)
func Int(key string, value int) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.IntValue(value),
	}
}

// Any возвращает аттрибут с любым значением (обертка)
func Any(key string, value any) slog.Attr {
	return slog.Attr{
		Key:   key,
		Value: slog.AnyValue(value),
	}
}
