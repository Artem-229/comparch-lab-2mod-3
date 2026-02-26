package logging

import (
	"log/slog"
	"os"
)

func InitLogger(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	return logger
}

func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
