package keycloak

import (
	"context"
	"errors"

	"github.com/Nerzal/gocloak/v13"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotActive   = errors.New("token not active")
	ErrTokenInvalidType = errors.New("token invalid type")
	ErrLogin            = errors.New("failed to login to keycloak")
	ErrLogout           = errors.New("failed to logout")
	ErrEmptyUser        = errors.New("user is nil")
	ErrGetUser          = errors.New("failed to get user")
	ErrCreateUser       = errors.New("failed to create user")
	ErrSetPassword      = errors.New("failed to set password")
	ErrGetRole          = errors.New("failed to get role by name")
	ErrSetRole          = errors.New("failed to assign role")
	ErrParseID          = errors.New("failed to parse user id")
)

type KeycloakRepository struct {
	cfg *config.IdentityAuthConfig
}

func NewKeycloakRepository(cfg *config.IdentityAuthConfig) storage.AuthRepository {
	return &KeycloakRepository{
		cfg: cfg,
	}
}

func loginRestAPIClient(ctx context.Context, baseURL, clientID, clientSecret, realm string) (client *gocloak.GoCloak, token *gocloak.JWT, err error) {
	client = gocloak.NewClient(baseURL)

	token, err = client.LoginClient(ctx, clientID, clientSecret, realm)
	if err != nil {
		return nil, nil, errors.Join(err, ErrLogin)
	}

	return
}
