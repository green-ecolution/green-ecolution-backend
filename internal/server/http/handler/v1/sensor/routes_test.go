package sensor

import (
	"context"
	"net/http"
	"testing"

	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("/v1/sensor", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockSensorService := serviceMock.NewMockSensorService(t)
			app := RegisterRoutes(mockSensorService)

			mockSensorService.EXPECT().GetAll(
				mock.Anything,
			).Return(TestSensorList, nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})
}
