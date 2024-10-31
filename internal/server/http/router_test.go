package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPublicRoutes(t *testing.T) {
    t.Run("/", func(t *testing.T) {
        t.Run("should return hello world", func(t *testing.T) {
            app := fiber.New()

            // given
            mockServices := setupMockServices(t)

            server := &Server{services: mockServices}
            server.publicRoutes(app)

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
        t.Run("/user/logout", func(t *testing.T) {
            app := fiber.New()

            mockServices := setupMockServices(t)
            server := &Server{services: mockServices}
            server.publicRoutes(app)

            refreshToken := "some-refresh-token"

            mockServices.AuthService.(*serviceMock.MockAuthService).EXPECT().
				LogoutRequest(
					mock.Anything,
					&domain.Logout{RefreshToken: refreshToken},
				).Return(nil)

            requestBody := fmt.Sprintf(`{"refresh_token": "%s"}`, refreshToken)

            // when
            req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/user/logout", bytes.NewBuffer([]byte(requestBody)))
            req.Header.Set("Content-Type", "application/json")
            assert.NoError(t, err)

            resp, err := app.Test(req)
            assert.NoError(t, err)
            defer resp.Body.Close()

            // then
            assert.Equal(t, http.StatusOK, resp.StatusCode)
        })

        t.Run("/user/login", func(t *testing.T) {
            app := fiber.New()

            mockServices := setupMockServices(t)
            server := &Server{services: mockServices}
            server.publicRoutes(app)

            redirectURL, err := url.Parse("http://localhost:3000/auth/callback")
            assert.NoError(t, err)

            loginRequest := &domain.LoginRequest{
                RedirectURL: redirectURL,
            }

            respURL, err := url.Parse("http://localhost:8080/auth/realms/realm/protocol/openid-connect/auth")
            assert.NoError(t, err)

            expected := &domain.LoginResp{
                LoginURL: respURL,
            }

			mockServices.AuthService.(*serviceMock.MockAuthService).EXPECT().
				LoginRequest(mock.Anything, mock.MatchedBy(func(req *domain.LoginRequest) bool {
					return req.RedirectURL.String() == loginRequest.RedirectURL.String()
				})).
				Return(expected, nil)

			// when
            req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/user/login?redirect_url=http://localhost:3000/auth/callback", nil)
            assert.NoError(t, err)

            resp, err := app.Test(req)
            assert.NoError(t, err)
            defer resp.Body.Close()

			// then
            assert.Equal(t, http.StatusOK, resp.StatusCode)

            body, err := io.ReadAll(resp.Body)
            assert.NoError(t, err)

            var loginResponse entities.LoginResponse
            err = json.Unmarshal(body, &loginResponse)
            assert.NoError(t, err)

            assert.Equal(t, expected.LoginURL.String(), loginResponse.LoginURL)
        })
    })
}

func setupMockServices(t *testing.T) *service.Services {
    mockInfoService := serviceMock.NewMockInfoService(t)
    mockClusterService := serviceMock.NewMockTreeClusterService(t)
    mockTreeService := serviceMock.NewMockTreeService(t)
    mockSensorService := serviceMock.NewMockSensorService(t)
    mockAuthService := serviceMock.NewMockAuthService(t)
    mockRegionService := serviceMock.NewMockRegionService(t)

    return &service.Services{
        InfoService:        mockInfoService,
        TreeClusterService: mockClusterService,
        TreeService:       mockTreeService,
        SensorService:     mockSensorService,
        AuthService:       mockAuthService,
        RegionService:     mockRegionService,
    }
}