package rest

import (
	"github.com/labstack/echo/v4"

	"github.com/rasulov-emirlan/poc-auth/internal/domains/auth"
)

type authHandler struct {
	service auth.Service
}

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h authHandler) Login(ctx echo.Context) error {
	var req AuthLoginRequest

	if err := ctx.Bind(&req); err != nil {
		return responsError(ctx, err)
	}

	if err := h.service.Login(ctx.Request().Context(), req.Email, req.Password); err != nil {
		return responsError(ctx, err)
	}

	return ctx.NoContent(200)
}
