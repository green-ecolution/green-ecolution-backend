package user

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	t.Run("Should login user sucessfully (200).", func(t *testing.T) {
		app := fiber.New()
		mockAuthService := serviceMock.NewMockAuthService(t)
		app.Get("/v1/user/login", Login(mockAuthService))

		parsedUrlRedirect, _ := url.Parse("http://example.com/redirect")
		parsedUrlResponse, _ := url.Parse("http://example.com/login")

		loginRequest := &entities.LoginRequest{
			RedirectURL: parsedUrlRedirect,
		}
		loginResponse := &entities.LoginResp{
			LoginURL: parsedUrlResponse,
		}
		mockAuthService.EXPECT().LoginRequest(mock.Anything, loginRequest).Return(loginResponse, nil)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/v1/user/login?redirect_url="+parsedUrlRedirect.String(), nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Should return 400 bad request for invalid request body.", func(t *testing.T) {

	})

	t.Run("Should return 500 for invalid login.", func(t *testing.T) {

	})

}
