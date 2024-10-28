package treecluster

import (
	"context"
	"net/http"
	"testing"

	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("/v1/region", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockClusterService := serviceMock.NewMockTreeClusterService(t)
			app := RegisterRoutes(mockClusterService)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)

			mockClusterService.EXPECT().GetAll(
				mock.Anything,
			).Return(TestClusterList, nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})
}