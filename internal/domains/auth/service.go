package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/rasulov-emirlan/poc-auth/internal/entities"
)

type (
	UsersRepository interface {
		Create(ctx context.Context, user entities.User) (entities.User, error)
		GetByEmail(ctx context.Context, email string) (entities.User, error)
		Update(ctx context.Context, user entities.User) (entities.User, error)
	}

	Service struct {
		usersRepo    UsersRepository
		fusionClient *fusionauth.FusionAuthClient
		log          *slog.Logger
	}
)

func NewService(ctx context.Context, cfg ServiceConfigs) (Service, error) {
	httpclient := &http.Client{
		Timeout: time.Second * 10,
	}
	baseUrl, err := url.Parse(cfg.Cfg.FusionAuth.Host)
	if err != nil {
		return Service{}, fmt.Errorf("failed to parse fusion auth host: %w", err)
	}

	authClient := fusionauth.NewClient(httpclient, baseUrl, cfg.Cfg.FusionAuth.ApiKey)

	return Service{
		usersRepo:    cfg.UsersRepository,
		fusionClient: authClient,
		log:          cfg.Logger,
	}, nil
}

func (s Service) Login(ctx context.Context, email, password string) error {
	var credentials fusionauth.LoginRequest

	credentials.LoginId = email
	credentials.Password = password

	authResponse, errors, err := s.fusionClient.Login(credentials)
	if err != nil {
		s.log.DebugContext(ctx, "failed to login", "error", err)
		return fmt.Errorf("failed to login: %w", err)
	}

	if errors != nil {
		s.log.DebugContext(ctx, "failed to login", "error", errors)
		return fmt.Errorf("failed to login: %s", errors.Error())
	}

	s.log.InfoContext(ctx, "authResponse", "response", authResponse)

	return nil
}
