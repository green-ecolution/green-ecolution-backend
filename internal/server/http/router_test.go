package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServices struct {
	InfoService        *serviceMock.MockInfoService
	TreeClusterService *serviceMock.MockTreeClusterService
	TreeService        *serviceMock.MockTreeService
	SensorService      *serviceMock.MockSensorService
	AuthService        *serviceMock.MockAuthService
	RegionService      *serviceMock.MockRegionService
}

func TestPublicRoutes(t *testing.T) {
	app := fiber.New()

	// given
	mockInfoService := serviceMock.NewMockInfoService(t)
	mockClusterService := serviceMock.NewMockTreeClusterService(t)
	mockTreeService := serviceMock.NewMockTreeService(t)
	mockSensorService := serviceMock.NewMockSensorService(t)
	mockAuthService := serviceMock.NewMockAuthService(t)
	mockRegionService := serviceMock.NewMockRegionService(t)

	mockServices := &service.Services{
		InfoService:        mockInfoService,
		TreeClusterService: mockClusterService,
		TreeService:       mockTreeService,
		SensorService:     mockSensorService,
		AuthService:       mockAuthService,
		RegionService:     mockRegionService,
	}

	server := &Server{services: mockServices}

	server.publicRoutes(app)

	t.Run("/", func(t *testing.T) {
		t.Run("should return hello world", func(t *testing.T) {
			// when
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			assert.NoError(t, err)
			resp, err := app.Test(req)
			defer resp.Body.Close()

			// then
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "Hello, World!", string(body))
		})
	})

	t.Run("/api/v1", func(t *testing.T) {
		t.Run("POST /user/logout", func(t *testing.T) {
			refreshToken := "some-refresh-token"

			mockAuthService.EXPECT().LogoutRequest(
				mock.Anything, 
				&domain.Logout{RefreshToken: refreshToken},
			).Return(nil)
		
			requestBody := fmt.Sprintf(`{"refresh_token": "%s"}`, refreshToken)
		
			// when
			req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/user/logout", bytes.NewBuffer([]byte(requestBody)))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
		
			resp, err := app.Test(req)
			defer resp.Body.Close()
				
			// then
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			mockAuthService.AssertExpectations(t)
		})
	})
}
