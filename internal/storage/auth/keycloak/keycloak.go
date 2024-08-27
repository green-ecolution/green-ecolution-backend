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

func (r *KeycloakRepository) loginRestAPIClient(ctx context.Context) (*gocloak.JWT, error) {
	client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)

	token, err := client.LoginClient(ctx, r.cfg.KeyCloak.ClientID, r.cfg.KeyCloak.ClientSecret, r.cfg.KeyCloak.Realm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to login to keycloak")
	}
	return token, nil
}
