package auth

import (
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/auth/keycloak"
)

func NewRepository(cfg *config.IdentityAuthConfig) *storage.Repository {
	authRepo := keycloak.NewKeycloakRepository(cfg)
	slog.Info("successfully initialized auth repository", "service", "keycloak")

	userRepo := keycloak.NewUserRepository(cfg)
	slog.Info("successfully initialized user repository")

	return &storage.Repository{
		Auth: authRepo,
		User: userRepo,
	}
}
