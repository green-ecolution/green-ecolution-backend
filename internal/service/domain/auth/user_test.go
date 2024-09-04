package auth

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	t.Run("should return result when success", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		inputUser := &entities.User{
			Username:    "username",
			Email:       "mail@foo.com",
			PhoneNumber: "phoneNumber",
			FirstName:   "firstName",
			LastName:    "lastName",
			EmployeeID:  "employeeID",
		}
		input := &entities.RegisterUser{
			User:     *inputUser,
			Password: "password",
			Roles:    &[]string{"viewer"},
		}

		expected := &entities.User{
			ID:            uuid.MustParse("6be4c752-94df-4719-99b1-ce58253eaf75"),
			CreatedAt:     time.Now(),
			Username:      inputUser.Username,
			FirstName:     inputUser.FirstName,
			LastName:      inputUser.LastName,
			Email:         inputUser.Email,
			EmployeeID:    inputUser.EmployeeID,
			PhoneNumber:   inputUser.PhoneNumber,
			EmailVerified: false,
			Avatar:        nil,
		}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		userRepo.EXPECT().Create(context.Background(), inputUser, input.Password, input.Roles).Return(expected, nil)
		resp, err := svc.Register(context.Background(), input)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("should return error when failed to register user", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		resp, err := svc.Register(context.Background(), &entities.RegisterUser{})

		// then
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should return error when validation error", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		resp, err := svc.Register(context.Background(), &entities.RegisterUser{})

		// then
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestLoginRequest(t *testing.T) {
	t.Run("should return login url", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{
			KeyCloak: config.KeyCloakConfig{
				BaseURL:      "http://localhost:8080/auth",
				Realm:        "realm",
				ClientID:     "backend_client",
				ClientSecret: "backend_secret",
				Frontend: config.KeyCloakFrontendConfig{
					AuthURL:      "http://localhost:8080/auth/realms/realm/protocol/openid-connect/auth",
					TokenURL:     "http://localhost:8080/auth/realms/realm/protocol/openid-connect/token",
					ClientID:     "frontend_client",
					ClientSecret: "frontend_secret",
				},
			},
		}

		redirectURL, _ := url.Parse("http://localhost:3000/auth/callback")
		loginRequest := &domain.LoginRequest{
			RedirectURL: redirectURL,
		}

		respURL, _ := url.Parse("http://localhost:8080/auth/realms/realm/protocol/openid-connect/auth")
		query := respURL.Query()
		query.Add("client_id", identityConfig.KeyCloak.Frontend.ClientID)
		query.Add("response_type", "code")
		query.Add("redirect_uri", loginRequest.RedirectURL.String())
		respURL.RawQuery = query.Encode()

		expected := &domain.LoginResp{
			LoginURL: respURL,
		}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		resp, err := svc.LoginRequest(context.Background(), loginRequest)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("should return error when failed to parse auth url in config", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{
			KeyCloak: config.KeyCloakConfig{
				Frontend: config.KeyCloakFrontendConfig{
					AuthURL: "not_a_valid_url",
				},
			},
		}

		loginRequest := &domain.LoginRequest{
			RedirectURL: &url.URL{
				Scheme: "http",
				Host:   "localhost:3000",
				Path:   "/auth/callback",
			},
		}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		_, err := svc.LoginRequest(context.Background(), loginRequest)

		// then
		assert.Error(t, err)
	})
}

func TestClientTokenCallback(t *testing.T) {
	t.Run("should return client token", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		loginCallback := &domain.LoginCallback{
			Code: "code",
			RedirectURL: &url.URL{
				Scheme: "http",
				Host:   "localhost:3000",
				Path:   "/auth/callback",
			},
		}

		expected := &domain.ClientToken{
			AccessToken: "access_token",
		}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		authRepo.EXPECT().GetAccessTokenFromClientCode(context.Background(), loginCallback.Code, loginCallback.RedirectURL.String()).Return(expected, nil)
		resp, err := svc.ClientTokenCallback(context.Background(), loginCallback)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expected, resp)
	})

	t.Run("should return error when validation error", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		loginCallback := &domain.LoginCallback{}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		_, err := svc.ClientTokenCallback(context.Background(), loginCallback)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when failed to get access token", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		loginCallback := &domain.LoginCallback{
			Code: "code",
			RedirectURL: &url.URL{
				Scheme: "http",
				Host:   "localhost:3000",
				Path:   "/auth/callback",
			},
		}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		authRepo.EXPECT().GetAccessTokenFromClientCode(context.Background(), loginCallback.Code, loginCallback.RedirectURL.String()).Return(nil, assert.AnError)
		_, err := svc.ClientTokenCallback(context.Background(), loginCallback)

		// then
		assert.Error(t, err)
	})
}
