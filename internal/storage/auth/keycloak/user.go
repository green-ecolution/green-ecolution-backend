package keycloak

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
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

	clientToken, err := loginRestAPIClient(ctx, r.cfg.OidcProvider.BaseURL, r.cfg.OidcProvider.Backend.ClientID, r.cfg.OidcProvider.Backend.ClientSecret, r.cfg.OidcProvider.DomainName)
	if err != nil {
		return nil, err
	}

	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)
	userID, err := client.CreateUser(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, *keyCloakUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	if err = client.SetPassword(ctx, clientToken.AccessToken, userID, r.cfg.OidcProvider.DomainName, password, false); err != nil {
		return nil, errors.Wrap(err, "failed to set password")
	}

	if roles != nil || len(roles) > 0 {
		kcRoles := make([]gocloak.Role, len(roles))
		for i, roleName := range roles {
			var roleKeyCloak *gocloak.Role
			roleNameLowerCase := strings.ToLower(roleName)
			roleKeyCloak, err = client.GetRealmRole(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, roleNameLowerCase)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to get role by name: '%v'", roleNameLowerCase))
			}

			if roleKeyCloak != nil {
				kcRoles[i] = *roleKeyCloak
			}
		}
		if err = client.AddRealmRoleToUser(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, userID, kcRoles); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to roles to user. roles '%v'", kcRoles))
		}
	}

	userKeyCloak, err := client.GetUserByID(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, userID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get created user by id: '%v'", userID))
	}

	return keyCloakUserToUser(userKeyCloak)
}

func (r *UserRepository) RemoveSession(ctx context.Context, refreshToken string) error {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)
	if err := client.Logout(ctx, r.cfg.OidcProvider.Frontend.ClientID, r.cfg.OidcProvider.Frontend.ClientSecret, r.cfg.OidcProvider.DomainName, refreshToken); err != nil {
		return errors.Wrap(err, "failed to logout")
	}
	return nil
}

func (r *UserRepository) getClientToken(ctx context.Context) (string, error) {
	clientToken, err := loginRestAPIClient(
		ctx,
		r.cfg.OidcProvider.BaseURL,
		r.cfg.OidcProvider.Backend.ClientID,
		r.cfg.OidcProvider.Backend.ClientSecret,
		r.cfg.OidcProvider.DomainName,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to log in to Keycloak")
	}
	return clientToken.AccessToken, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)
	token, err := r.getClientToken(ctx)
	if err != nil {
		return nil, err
	}

	users, err := client.GetUsers(ctx, token, r.cfg.OidcProvider.DomainName, gocloak.GetUsersParams{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users from Keycloak")
	}

	allUsers := make([]*entities.User, 0, len(users))
	for _, kcUser := range users {
		user, err := keyCloakUserToUser(kcUser)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert Keycloak user")
		}
		allUsers = append(allUsers, user)
	}

	return allUsers, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)
	token, err := r.getClientToken(ctx)
	if err != nil {
		return nil, err
	}

	kcUser, err := client.GetUserByID(ctx, token, r.cfg.OidcProvider.DomainName, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from Keycloak")
	}

	user, err := keyCloakUserToUser(kcUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert Keycloak user")
	}

	return user, nil
}

func keyCloakUserToUser(user *gocloak.User) (*entities.User, error) {
	userID, err := uuid.Parse(*user.ID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse user id: '%v'", *user.ID))
	}
	var phoneNumber, employeeID string
	if user.Attributes != nil {
		if val, ok := (*user.Attributes)["phone_number"]; ok && len(val) > 0 {
			phoneNumber = val[0]
		}

		if val, ok := (*user.Attributes)["employee_id"]; ok && len(val) > 0 {
			employeeID = val[0]
		}
	}
	const millisecondsInSecond = 1000
	return &entities.User{
		ID:          userID,
		CreatedAt:   time.Unix(*user.CreatedTimestamp/millisecondsInSecond, 0),
		Username:    *user.Username,
		FirstName:   *user.FirstName,
		LastName:    *user.LastName,
		Email:       *user.Email,
		PhoneNumber: phoneNumber,
		EmployeeID:  employeeID,
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
