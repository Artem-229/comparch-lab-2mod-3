package logging

import (
	"io"
	"log/slog"
	"os"
)

var (
	emptyLogger = slog.New(slog.NewTextHandler(io.Discard, nil))
)

func InitLogger(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	return logger
}

func EmptyLogger() *slog.Logger {
	return emptyLogger
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
