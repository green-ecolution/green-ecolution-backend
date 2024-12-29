package service

import (
	"testing"
	"time"

	serviceMock "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	"github.com/stretchr/testify/assert"
)

func TestAllServiceReady(t *testing.T) {
	t.Run("should return true if all service implemented the ServiceReady interface", func(t *testing.T) {
		// given
		infoSvc := serviceMock.NewMockInfoService(t)
		treeSvc := serviceMock.NewMockTreeService(t)
		authSvc := serviceMock.NewMockAuthService(t)
		regionSvc := serviceMock.NewMockRegionService(t)
		treeClusterSvc := serviceMock.NewMockTreeClusterService(t)
		sensorSvc := serviceMock.NewMockSensorService(t)
		vehicleSvc := serviceMock.NewMockVehicleService(t)
		pluginSvc := serviceMock.NewMockPluginService(t)
		wateringPlanSvc := serviceMock.NewMockWateringPlanService(t)
		svc := Services{
			InfoService:         infoSvc,
			TreeService:         treeSvc,
			AuthService:         authSvc,
			RegionService:       regionSvc,
			TreeClusterService:  treeClusterSvc,
			SensorService:       sensorSvc,
			VehicleService:      vehicleSvc,
			PluginService:       pluginSvc,
			WateringPlanService: wateringPlanSvc,
		}

		// when
		infoSvc.EXPECT().Ready().Return(true)
		treeSvc.EXPECT().Ready().Return(true)
		authSvc.EXPECT().Ready().Return(true)
		regionSvc.EXPECT().Ready().Return(true)
		treeClusterSvc.EXPECT().Ready().Return(true)
		sensorSvc.EXPECT().Ready().Return(true)
		vehicleSvc.EXPECT().Ready().Return(true)
		pluginSvc.EXPECT().Ready().Return(true)
		wateringPlanSvc.EXPECT().Ready().Return(true)

		ready := svc.AllServicesReady()

		// then
		assert.True(t, ready)
	})
}

func TestNewError(t *testing.T) {
	t.Run("should create an error with correct fields", func(t *testing.T) {
		// when
		code := NotFound
		message := "entity not found"

		err := NewError(code, message)

		// validate
		assert.NotNil(t, err)
		assert.Equal(t, code, err.Code)
		assert.Equal(t, message, err.Message)
		assert.NotEmpty(t, err.File)
		assert.NotZero(t, err.Line)
		assert.NotEmpty(t, err.Timestamp)

		_, parseErr := time.Parse(time.RFC3339, err.Timestamp)
		assert.NoError(t, parseErr, "timestamp should be in RFC3339 format")
	})
	t.Run("should capture the correct caller file and line", func(t *testing.T) {
		// when
		err := NewError(BadRequest, "validation failed")

		// validate
		assert.Contains(t, err.File, "service_test.go", "error should capture the correct file")
		assert.True(t, err.Line > 0, "line number should be greater than zero")
	})
	t.Run("should format the error string correctly", func(t *testing.T) {
		// when
		code := InternalError
		message := "unexpected internal error"

		err := NewError(code, message)
		formattedError := err.Error()

		// validate
		assert.Contains(t, formattedError, "[500]")
		assert.Contains(t, formattedError, message)
		assert.Contains(t, formattedError, err.File)
		assert.Contains(t, formattedError, err.Timestamp)
	})

}
