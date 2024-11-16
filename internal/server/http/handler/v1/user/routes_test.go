package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("v1/user", func(t *testing.T) {
		t.Log("v1/user should call GET handler (not implemented)") // TODO: implement testcase

		t.Run("v1/user should call POST", func(t *testing.T) {
			mockUserService := serviceMock.NewMockAuthService(t)
			app := RegisterRoutes(mockUserService)
			expected := &domain.User{
				ID:            uuid.MustParse("6be4c752-94df-4719-99b1-ce58253eaf75"),
				CreatedAt:     time.Now(),
				Username:      "toni_tester",
				FirstName:     "Toni",
				LastName:      "Tester",
				Email:         "dev@green-ecolution.de",
				EmployeeID:    "123456",
				PhoneNumber:   "+49 123456",
				EmailVerified: false,
				Avatar:        nil,
			}

			mockUserService.EXPECT().Register(
				mock.Anything,
				mock.AnythingOfType("*entities.RegisterUser"),
			).Return(expected, nil)

			// when
			body, _ := json.Marshal(entities.UserRegisterRequest{
				Username:    "toni_tester",
				FirstName:   "Toni",
				LastName:    "Tester",
				Email:       "dev@green-ecolution.de",
				EmployeeID:  "123456",
				PhoneNumber: "+49 123456",
				Password:    "test",
			})
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			var got entities.UserResponse
			err = json.NewDecoder(resp.Body).Decode(&got)
			assert.NoError(t, err)
			assert.Equal(t, expected.ID.String(), got.ID)
			assert.Equal(t, expected.Username, got.Username)
			assert.Equal(t, expected.FirstName, got.FirstName)
			assert.Equal(t, expected.LastName, got.LastName)
			assert.Equal(t, expected.Email, got.Email)
			assert.Equal(t, expected.EmployeeID, got.EmployeeID)
			assert.Equal(t, expected.PhoneNumber, got.PhoneNumber)
		})
	})

	t.Run("v1/user/:id", func(t *testing.T) {
		t.Log("should call GET handler (not implemented)")    // TODO: implement testcase
		t.Log("should call PUT handler (not implemented)")    // TODO: implement testcase
		t.Log("should call DELETE handler (not implemented)") // TODO: implement testcase
	})

	t.Run("v1/user/:id/roles", func(t *testing.T) {
		t.Log("should call GET handler (not implemented)") // TODO: implement testcase
	})
}

func TestRegisterPublicRoutes(t *testing.T) {
	t.Run("v1/user/logout should call POST handler", func(t *testing.T) {
		mockUserService := serviceMock.NewMockAuthService(t)
		app := RegisterPublicRoutes(mockUserService)

		mockUserService.EXPECT().LogoutRequest(
			mock.Anything,
			mock.AnythingOfType("*entities.Logout"),
		).Return(nil)

		// when
		body, _ := json.Marshal(entities.LogoutRequest{
			RefreshToken: "refresh-token",
		})
		req := httptest.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// then
		resp, err := app.Test(req)
		defer resp.Body.Close()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("v1/user/login should call GET handler", func(t *testing.T) {
		// given
		mockUserService := serviceMock.NewMockAuthService(t)
		app := RegisterPublicRoutes(mockUserService)
		loginURL, _ := url.Parse("http://localhost:8080/auth/realms/green-ecolution/protocol/openid-connect/auth?client_id=green-ecolution-frontend&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Flogin&response_type=code&scope=openid%20profile%20email&state=state&nonce=nonce")
		expected := &domain.LoginResp{
			LoginURL: loginURL,
		}

		mockUserService.EXPECT().LoginRequest(
			mock.Anything,
			mock.AnythingOfType("*entities.LoginRequest"),
		).Return(expected, nil)

		// when
		req := httptest.NewRequest(http.MethodGet, "/login?redirect_url=http%3A%2F%2Flocalhost%3A3000%2Flogin", nil)
		req.Header.Set("Content-Type", "application/json")

		// then
		resp, err := app.Test(req)
		defer resp.Body.Close()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var got entities.LoginResponse
		err = json.NewDecoder(resp.Body).Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, expected.LoginURL.String(), got.LoginURL)
	})

	t.Run("v1/user/login/token should call POST handler", func(t *testing.T) {
		mockUserService := serviceMock.NewMockAuthService(t)
		app := RegisterPublicRoutes(mockUserService)
		expected := &domain.ClientToken{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600,
			TokenType:    "Bearer",
		}

		mockUserService.EXPECT().ClientTokenCallback(
			mock.Anything,
			mock.AnythingOfType("*entities.LoginCallback"),
		).Return(expected, nil)

		// when
		body, _ := json.Marshal(entities.LoginTokenRequest{
			Code: "code",
		})
		req := httptest.NewRequest(http.MethodPost, "/login/token?redirect_url=http%3A%2F%2Flocalhost%3A3000%2Flogin", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// then
		resp, err := app.Test(req)
		defer resp.Body.Close()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var got entities.ClientTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, expected.AccessToken, got.AccessToken)
		assert.Equal(t, expected.RefreshToken, got.RefreshToken)
		assert.Equal(t, expected.ExpiresIn, got.ExpiresIn)
		assert.Equal(t, expected.TokenType, got.TokenType)
	})

	t.Run("v1/user/token/refresh should call POST handler", func(t *testing.T) {
		mockUserService := serviceMock.NewMockAuthService(t)
		app := RegisterPublicRoutes(mockUserService)
		refreshToken := generateJWT(t, "user123")

		expectedResponse := &domain.ClientToken{
			AccessToken:  "valid_access_token",
			RefreshToken: refreshToken,
			ExpiresIn:    3600,
			TokenType:    "Bearer",
		}
		mockUserService.EXPECT().RefreshToken(mock.Anything, refreshToken).Return(expectedResponse, nil)

		// when
		body, _ := json.Marshal(entities.RefreshTokenRequest{
			RefreshToken: refreshToken,
		})
		req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// then
		resp, err := app.Test(req)
		defer resp.Body.Close()
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var got entities.ClientTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&got)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.AccessToken, got.AccessToken)
		assert.Equal(t, expectedResponse.RefreshToken, got.RefreshToken)
		assert.Equal(t, expectedResponse.ExpiresIn, got.ExpiresIn)
		assert.Equal(t, expectedResponse.TokenType, got.TokenType)
	})

}
