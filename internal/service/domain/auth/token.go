package auth

import (
	"context"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

func (s *AuthService) RetrospectToken(ctx context.Context, token string) (*domain.IntroSpectTokenResult, error) {
	log := logger.GetLogger(ctx)
	result, err := s.authRepository.RetrospectToken(ctx, token)
	if err != nil {
		log.Error("failed to retrospect token", "token", token, "error", err)
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to retrospect token").Error())
	}

	return result, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, token string) (*domain.ClientToken, error) {
	log := logger.GetLogger(ctx)
	result, err := s.authRepository.RefreshToken(ctx, token)
	if err != nil {
		log.Error("failed to refresh token", "token", token, "error", err)
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to refresh token").Error())
	}

	return result, nil
}
