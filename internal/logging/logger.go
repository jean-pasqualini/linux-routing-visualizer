package logging

import (
	"log/slog"
	"os"
	"time"
)
import "github.com/lmittmann/tint"

type Logger = slog.Logger

func New(name string) *Logger {
	consoleHandler := tint.NewHandler(os.Stdout, &tint.Options{
		TimeFormat: time.RFC3339,
		Level:      slog.LevelDebug,
	})

	return slog.New(consoleHandler)
}

func NewFilelogger() *Logger {
	file, _ := os.OpenFile("/app/log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	fileHandler := tint.NewHandler(file, &tint.Options{
		TimeFormat: time.RFC3339,
		Level:      slog.LevelDebug,
	})

	return slog.New(fileHandler)
}
