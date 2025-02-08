package testutils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

func (s *KeycloakTestSuite) LoginAdminAndGetToken(t testing.TB) *gocloak.JWT {
	t.Helper()
	identityConfig := s.IdentityConfig(t, context.Background())
	client := gocloak.NewClient(identityConfig.OidcProvider.BaseURL)
	token, err := client.Login(context.Background(), identityConfig.OidcProvider.Backend.ClientID, identityConfig.OidcProvider.Backend.ClientSecret, identityConfig.OidcProvider.DomainName, s.User, s.Password)
	if err != nil {
		t.Fatalf("loginAdminAndGetToken::failed to get token: %v", err)
	}
	return token
}

func (s *KeycloakTestSuite) LoginUser(t testing.TB, user *entities.User) *gocloak.JWT {
	t.Helper()
	identityConfig := s.IdentityConfig(t, context.Background())
	client := gocloak.NewClient(identityConfig.OidcProvider.BaseURL)
	token, err := client.Login(context.Background(), identityConfig.OidcProvider.Frontend.ClientID, identityConfig.OidcProvider.Frontend.ClientSecret, identityConfig.OidcProvider.DomainName, user.Username, "test")
	if err != nil {
		t.Fatalf("loginUser::failed to get token: %v", err)
	}

	return token
}

func (s *KeycloakTestSuite) EnsureUserExists(t testing.TB, user *entities.User) string {
	t.Helper()
	identityConfig := s.IdentityConfig(t, context.Background())
	client := gocloak.NewClient(identityConfig.OidcProvider.BaseURL)
	token, err := client.LoginClient(context.Background(), identityConfig.OidcProvider.Backend.ClientID, identityConfig.OidcProvider.Backend.ClientSecret, identityConfig.OidcProvider.DomainName)
	if err != nil {
		t.Fatalf("ensureUserExists::failed to get token: %v", err)
	}

	attribute := make(map[string][]string)
	attribute["phone_number"] = []string{user.PhoneNumber}
	attribute["employee_id"] = []string{user.EmployeeID}

	kcUser := gocloak.User{
		ID:         gocloak.StringP(user.ID.String()),
		Username:   gocloak.StringP(user.Username),
		FirstName:  gocloak.StringP(user.FirstName),
		LastName:   gocloak.StringP(user.LastName),
		Email:      gocloak.StringP(user.Email),
		Enabled:    gocloak.BoolP(true),
		Attributes: &attribute,
	}

	realm := identityConfig.OidcProvider.DomainName

	userID, err := client.CreateUser(context.Background(), token.AccessToken, realm, kcUser)
	if err != nil {
		t.Log("ensureUserExists::failed to create user. maybe user already exists. error: ", err)
	}

	if err = client.SetPassword(context.Background(), token.AccessToken, userID, realm, "test", false); err != nil {
		t.Fatalf("ensureUserExists::failed to set password: %v", err)
	}

	// set realm role to user
	if len(user.Roles) > 0 {
		ensureUserRolesExists(client, token.AccessToken, realm, userID, user.Roles)
	}

	return userID
}

func ensureUserRolesExists(client *gocloak.GoCloak, accessToken string, realm string, userID string, userRoles []entities.UserRole) {
	var roles []gocloak.Role

	if len(userRoles) > 0 {
		for _, userRole := range userRoles {
			roleName := string(userRole)
			role, err := client.GetRealmRole(context.Background(), accessToken, realm, roleName)

			if err != nil && strings.Contains(err.Error(), "404") {
				log.Printf("ensureUserRolesExists::role %s does not exist. Creating it...", roleName)

				newRole := gocloak.Role{Name: gocloak.StringP(roleName)}
				_, err = client.CreateRealmRole(context.Background(), accessToken, realm, newRole)
				if err != nil {
					log.Fatalf("ensureUserRolesExists::failed to create role %s: %v", roleName, err)
				}

				role, err = client.GetRealmRole(context.Background(), accessToken, realm, roleName)
				if err != nil {
					log.Fatalf("ensureUserRolesExists::failed to retrieve newly created role %s: %v", roleName, err)
				}

				roles = append(roles, *role)
			} else if err != nil {
				log.Fatalf("ensureUserRolesExists::failed to check role %s: %v", roleName, err)
			}
		}

		if len(roles) > 0 {
			err := client.AddRealmRoleToUser(context.Background(), accessToken, realm, userID, roles)
			if err != nil {
				log.Fatalf("ensureUserRolesExists::failed to assign roles to user: %v", err)
			}
		}
	}
}

func (s *KeycloakTestSuite) TestUserToCreateFunc() []*entities.User {
	n := 20
	users := make([]*entities.User, n)
	for i := 0; i < n; i++ {
		users[i] = &entities.User{
			Username:    fmt.Sprintf("test%d", i),
			FirstName:   "Toni",
			LastName:    "Tester",
			Email:       fmt.Sprintf("test%d@green-ecolution.de", i),
			EmployeeID:  "123456",
			PhoneNumber: "+49 123456",
		}
	}
	return users
}
