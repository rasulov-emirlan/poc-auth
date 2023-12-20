package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/rasulov-emirlan/poc-auth/config"
	"github.com/rasulov-emirlan/poc-auth/internal/domains/auth"
)

type server struct {
	srvr *http.Server
}

type ServerConfigs struct {
	Cfg        config.Config
	AuthDomain auth.Service
}

func NewServer(cfg ServerConfigs) server {
	srvr := http.Server{
		Addr:         cfg.Cfg.Server.Port,
		ReadTimeout:  cfg.Cfg.Server.TimeoutRead,
		WriteTimeout: cfg.Cfg.Server.TimeoutWrite,
	}

	router := echo.New()
	router.Use(middleware.Gzip())
	router.Use(middleware.CORS())
	router.Use(middlewareRequestID)
	router.Use(slogecho.New(slog.Default()))
	router.Use(middleware.RemoveTrailingSlash())

	authHandler := authHandler{
		service: cfg.AuthDomain,
	}

	router.POST("/auth/login", authHandler.Login)

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
