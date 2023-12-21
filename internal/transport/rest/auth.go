package rest

import (
	"net/http"
	"strings"
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

	return ctx.JSON(http.StatusOK, AuthLoginResponse{
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

	return ctx.JSON(http.StatusOK, AuthLoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h authHandler) ForgotPassword(ctx echo.Context) error {
	email := ctx.Param("email")

	if err := h.service.ForgotPassword(ctx.Request().Context(), email); err != nil {
		return responsError(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h authHandler) ResetPassword(ctx echo.Context) error {
	var req AuthLoginRequest

	if err := ctx.Bind(&req); err != nil {
		return responsError(ctx, err)
	}

	token := ctx.Param("token")

	if err := h.service.ResetPassword(ctx.Request().Context(), req.Email, req.Password, token); err != nil {
		return responsError(ctx, err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h authHandler) Refresh(ctx echo.Context) error {
	refreshToken, err := ctx.Cookie(refreshTokenKey)
	if err != nil {
		return responsError(ctx, err)
	}

	res, err := h.service.RefreshToken(ctx.Request().Context(), refreshToken.Value)
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

	return ctx.JSON(http.StatusOK, AuthLoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h authHandler) middlewareExtractUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		accessToken := ctx.Request().Header.Get("Authorization")
		if accessToken == "" {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "access token required"})
		}

		accessParts := strings.Split(accessToken, " ")
		if len(accessParts) != 2 {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid access token"})
		}

		if accessParts[0] != "Bearer" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "token has to be bearer"})
		}

		res, err := h.service.VerifyToken(ctx.Request().Context(), accessParts[1])
		if err != nil {
			return responsError(ctx, err)
		}

		ctx.Set("user", res)

		return next(ctx)
	}
}
