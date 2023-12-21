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
		usersRepo     UsersRepository
		fusionClient  *fusionauth.FusionAuthClient
		applicationId string
		log           *slog.Logger
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
		usersRepo:     cfg.UsersRepository,
		fusionClient:  authClient,
		applicationId: cfg.Cfg.FusionAuth.AppId,
		log:           cfg.Logger,
	}, nil
}

func (s Service) Login(ctx context.Context, email, password string) (Session, error) {
	var credentials fusionauth.LoginRequest

	credentials.LoginId = email
	credentials.Password = password

	authResponse, errors, err := s.fusionClient.Login(credentials)
	if err != nil {
		s.log.DebugContext(ctx, "failed to login", "error", err)
		return Session{}, fmt.Errorf("failed to login: %w", err)
	}

	if errors != nil {
		s.log.DebugContext(ctx, "failed to login", "error", errors)
		return Session{}, fmt.Errorf("failed to login: %s", errors.Error())
	}

	s.log.DebugContext(ctx, "Logged user in fusionauth", "email", email, "response", authResponse.Token)

	return Session{authResponse.Token, authResponse.RefreshToken}, nil
}

func (s Service) Register(ctx context.Context, email, password, firstname, lastname string) (Session, error) {
	var user fusionauth.RegistrationRequest

	user.Registration = fusionauth.UserRegistration{
		ApplicationId: s.applicationId,
	}
	user.User.Email = email
	user.User.Password = password

	res, errors, err := s.fusionClient.Register("", user)
	if err != nil {
		s.log.DebugContext(ctx, "failed to register", "error", err)
		return Session{}, fmt.Errorf("failed to register: %w", err)
	}

	if errors != nil {
		s.log.DebugContext(ctx, "failed to register", "error", errors)
		return Session{}, fmt.Errorf("failed to register: %s", errors.Error())
	}

	s.log.DebugContext(ctx, "Registered user in fusionauth", "email", email, "response", res)

	u, err := s.usersRepo.Create(ctx, entities.User{
		Email:     email,
		Firstname: firstname,
		Lastname:  lastname,
	})
	if err != nil {
		if res, _, err := s.fusionClient.DeleteUser(res.User.Id); err != nil {
			s.log.ErrorContext(ctx, "Failed to delete user", "error", err, "response", res)
		}
		s.log.DebugContext(ctx, "Failed to create user in db, deleted it in fustionauth", "error", err)
		return Session{}, fmt.Errorf("failed to create user: %w", err)
	}

	s.log.DebugContext(ctx, "Created user in database", "email", email, "user", u)

	return Session{res.Token, res.RefreshToken}, nil
}
