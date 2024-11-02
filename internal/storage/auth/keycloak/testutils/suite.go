package testutils

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	keycloak "github.com/stillya/testcontainers-keycloak"
)

type KeycloakTestSuite struct {
	Container        *keycloak.KeycloakContainer
	RealmName        string
	BackendClientId  string
	FrontendClientId string
	User             string
	Password         string
}

var (
	User     = "admin"
	Password = "admin"
)

func SetupKeycloakTestSuite(ctx context.Context) *KeycloakTestSuite {
	container := startKeycloakContainer(ctx)
	return &KeycloakTestSuite{
		Container:        container,
		User:             User,
		Password:         Password,
		RealmName:        "green-ecolution-test",
		BackendClientId:  "green-ecolution-backend",
		FrontendClientId: "green-ecolution-frontend",
	}
}

func startKeycloakContainer(ctx context.Context) *keycloak.KeycloakContainer {
	realmImportFile := fmt.Sprintf("%s/internal/storage/auth/keycloak/testutils/testdata/realm-export.json", utils.RootDir())
	keycloakContainer, err := keycloak.Run(ctx,
		"keycloak/keycloak:25.0",
		keycloak.WithContextPath("/auth"),
		keycloak.WithAdminUsername(User),
		keycloak.WithAdminPassword(Password),
		keycloak.WithRealmImportFile(realmImportFile),
	)
	if err != nil {
		log.Fatalf("Could not start keycloak container: %s", err)
	}

	return keycloakContainer
}

func (s *KeycloakTestSuite) IdentityConfig(t testing.TB, ctx context.Context) *config.IdentityAuthConfig {
	t.Helper()
	backendClient := s.GetBackendClient(t, ctx)
	frontendClient := s.GetFrontendClient(t, ctx)

	return &config.IdentityAuthConfig{
		KeyCloak: config.KeyCloakConfig{
			BaseURL:        s.GetAuthServerURL(t, ctx),
			Realm:          s.RealmName,
			ClientID:       s.BackendClientId,
			ClientSecret:   *backendClient.Secret,
			RealmPublicKey: "",
			Frontend: config.KeyCloakFrontendConfig{
				ClientID:     s.FrontendClientId,
				ClientSecret: *frontendClient.Secret,
				AuthURL:      s.GetAuthServerURL(t, ctx) + "/realms/" + s.RealmName + "/protocol/openid-connect/auth",
				TokenURL:     s.GetAuthServerURL(t, ctx) + "/realms/" + s.RealmName + "/protocol/openid-connect/token",
			},
		},
	}
}

func (s *KeycloakTestSuite) GetAdminClient(t testing.TB, ctx context.Context) *keycloak.AdminClient {
	t.Helper()
	adminClient, err := s.Container.GetAdminClient(ctx)
	if err != nil {
		t.Fatalf("Could not get admin client: %s", err)
	}

	return adminClient
}

func (s *KeycloakTestSuite) GetAuthServerURL(t testing.TB, ctx context.Context) string {
	t.Helper()
	url, err := s.Container.GetAuthServerURL(ctx)
	if err != nil {
		t.Fatalf("Could not get auth server URL: %s", err)
	}
	return url
}

func (s *KeycloakTestSuite) GetBackendClient(t testing.TB, ctx context.Context) *keycloak.Client {
	t.Helper()
	adminClient, err := s.Container.GetAdminClient(ctx)
	if err != nil {
		t.Fatalf("Could not get admin client: %s", err)
	}

	backendClient, err := adminClient.GetClient(s.RealmName, s.BackendClientId)
	if err != nil {
		t.Fatalf("Could not get backend client: %s", err)
	}

	return backendClient
}

func (s *KeycloakTestSuite) GetFrontendClient(t testing.TB, ctx context.Context) *keycloak.Client {
	t.Helper()
	adminClient, err := s.Container.GetAdminClient(ctx)
	if err != nil {
		t.Fatalf("Could not get admin client: %s", err)
	}

	frontendClient, err := adminClient.GetClient(s.RealmName, s.FrontendClientId)
	if err != nil {
		t.Fatalf("Could not get frontend client: %s", err)
	}

	return frontendClient
}

func (s *KeycloakTestSuite) Terminate(ctx context.Context) {
	if err := s.Container.Terminate(ctx); err != nil {
		log.Fatalf("Could not terminate container: %s", err)
	}
}
