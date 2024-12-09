package auth

import (
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			Roles:    []string{"viewer"},
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
		assert.Error(t, err, "500: failed to register user")
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
		assert.Error(t, err, "400: validation error")
		assert.Nil(t, resp)
	})
}

func TestLoginRequest(t *testing.T) {
	t.Run("should return login url", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{
			OidcProvider: config.OidcProvider{
				BaseURL:    "http://localhost:8080/auth",
				DomainName: "realm",
				AuthURL:    "http://localhost:8080/auth/realms/realm/protocol/openid-connect/auth",
				TokenURL:   "http://localhost:8080/auth/realms/realm/protocol/openid-connect/token",
				Backend: config.OidcClient{
					ClientID:     "backend_client",
					ClientSecret: "backend_secret",
				},
				Frontend: config.OidcClient{
					ClientID:     "frontend_client",
					ClientSecret: "frontend_secret",
				},
			},
		}

		redirectURL, _ := url.Parse("http://localhost:3000/auth/callback")
		loginRequest := &entities.LoginRequest{
			RedirectURL: redirectURL,
		}

		respURL, _ := url.Parse("http://localhost:8080/auth/realms/realm/protocol/openid-connect/auth")
		query := respURL.Query()
		query.Add("client_id", identityConfig.OidcProvider.Frontend.ClientID)
		query.Add("response_type", "code")
		query.Add("redirect_uri", loginRequest.RedirectURL.String())
		respURL.RawQuery = query.Encode()

		expected := &entities.LoginResp{
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
			OidcProvider: config.OidcProvider{
				AuthURL: "not_a_valid_url",
			},
		}

		loginRequest := &entities.LoginRequest{
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
		assert.Error(t, err, "500: failed to parse auth url in config")
	})
}

func TestClientTokenCallback(t *testing.T) {
	t.Run("should return client token", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		loginCallback := &entities.LoginCallback{
			Code: "code",
			RedirectURL: &url.URL{
				Scheme: "http",
				Host:   "localhost:3000",
				Path:   "/auth/callback",
			},
		}

		expected := &entities.ClientToken{
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
		loginCallback := &entities.LoginCallback{}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		// when
		_, err := svc.ClientTokenCallback(context.Background(), loginCallback)

		// then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "400: validation error")
	})

	t.Run("should return error when failed to get access token", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}
		loginCallback := &entities.LoginCallback{
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
		assert.Error(t, err, "500: failed to get access token")
	})
}

func TestLogoutRequest(t *testing.T) {
	t.Run("should succeed when logout request is valid", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		logoutRequest := &entities.Logout{RefreshToken: "valid-refresh-token"}
		userRepo.EXPECT().RemoveSession(mock.Anything, logoutRequest.RefreshToken).Return(nil)

		// then
		err := svc.LogoutRequest(context.Background(), logoutRequest)

		// when
		assert.NoError(t, err)
	})

	t.Run("should return error when validation fails", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		invalidLogoutRequest := &entities.Logout{RefreshToken: ""}

		// when
		err := svc.LogoutRequest(context.Background(), invalidLogoutRequest)

		// then
		assert.Error(t, err, "400: validation error")
	})

	t.Run("should return error when session removal fails", func(t *testing.T) {
		// given
		identityConfig := &config.IdentityAuthConfig{}

		authRepo := storageMock.NewMockAuthRepository(t)
		userRepo := storageMock.NewMockUserRepository(t)
		svc := NewAuthService(authRepo, userRepo, identityConfig)

		logoutRequest := &entities.Logout{RefreshToken: "valid-refresh-token"}
		userRepo.EXPECT().RemoveSession(mock.Anything, logoutRequest.RefreshToken).Return(errors.New("Internal error"))

		// when
		err := svc.LogoutRequest(context.Background(), logoutRequest)

		// then
		assert.Error(t, err, "500: failed to remove session")
	})
}
