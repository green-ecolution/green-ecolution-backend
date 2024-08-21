package keycloak

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
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

func (r *KeycloakRepository) CreateUser(ctx context.Context, user *auth.User, password string, roles *[]string) (*auth.User, error) {
  slog.Debug("Creating user in keycloak", "user", user)
	keyCloakUser := userToKeyCloakUser(user, password)

	token, err := r.loginRestAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)
	userID, err := client.CreateUser(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, *keyCloakUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	err = client.SetPassword(ctx, token.AccessToken, userID, r.cfg.KeyCloak.Realm, password, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set password")
	}

	kcRoles := make([]gocloak.Role, len(*roles))
	for _, roleName := range *roles {
		roleNameLowerCase := strings.ToLower(roleName)
		roleKeyCloak, err := client.GetRealmRole(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, roleNameLowerCase)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get role by name: '%v'", roleNameLowerCase))
		}
		kcRoles = append(kcRoles, *roleKeyCloak)
	}

	err = client.AddRealmRoleToUser(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, userID, kcRoles)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to roles to user. roles '%v'", kcRoles))
	}

	userKeyCloak, err := client.GetUserByID(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, userID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get created user by id: '%v'", userID))
	}

	return keyCloakUserToUser(userKeyCloak), nil
}

func keyCloakUserToUser(user *gocloak.User) *auth.User {
	return &auth.User{
		ID:          *user.ID,
		Username:    *user.Username,
		FirstName:   *user.FirstName,
		LastName:    *user.LastName,
		Email:       *user.Email,
	}
}

func userToKeyCloakUser(user *auth.User, password string) *gocloak.User {
	attribute := make(map[string][]string)
	attribute["phone_number"] = []string{user.PhoneNumber}
	attribute["employee_id"] = []string{user.EmployeeID}

	return &gocloak.User{
		Username:   gocloak.StringP(user.Username),
		FirstName:  gocloak.StringP(user.FirstName),
		LastName:   gocloak.StringP(user.LastName),
		Email:      gocloak.StringP(user.Email),
		Enabled:    gocloak.BoolP(true),
		Attributes: &attribute,
	}
}
