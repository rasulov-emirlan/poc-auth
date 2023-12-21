package auth

import (
	"log/slog"

	"github.com/rasulov-emirlan/poc-auth/config"
)

type ServiceConfigs struct {
	UsersRepository UsersRepository
	Logger          *slog.Logger
	Cfg             config.Config
}

type Session struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
