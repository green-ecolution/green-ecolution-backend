package keycloak

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
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
	slog.Debug("Creating user in keycloak", "user", user, "roles", roles)
	if user == nil {
		return nil, ErrEmptyUser
	}

	keyCloakUser := userToKeyCloakUser(user)

	client, clientToken, err := loginRestAPIClient(ctx, r.cfg.OidcProvider.BaseURL, r.cfg.OidcProvider.Backend.ClientID, r.cfg.OidcProvider.Backend.ClientSecret, r.cfg.OidcProvider.DomainName)
	if err != nil {
		return nil, err
	}

	userID, err := client.CreateUser(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, *keyCloakUser)
	if err != nil {
		return nil, errors.Join(err, ErrCreateUser)
	}

	if err = client.SetPassword(ctx, clientToken.AccessToken, userID, r.cfg.OidcProvider.DomainName, password, false); err != nil {
		return nil, errors.Join(err, ErrSetPassword)
	}

	if roles != nil || len(roles) > 0 {
		kcRoles := make([]gocloak.Role, len(roles))
		for i, roleName := range roles {
			var roleKeyCloak *gocloak.Role
			roleNameLowerCase := strings.ToLower(roleName)
			roleKeyCloak, err = client.GetRealmRole(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, roleNameLowerCase)
			if err != nil {
				slog.Error("failed to get role by name", "role", roleNameLowerCase, "err", err)
				return nil, errors.Join(err, ErrGetRole)
			}

			if roleKeyCloak != nil {
				kcRoles[i] = *roleKeyCloak
			}
		}
		if err = client.AddRealmRoleToUser(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, userID, kcRoles); err != nil {
			slog.Error("failed to assign roles to user", "role", kcRoles, "err", err)
			return nil, errors.Join(err, ErrSetRole)
		}
	}

	userKeyCloak, err := client.GetUserByID(ctx, clientToken.AccessToken, r.cfg.OidcProvider.DomainName, userID)
	if err != nil {
		slog.Error("failed to get recently created user by id", "user_id", userID, "err", err)
		return nil, errors.Join(err, ErrGetUser)
	}

	return keyCloakUserToUser(userKeyCloak)
}

func (r *UserRepository) RemoveSession(ctx context.Context, refreshToken string) error {
	client := gocloak.NewClient(r.cfg.OidcProvider.BaseURL)
	if err := client.Logout(ctx, r.cfg.OidcProvider.Frontend.ClientID, r.cfg.OidcProvider.Frontend.ClientSecret, r.cfg.OidcProvider.DomainName, refreshToken); err != nil {
		return errors.Join(err, ErrLogout)
	}
	return nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	client, token, err := loginRestAPIClient(ctx, r.cfg.OidcProvider.BaseURL, r.cfg.OidcProvider.Backend.ClientID, r.cfg.OidcProvider.Backend.ClientSecret, r.cfg.OidcProvider.DomainName)
	if err != nil {
		return nil, err
	}

	users, err := client.GetUsers(ctx, token.AccessToken, r.cfg.OidcProvider.DomainName, gocloak.GetUsersParams{})
	if err != nil {
		return nil, errors.Join(err, ErrGetUser)
	}

	allUsers := make([]*entities.User, len(users))
	for i, kcUser := range users {
		user, err := keyCloakUserToUser(kcUser)
		if err != nil {
			return nil, err
		}
		allUsers[i] = user
	}

	return allUsers, nil
}

func (r *UserRepository) GetByIDs(ctx context.Context, ids []string) ([]*entities.User, error) {
	client, token, err := loginRestAPIClient(ctx, r.cfg.OidcProvider.BaseURL, r.cfg.OidcProvider.Backend.ClientID, r.cfg.OidcProvider.Backend.ClientSecret, r.cfg.OidcProvider.DomainName)
	if err != nil {
		return nil, err
	}

	users := make([]*entities.User, len(ids))
	for i, id := range ids {
		kcUser, err := client.GetUserByID(ctx, token.AccessToken, r.cfg.OidcProvider.DomainName, id)
		if err != nil {
			return nil, err
		}

		user, err := keyCloakUserToUser(kcUser)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

func keyCloakUserToUser(user *gocloak.User) (*entities.User, error) {
	userID, err := uuid.Parse(*user.ID)
	if err != nil {
		slog.Error("failed to parse user id", "user_id", *user.ID, "err", err)
		return nil, ErrParseID
	}
	var phoneNumber, employeeID, drivingLicenseClass string
	var userRoles []string
	if user.Attributes != nil {
		if val, ok := (*user.Attributes)["phone_number"]; ok && len(val) > 0 {
			phoneNumber = val[0]
		}

		if val, ok := (*user.Attributes)["employee_id"]; ok && len(val) > 0 {
			employeeID = val[0]
		}

		if val, ok := (*user.Attributes)["driving_license_class"]; ok && len(val) > 0 {
			drivingLicenseClass = val[0]
		}

		if val, ok := (*user.Attributes)["user_roles"]; ok && len(val) > 0 {
			userRoles = val

		}
	}
	var roles []entities.Role
	for _, roleName := range userRoles {
		roles = append(roles, entities.Role{Name: entities.ParseUserRole(roleName)})
	}

	const millisecondsInSecond = 1000
	return &entities.User{
		ID:             userID,
		CreatedAt:      time.Unix(*user.CreatedTimestamp/millisecondsInSecond, 0),
		Username:       *user.Username,
		FirstName:      *user.FirstName,
		LastName:       *user.LastName,
		Email:          *user.Email,
		PhoneNumber:    phoneNumber,
		EmployeeID:     employeeID,
		Roles:          roles,
		DrivingLicense: entities.DrivingLicense(drivingLicenseClass),
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
