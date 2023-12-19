package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/poc-auth/config"
	"github.com/rasulov-emirlan/poc-auth/internal/storage/mongodb"
	"github.com/rasulov-emirlan/poc-auth/pkg/logging"
)

type application struct {
	ctx    context.Context
	cancel context.CancelFunc

	logger *slog.Logger

	cfg config.Config

	mdb mongodb.RepoCombiner

	cleanupFuncs []func()
}

func Run() {
	a := application{}
	defer a.cleanup()

	if err := a.init(); err != nil {
		panic(err)
	}

	if err := a.initDB(); err != nil {
		a.fatal("failed to init db", err)
	}

	if err := a.initHttp(); err != nil {
		a.fatal("failed to init http", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	a.logger.InfoContext(a.ctx, "shutting down")
}

func (a *application) fatal(msg string, err error) {
	a.logger.ErrorContext(a.ctx, msg, "err", err.Error())
	a.cleanup()
	os.Exit(1)
}

func (a *application) cleanup() {
	for _, f := range a.cleanupFuncs {
		f()
	}
}

func (a *application) init() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx
	a.cancel = cancel

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	a.cfg = cfg

	a.logger = logging.NewLogger(cfg.LogLevel)

	return nil
}
