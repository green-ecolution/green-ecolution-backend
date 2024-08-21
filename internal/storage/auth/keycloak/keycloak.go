package keycloak

import (
	"context"
	"fmt"
	"net/url"
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

func (r *KeycloakRepository) CreateUser(ctx context.Context, user *auth.User, password, role string) (*auth.User, error) {
	keyCloakUser := userToKeyCloakUser(user, password, role)

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

	roleNameLowerCase := strings.ToLower(role)
	roleKeyCloak, err := client.GetRealmRole(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, roleNameLowerCase)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get role by name: '%v'", roleNameLowerCase))
	}
	err = client.AddRealmRoleToUser(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, userID, []gocloak.Role{*roleKeyCloak})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to add role to user: '%v'", roleNameLowerCase))
	}

  userKeyCloak, err := client.GetUserByID(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, userID)
  if err != nil {
    return nil, errors.Wrap(err, fmt.Sprintf("failed to get created user by id: '%v'", userID))
  }

  return keyCloakUserToUser(userKeyCloak), nil
}

func keyCloakUserToUser(user *gocloak.User) *auth.User {
  return &auth.User{
    ID:             *user.ID,
    Username:       *user.Username,
    FirstName:      *user.FirstName,
    LastName:       *user.LastName,
    Email:          *user.Email,
    PhoneNumber:    (*user.Attributes)["phone_number"][0],
    EmployeeID:     (*user.Attributes)["employee_id"][0],
    ProfileImageURL: &url.URL{
      Host: (*user.Attributes)["profile_image"][0],
    },
  }
}

func userToKeyCloakUser(user *auth.User, password, role string) *gocloak.User {
	attribute := make(map[string][]string)
	attribute["phone_number"] = []string{user.PhoneNumber}
	attribute["employee_id"] = []string{user.EmployeeID}
	attribute["profile_image"] = []string{user.ProfileImageURL.Host}

	return &gocloak.User{
		Username:   &user.Username,
		FirstName:  &user.FirstName,
		LastName:   &user.LastName,
		Email:      &user.Email,
		Enabled:    gocloak.BoolP(true),
		Attributes: &attribute,
	}
}
