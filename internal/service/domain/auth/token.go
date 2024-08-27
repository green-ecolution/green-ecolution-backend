package auth

import (
	"context"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

func (s *AuthService) RetrospectToken(ctx context.Context, token string) (*domain.IntroSpectTokenResult, error) {
	result, err := s.authRepository.RetrospectToken(ctx, token)
	if err != nil {
		return nil, service.NewError(service.InternalError, errors.Wrap(err, "failed to retrospect token").Error())
	}

	return result, nil
}
