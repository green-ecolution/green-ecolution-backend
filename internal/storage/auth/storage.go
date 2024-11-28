package auth

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/auth/keycloak"
)

func NewRepository(cfg *config.IdentityAuthConfig) *storage.Repository {
	authRepo := keycloak.NewKeycloakRepository(cfg)
	userRepo := keycloak.NewUserRepository(cfg)

	return &storage.Repository{
		Auth: authRepo,
		User: userRepo,
	}
}
