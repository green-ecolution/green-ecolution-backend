package keycloak

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/pkg/errors"
)

func (r *KeycloakRepository) CreateUser(ctx context.Context, user *entities.User, password string, roles *[]string) (*entities.User, error) {
	slog.Debug("Creating user in keycloak", "user", user)
	keyCloakUser := userToKeyCloakUser(user)

	token, err := r.loginRestAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)
	userID, err := client.CreateUser(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, *keyCloakUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	if err = client.SetPassword(ctx, token.AccessToken, userID, r.cfg.KeyCloak.Realm, password, false); err != nil {
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

	if err = client.AddRealmRoleToUser(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, userID, kcRoles); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to roles to user. roles '%v'", kcRoles))
	}

	userKeyCloak, err := client.GetUserByID(ctx, token.AccessToken, r.cfg.KeyCloak.Realm, userID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get created user by id: '%v'", userID))
	}

	return keyCloakUserToUser(userKeyCloak)
}

func keyCloakUserToUser(user *gocloak.User) (*entities.User, error) {
	userID, err := uuid.Parse(*user.ID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse user id: '%v'", *user.ID))
	}
	return &entities.User{
		ID:        userID,
		Username:  *user.Username,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
	}, nil
}

func userToKeyCloakUser(user *entities.User) *gocloak.User {
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
