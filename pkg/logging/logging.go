package logging

import (
	"context"
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

type key string

const ReqIdKey key = "request-id"

type handler struct {
	slog.Handler
}

func (h handler) Handle(ctx context.Context, r slog.Record) error {
	if reqId, ok := ctx.Value(ReqIdKey).(string); ok {
		r.Add(string(ReqIdKey), slog.StringValue(reqId))
	}
	return h.Handler.Handle(ctx, r)
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
		l = slog.New(handler{tint.NewHandler(os.Stdout, &tint.Options{
			Level: level,
		})})
	} else {
		l = slog.New(handler{slog.NewTextHandler(os.Stdout, &options)})
	}
	slog.SetDefault(l)
	return l
}
