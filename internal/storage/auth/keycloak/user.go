package keycloak

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/pkg/errors"
)

type UserRepository struct {
	cfg *config.IdentityAuthConfig
}

func NewUserRepository(cfg *config.IdentityAuthConfig) storage.UserRepository {
	return &UserRepository{
		cfg: cfg,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User, password string, roles []string) (*entities.User, error) {
	slog.Debug("Creating user in keycloak", "user", user)
	if user == nil {
		return nil, errors.New("user is nil")
	}

	keyCloakUser := userToKeyCloakUser(user)

	clientToken, err := loginRestAPIClient(ctx, r.cfg.KeyCloak.BaseURL, r.cfg.KeyCloak.ClientID, r.cfg.KeyCloak.ClientSecret, r.cfg.KeyCloak.Realm)
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)
	userID, err := client.CreateUser(ctx, clientToken.AccessToken, r.cfg.KeyCloak.Realm, *keyCloakUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	if err = client.SetPassword(ctx, clientToken.AccessToken, userID, r.cfg.KeyCloak.Realm, password, false); err != nil {
		return nil, errors.Wrap(err, "failed to set password")
	}

	if roles != nil || len(roles) > 0 {
		kcRoles := make([]gocloak.Role, len(roles))
		for i, roleName := range roles {
			var roleKeyCloak *gocloak.Role
			roleNameLowerCase := strings.ToLower(roleName)
			roleKeyCloak, err = client.GetRealmRole(ctx, clientToken.AccessToken, r.cfg.KeyCloak.Realm, roleNameLowerCase)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to get role by name: '%v'", roleNameLowerCase))
			}

			if roleKeyCloak != nil {
				kcRoles[i] = *roleKeyCloak
			}
		}
		if err = client.AddRealmRoleToUser(ctx, clientToken.AccessToken, r.cfg.KeyCloak.Realm, userID, kcRoles); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to roles to user. roles '%v'", kcRoles))
		}
	}

	userKeyCloak, err := client.GetUserByID(ctx, clientToken.AccessToken, r.cfg.KeyCloak.Realm, userID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get created user by id: '%v'", userID))
	}

	return keyCloakUserToUser(userKeyCloak)
}

func (r *UserRepository) RemoveSession(ctx context.Context, refreshToken string) error {
	client := gocloak.NewClient(r.cfg.KeyCloak.BaseURL)
	if err := client.Logout(ctx, r.cfg.KeyCloak.Frontend.ClientID, r.cfg.KeyCloak.Frontend.ClientSecret, r.cfg.KeyCloak.Realm, refreshToken); err != nil {
		return errors.Wrap(err, "failed to logout")
	}
	return nil
}

func keyCloakUserToUser(user *gocloak.User) (*entities.User, error) {
	userID, err := uuid.Parse(*user.ID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse user id: '%v'", *user.ID))
	}
	var phone_number string
	var employee_id string
	if user.Attributes != nil {
		phone_number = (*user.Attributes)["phone_number"][0]
		employee_id = (*user.Attributes)["employee_id"][0]
	}

	return &entities.User{
		ID:          userID,
		CreatedAt:   time.Unix(*user.CreatedTimestamp, 0),
		Username:    *user.Username,
		FirstName:   *user.FirstName,
		LastName:    *user.LastName,
		Email:       *user.Email,
		PhoneNumber: phone_number,
		EmployeeID:  employee_id,
	}, nil
}

func userToKeyCloakUser(user *entities.User) *gocloak.User {
	attribute := make(map[string][]string)
	attribute["phone_number"] = []string{user.PhoneNumber}
	attribute["employee_id"] = []string{user.EmployeeID}

	return &gocloak.User{
    ID:         gocloak.StringP(user.ID.String()),
		Username:   gocloak.StringP(user.Username),
		FirstName:  gocloak.StringP(user.FirstName),
		LastName:   gocloak.StringP(user.LastName),
		Email:      gocloak.StringP(user.Email),
		Enabled:    gocloak.BoolP(true),
		Attributes: &attribute,
	}
}
