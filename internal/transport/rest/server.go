//go:generate swag init -g ./internal/transport/rest/server.go
package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/rasulov-emirlan/poc-auth/docs"
	slogecho "github.com/samber/slog-echo"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/rasulov-emirlan/poc-auth/config"
	"github.com/rasulov-emirlan/poc-auth/internal/domains/auth"
)

// @title POC-Auth API
// @description This is a sample server for POC-Auth API.
// @version 1.0
// @BasePath /

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

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	slog.Default().DebugContext(context.Background(), "Server started on port "+cfg.Cfg.Server.Port)

	authHandler := authHandler{
		service: cfg.AuthDomain,
	}

	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/refresh", authHandler.Refresh)
	router.POST("/auth/forgot-password/:email", authHandler.ForgotPassword)
	router.POST("/auth/reset-password/:token", authHandler.ResetPassword)

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
