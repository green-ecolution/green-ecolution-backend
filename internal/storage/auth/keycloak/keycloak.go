package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/pkg/errors"
)

type KeycloakRepository struct {
	cfg *config.IdentityAuthConfig
}

func NewKeycloakRepository(cfg *config.IdentityAuthConfig) storage.AuthRepository {
	return &KeycloakRepository{
		cfg: cfg,
	}
}

func loginRestAPIClient(ctx context.Context, baseURL, clientID, clientSecret, realm string) (*gocloak.JWT, error) {
	client := gocloak.NewClient(baseURL)

	token, err := client.LoginClient(ctx, clientID, clientSecret, realm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login to keycloak")
	}
	return token, nil
}
