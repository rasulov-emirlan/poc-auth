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
	if len(password) < 8 {
		return Session{}, ErrPasswordTooShort
	}
	if len(firstname) < 1 || len(lastname) < 1 {
		return Session{}, ErrFirstnameOrLastnameTooShort
	}
	var user fusionauth.RegistrationRequest

	user.Registration = fusionauth.UserRegistration{
		ApplicationId: s.applicationId,
	}
	user.User.Email = email
	user.User.Password = password
	user.User.FirstName = firstname
	user.User.LastName = lastname

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

func (s Service) ForgotPassword(ctx context.Context, email string) error {
	var req fusionauth.ForgotPasswordRequest

	req.ApplicationId = s.applicationId
	req.Email = email

	res, errors, err := s.fusionClient.ForgotPassword(req)
	if err != nil {
		s.log.DebugContext(ctx, "failed to forgot password", "error", err)
		return fmt.Errorf("failed to forgot password: %w", err)
	}

	if errors != nil {
		s.log.DebugContext(ctx, "failed to forgot password", "error", errors)
		return fmt.Errorf("failed to forgot password: %s", errors.Error())
	}

	s.log.DebugContext(ctx, "Forgot password", "email", email, "change_id", res.ChangePasswordId)

	return nil
}

func (s Service) ResetPassword(ctx context.Context, password, token string) error {
	var req fusionauth.ChangePasswordRequest

	req.ApplicationId = s.applicationId
	req.Password = password

	res, errors, err := s.fusionClient.ChangePassword(token, req)
	if err != nil {
		s.log.DebugContext(ctx, "failed to reset password", "error", err)
		return fmt.Errorf("failed to reset password: %w", err)
	}

	if errors != nil {
		s.log.DebugContext(ctx, "failed to reset password", "error", errors)
		return fmt.Errorf("failed to reset password: %s", errors.Error())
	}

	s.log.DebugContext(ctx, "Reset password", "response", res)

	return nil
}

func (s Service) RefreshToken(ctx context.Context, refreshToken string) (Session, error) {
	var req fusionauth.RefreshRequest

	req.RefreshToken = refreshToken

	res, errors, err := s.fusionClient.ExchangeRefreshTokenForJWTWithContext(ctx, req)
	if err != nil {
		s.log.DebugContext(ctx, "failed to refresh token", "error", err)
		return Session{}, fmt.Errorf("failed to refresh token: %w", err)
	}

	if errors != nil {
		s.log.DebugContext(ctx, "failed to refresh token", "error", errors)
		return Session{}, fmt.Errorf("failed to refresh token: %s", errors.Error())
	}

	s.log.DebugContext(ctx, "Refreshed token", "response", res)

	return Session{res.Token, res.RefreshToken}, nil
}

func (s Service) VerifyToken(ctx context.Context, tokenString string) (entities.User, error) {
	res, errs, err := s.fusionClient.RetrieveUserUsingJWTWithContext(ctx, tokenString)
	if err != nil {
		s.log.DebugContext(ctx, "Failed to verify token", "error", err.Error(), "errors", errs.Error())
		return entities.User{}, fmt.Errorf("failed to verify token: %w", err)
	}

	if errs != nil {
		s.log.DebugContext(ctx, "Failed to verify token", "error", errs.Error())
		return entities.User{}, fmt.Errorf("failed to verify token: %s", errs.Error())
	}

	s.log.DebugContext(ctx, "Verified token", "response", res)

	if res.User.Email != "" {
		return entities.User{
			Email:     res.User.Email,
			Firstname: res.User.FirstName,
			Lastname:  res.User.LastName,
		}, nil
	}

	return entities.User{}, fmt.Errorf("invalid token")
}
