package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var levels = map[string]slog.Level{
	"dev":   slog.LevelDebug,
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func NewLogger(logLevel string) *slog.Logger {
	level, ok := levels[logLevel]
	if !ok {
		level = slog.LevelInfo
	}
	var l *slog.Logger
	options := slog.HandlerOptions{
		Level: level,
	}

	if logLevel == "dev" {
		l = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level: level,
		}))
	} else {
		l = slog.New(slog.NewTextHandler(os.Stdout, &options))
	}
	slog.SetDefault(l)
	return l
}
