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
