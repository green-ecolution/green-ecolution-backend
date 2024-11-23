package tree

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
)

func TestRegisterTreeRoutes(t *testing.T) {
	t.Run("/v1/tree", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockTreeService := serviceMock.NewMockTreeService(t)
			app := RegisterRoutes(mockTreeService)

			mockTreeService.EXPECT().GetAll(
				mock.Anything,
			).Return(TestTrees, nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("should call POST", func(t *testing.T) {
			mockTreeService := serviceMock.NewMockTreeService(t)
			app := RegisterRoutes(mockTreeService)

			mockTreeService.EXPECT().Create(
				mock.Anything,
				mock.AnythingOfType("*entities.TreeCreate"),
			).Return(TestTrees[0], nil)

			// when
			body, _ := json.Marshal(TestTreeCreateRequest)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	})

	t.Run("/v1/tree/:id", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockTreeService := serviceMock.NewMockTreeService(t)
			app := RegisterRoutes(mockTreeService)

			mockTreeService.EXPECT().GetByID(
				mock.Anything,
				int32(1),
			).Return(TestTrees[0], nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("should call PUT handler", func(t *testing.T) {
			mockTreeService := serviceMock.NewMockTreeService(t)
			app := RegisterRoutes(mockTreeService)

			mockTreeService.EXPECT().Update(
				mock.Anything,
				int32(1),
				mock.Anything,
			).Return(TestTrees[0], nil)

			// when
			body, _ := json.Marshal(TestTreeUpdateRequest)
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodPut, "/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("should call DELETE", func(t *testing.T) {
			mockTreeService := serviceMock.NewMockTreeService(t)
			app := RegisterRoutes(mockTreeService)

			mockTreeService.EXPECT().Delete(
				mock.Anything,
				int32(1),
			).Return(nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		})
	})

	t.Run("/v1/tree/sensor/:sensor_id", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			mockTreeService := serviceMock.NewMockTreeService(t)
			app := RegisterRoutes(mockTreeService)

			mockTreeService.EXPECT().GetBySensorID(
				mock.Anything,
				int32(1),
			).Return(TestTrees[0], nil)

			// when
			req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/sensor/1", nil)

			// then
			resp, err := app.Test(req)
			defer resp.Body.Close()
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})

	t.Run("/v1/tree/:id/images", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			//TODO: Implement test for GetTreeImages
		})

		t.Run("should call POST handler", func(t *testing.T) {
			//TODO: Implement test for AddTreeImage
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			//TODO: Implement test for RemoveTreeImage
		})
	})

	t.Run("/v1/tree/:id/sensor", func(t *testing.T) {
		t.Run("should call GET handler", func(t *testing.T) {
			//TODO: Implement test for GetTreeSensor
		})

		t.Run("should call POST handler", func(t *testing.T) {
			//TODO: Implement test for AddTreeSensor
		})

		t.Run("should call DELETE handler", func(t *testing.T) {
			//TODO: Implement test for RemoveTreeSensor
		})
	})
}
