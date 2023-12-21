package rest

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/rasulov-emirlan/poc-auth/internal/domains/auth"
)

const refreshTokenKey = "refresh_token"

type authHandler struct {
	service auth.Service
}

type (
	AuthLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AuthRegisterRequest struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	}

	AuthLoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

func (h authHandler) Login(ctx echo.Context) error {
	var req AuthLoginRequest

	if err := ctx.Bind(&req); err != nil {
		return responsError(ctx, err)
	}

	res, err := h.service.Login(ctx.Request().Context(), req.Email, req.Password)
	if err != nil {
		return responsError(ctx, err)
	}

	ctx.SetCookie(&http.Cookie{
		Name:     refreshTokenKey,
		Value:    res.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HttpOnly: true,
		Secure:   true,
	})

	return ctx.JSON(200, AuthLoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h authHandler) Register(ctx echo.Context) error {
	var req AuthRegisterRequest

	if err := ctx.Bind(&req); err != nil {
		return responsError(ctx, err)
	}

	res, err := h.service.Register(ctx.Request().Context(), req.Email, req.Password, req.Firstname, req.Lastname)
	if err != nil {
		return responsError(ctx, err)
	}

	ctx.SetCookie(&http.Cookie{
		Name:     refreshTokenKey,
		Value:    res.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HttpOnly: true,
		Secure:   true,
	})

	return ctx.JSON(200, AuthLoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}
