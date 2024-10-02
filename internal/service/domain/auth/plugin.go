package auth

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (s *AuthService) AuthPlugin(ctx context.Context, plugin *entities.AuthPlugin) (*entities.ClientToken, error) {
	token, err := s.authRepository.GetAccessTokenFromPassword(ctx, plugin.Username, plugin.Password)
	if err != nil {
		return nil, err
	}
	return token, nil
}
