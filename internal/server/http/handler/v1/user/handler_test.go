package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	t.Run("should login user sucessfully", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Get("/v1/user/login", Login(mockAuthService))

		parsedUrlRedirect, _ := url.Parse("http://example.com/redirect")
		parsedUrlResponse, _ := url.Parse("http://example.com/login")

		loginRequest := &domain.LoginRequest{
			RedirectURL: parsedUrlRedirect,
		}
		loginResponse := &domain.LoginResp{
			LoginURL: parsedUrlResponse,
		}
		mockAuthService.EXPECT().LoginRequest(mock.Anything, loginRequest).Return(loginResponse, nil)

		// when
		req := httptest.NewRequest(http.MethodGet, "/v1/user/login?redirect_url="+parsedUrlRedirect.String(), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Should return 400 bad request for invalid request body.", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Get("/v1/user/login", Login(mockAuthService))

		// when
		req := httptest.NewRequest(http.MethodGet, "/v1/user/login", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Should return 500 for invalid login.", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Get("/v1/user/login", Login(mockAuthService))

		parsedUrlRedirect, _ := url.Parse("http://example.com/redirect")

		loginRequest := &domain.LoginRequest{
			RedirectURL: parsedUrlRedirect,
		}
		mockAuthService.EXPECT().LoginRequest(mock.Anything, loginRequest).Return(nil, errors.New("service error"))

		// when
		req := httptest.NewRequest(http.MethodGet, "/v1/user/login?redirect_url="+parsedUrlRedirect.String(), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Should return 400 for invalid redirect url.", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Get("/v1/user/login", Login(mockAuthService))

		// when
		req := httptest.NewRequest(http.MethodGet, "/v1/user/login?redirect_url=invalid-url", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})
}

func TestLogout(t *testing.T) {
	t.Run("should logout user sucessfully", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Post("/v1/user/logout", Logout(mockAuthService))

		domainEntity := domain.Logout{
			RefreshToken: "valid_refresh_token",
		}

		// when
		mockAuthService.EXPECT().LogoutRequest(mock.Anything, &domainEntity).Return(nil)
		reqBody, _ := json.Marshal(entities.LogoutRequest{
			RefreshToken: "valid_refresh_token",
		})
		req := httptest.NewRequest(http.MethodPost, "/v1/user/logout", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Should return 400 bad request for invalid request body.", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Post("/v1/user/logout", Logout(mockAuthService))

		// when
		req := httptest.NewRequest(http.MethodPost, "/v1/user/logout", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Should return 500 for invalid logout.", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Post("/v1/user/logout", Logout(mockAuthService))

		domainEntity := domain.Logout{
			RefreshToken: "valid_refresh_token",
		}

		// when
		mockAuthService.EXPECT().LogoutRequest(mock.Anything, &domainEntity).Return(errors.New("service error"))
		reqBody, _ := json.Marshal(entities.LogoutRequest{
			RefreshToken: "valid_refresh_token",
		})
		req := httptest.NewRequest(http.MethodPost, "/v1/user/logout", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})
}

func TestRequestToken(t *testing.T) {
	t.Run("should request token sucessfully", func(t *testing.T) {
		// given
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Post("/v1/user/token", RequestToken(mockAuthService))

		redirectURL, _ := url.Parse("http://example.com/redirect")
		domainEntity := domain.LoginCallback{
			Code:        "valid_code",
			RedirectURL: redirectURL,
		}

		expectedResponse := &domain.ClientToken{
			AccessToken:  "valid_access_token",
			RefreshToken: "valid_refresh_token",
			ExpiresIn:    3600,
			TokenType:    "Bearer",
		}

		// when
		mockAuthService.EXPECT().ClientTokenCallback(mock.Anything, &domainEntity).Return(expectedResponse, nil)
		reqBody, _ := json.Marshal(entities.LoginTokenRequest{
			Code: "valid_code",
		})
		req := httptest.NewRequest(http.MethodPost, "/v1/user/token?redirect_url="+redirectURL.String(), bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)

		var response entities.ClientTokenResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse.AccessToken, response.AccessToken)
		assert.Equal(t, expectedResponse.RefreshToken, response.RefreshToken)
		assert.Equal(t, expectedResponse.ExpiresIn, response.ExpiresIn)
		assert.Equal(t, expectedResponse.TokenType, response.TokenType)
	})

  t.Run("Should return 400 bad request for invalid request body.", func(t *testing.T) {
    // given
    app := fiber.New()
    mockAuthService := serviceMock.NewMockAuthService(t)
    app.Post("/v1/user/token", RequestToken(mockAuthService))

    // when
    req := httptest.NewRequest(http.MethodPost, "/v1/user/token", nil)
    resp, err := app.Test(req, -1)
    defer resp.Body.Close()

    // then
    assert.Nil(t, err)
    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    mockAuthService.AssertExpectations(t)
  })

  t.Run("Should return 500 for invalid request token.", func(t *testing.T) {
    // given
    app := fiber.New()
    mockAuthService := serviceMock.NewMockAuthService(t)
    app.Post("/v1/user/token", RequestToken(mockAuthService))

    redirectURL, _ := url.Parse("http://example.com/redirect")
    domainEntity := domain.LoginCallback{
      Code:        "valid_code",
      RedirectURL: redirectURL,
    }

    // when
    mockAuthService.EXPECT().ClientTokenCallback(mock.Anything, &domainEntity).Return(nil, errors.New("service error"))
    reqBody, _ := json.Marshal(entities.LoginTokenRequest{
      Code: "valid_code",
    })
    req := httptest.NewRequest(http.MethodPost, "/v1/user/token?redirect_url="+redirectURL.String(), bytes.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(req, -1)
    defer resp.Body.Close()

    // then
    assert.Nil(t, err)
    assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
    mockAuthService.AssertExpectations(t)
  })

  t.Run("Should return 400 for invalid redirect url.", func(t *testing.T) {
    // given
    app := fiber.New()
    mockAuthService := serviceMock.NewMockAuthService(t)
    app.Post("/v1/user/token", RequestToken(mockAuthService))

    // when
    req := httptest.NewRequest(http.MethodPost, "/v1/user/token?redirect_url=invalid-url", nil)
    resp, err := app.Test(req, -1)
    defer resp.Body.Close()

    // then
    assert.Nil(t, err)
    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    mockAuthService.AssertExpectations(t)
  })

  t.Run("Should return 400 for invalid code.", func(t *testing.T) {
    // given
    app := fiber.New()
    mockAuthService := serviceMock.NewMockAuthService(t)
    app.Post("/v1/user/token", RequestToken(mockAuthService))

    // when
    req := httptest.NewRequest(http.MethodPost, "/v1/user/token", nil)
    resp, err := app.Test(req, -1)
    defer resp.Body.Close()

    // then
    assert.Nil(t, err)
    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    mockAuthService.AssertExpectations(t)
  })
}
  
