package info_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/info"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAppInfo(t *testing.T) {
	t.Run("should return app info successfully", func(t *testing.T) {
		mockInfoService := serviceMock.NewMockInfoService(t)
		app := fiber.New()
		handler := info.GetAppInfo(mockInfoService)

		repoURL, _ := url.Parse("https://github.com/green-ecolution/green-ecolution-backend")
		serverURL, _ := url.Parse("http://localhost")

		expectedInfo := &entities.App{
			Version:   "1.0.0",
			GoVersion: "go1.23.2",
			BuildTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			Git: entities.Git{
				Branch:     "main",
				Commit:     "abcd1234",
				Repository: repoURL,
			},
			Server: entities.Server{
				OS:        "linux",
				Arch:      "amd64",
				Hostname:  "localhost",
				URL:       serverURL,
				IP:        net.ParseIP("127.0.0.1"),
				Port:      8080,
				Interface: "eth0",
				Uptime:    24 * time.Hour,
			},
		}

		mockInfoService.EXPECT().GetAppInfoResponse(mock.Anything).Return(expectedInfo, nil)
		app.Get("/v1/info", handler)

		// Act
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/info", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response serverEntities.AppInfoResponse
		err = utils.ParseJSONResponse(resp, &response)
		assert.NoError(t, err)

		mockInfoService.AssertExpectations(t)
	})

	t.Run("should return 500 when service returns an error", func(t *testing.T) {
		mockInfoService := serviceMock.NewMockInfoService(t)
		app := fiber.New()
		handler := info.GetAppInfo(mockInfoService)

		mockInfoService.EXPECT().GetAppInfoResponse(mock.Anything).Return(nil, errors.New("service error"))
		app.Get("/v1/info", handler)

		// when
		req, _ := http.NewRequestWithContext(context.Background(), "GET", "/v1/info", nil)
		resp, err := app.Test(req, -1)
		defer resp.Body.Close()

		// then
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockInfoService.AssertExpectations(t)
	})
}
