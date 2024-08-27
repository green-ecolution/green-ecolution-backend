package auth

import (
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/auth/keycloak"
)

func NewRepository(cfg *config.IdentityAuthConfig) *storage.Repository {
	authRepo := keycloak.NewKeycloakRepository(cfg)

	return &storage.Repository{
		Auth: authRepo,
	}
}
