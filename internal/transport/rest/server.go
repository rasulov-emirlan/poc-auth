package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/rasulov-emirlan/poc-auth/config"
)

type server struct {
	srvr *http.Server
}

func NewServer(cfg config.Config) server {
	srvr := http.Server{
		Addr:         cfg.Server.Port,
		ReadTimeout:  cfg.Server.TimeoutRead,
		WriteTimeout: cfg.Server.TimeoutWrite,
	}

	router := echo.New()
	router.Use(middleware.Gzip())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(slogecho.New(slog.Default()))
	router.Use(middleware.RemoveTrailingSlash())

	srvr.Handler = router

	return server{
		srvr: &srvr,
	}
}

func (s server) Start() error {
	if err := s.srvr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s server) Shutdown(ctx context.Context) error {
	return s.srvr.Shutdown(ctx)
}
