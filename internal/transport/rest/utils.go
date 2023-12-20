package rest

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/poc-auth/internal/domains/auth"
	"github.com/rasulov-emirlan/poc-auth/pkg/logging"
)

func responsError(ctx echo.Context, err error) error {
	switch err {
	case auth.ErrEmailNotFound:
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	case auth.ErrEmailTaken:
		return ctx.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	default:
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}

func middlewareRequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		reqID := req.Header.Get(echo.HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.SetRequest(req.WithContext(req.Context()))
		c.Response().Header().Set(echo.HeaderXRequestID, reqID)
		ctx := req.Context()
		ctx = context.WithValue(ctx, logging.ReqIdKey, reqID)
		c.SetRequest(req.WithContext(ctx))
		return next(c)
	}
}
