package auth

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/pkg/errors"
)

type AuthService struct {
	authRepository storage.AuthRepository
	validator      *validator.Validate
	cfg            *config.Config
}

func NewAuthService(repo storage.AuthRepository, cfg *config.Config) service.AuthService {
	return &AuthService{
		validator:      validator.New(),
		authRepository: repo,
		cfg:            cfg,
	}
}

func (s *AuthService) RetrospectToken(ctx context.Context, token string) (*auth.IntroSpectTokenResult, error) {
	result, err := s.authRepository.RetrospectToken(ctx, token)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to retrospect token").Error())
	}

	return result, nil
}

func (s *AuthService) Register(ctx context.Context, user *auth.RegisterUser) (*auth.User, error) {
	if err := s.validator.Struct(user); err != nil {
		return nil, service.NewError(service.BadRequest, errors.Wrap(err, "validation error").Error())
	}

	createdUser, err := s.authRepository.CreateUser(ctx, &user.User, user.Password, user.Roles)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to create user").Error())
	}

	return createdUser, nil
}

func (s *AuthService) Ready() bool {
	return s.authRepository != nil
}
