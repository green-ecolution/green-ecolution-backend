package vehicle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVehicleService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("should return all vehicles when successful", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		expectedVehicles := getTestVehicles()
		vehicleRepo.EXPECT().GetAll(ctx).Return(expectedVehicles, nil)

		// when
		vehicles, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicles, vehicles)
	})

	t.Run("should return empty slice when no vehicles are found", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		vehicleRepo.EXPECT().GetAll(ctx).Return([]*entities.Vehicle{}, nil)

		// when
		vehicles, err := svc.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Empty(t, vehicles)
	})

	t.Run("should return error when GetAll fails", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		expectedErr := errors.New("GetAll failed")
		vehicleRepo.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		// when
		vehicles, err := svc.GetAll(ctx)

		// then
		assert.Error(t, err)
		assert.Nil(t, vehicles)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestVehicleService_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("should return vehicle when found", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		id := int32(1)
		expectedVehicle := getTestVehicles()[0]
		vehicleRepo.EXPECT().GetByID(ctx, id).Return(expectedVehicle, nil)

		// when
		vehicle, err := svc.GetByID(ctx, id)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle, vehicle)
	})

	t.Run("should return error if vehicle not found", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		id := int32(1)
		expectedErr := storage.ErrEntityNotFound
		vehicleRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		vehicle, err := svc.GetByID(ctx, id)

		// then
		assert.Error(t, err)
		assert.Nil(t, vehicle)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestVehicleService_GetByPlate(t *testing.T) {
	ctx := context.Background()

	t.Run("should return vehicle when found", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		plate := "FL TBZ 1234"
		expectedVehicle := getTestVehicles()[0]
		vehicleRepo.EXPECT().GetByPlate(ctx, plate).Return(expectedVehicle, nil)

		// when
		vehicle, err := svc.GetByPlate(ctx, plate)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle, vehicle)
	})

	t.Run("should return error if vehicle not found", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		plate := "FL TBZ 1234"
		expectedErr := storage.ErrEntityNotFound
		vehicleRepo.EXPECT().GetByPlate(ctx, plate).Return(nil, expectedErr)

		// when
		vehicle, err := svc.GetByPlate(ctx, plate)

		// then
		assert.Error(t, err)
		assert.Nil(t, vehicle)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestVehicleService_Create(t *testing.T) {
	ctx := context.Background()
	input := &entities.VehicleCreate{
		NumberPlate:   "FL TBZ 123",
		Description:   "Test description",
		Status:        entities.VehicleStatusActive,
		Type:          entities.VehicleTypeTrailer,
		WaterCapacity: 2000.5,
	}

	t.Run("should successfully create a new vehicle", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		expectedVehicle := getTestVehicles()[0]

		vehicleRepo.EXPECT().Create(
			ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(expectedVehicle, nil)

		// when
		result, err := svc.Create(ctx, input)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle, result)
	})

	t.Run("should return an error when creating vehicle fails", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		expectedErr := errors.New("Failed to create cluster")

		vehicleRepo.EXPECT().Create(
			ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedErr)

		// when
		result, err := svc.Create(ctx, input)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestVehicleService_Update(t *testing.T) {
	ctx := context.Background()
	vehicleID := int32(1)
	input := &entities.VehicleUpdate{
		NumberPlate:   "FL TBZ 123",
		Description:   "Test description",
		Status:        entities.VehicleStatusActive,
		Type:          entities.VehicleTypeTrailer,
		WaterCapacity: 2000.5,
	}

	t.Run("should successfully update a vehicle", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		expectedVehicle := getTestVehicles()[0]

		vehicleRepo.EXPECT().GetByID(
			ctx,
			vehicleID,
		).Return(expectedVehicle, nil)

		vehicleRepo.EXPECT().Update(
			ctx,
			vehicleID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(expectedVehicle, nil)

		// when
		result, err := svc.Update(ctx, vehicleID, input)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedVehicle, result)
	})

	t.Run("should return an error when vehicle is not found", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		vehicleRepo.EXPECT().GetByID(
			ctx,
			vehicleID,
		).Return(nil, storage.ErrVehicleNotFound)

		// when
		result, err := svc.Update(ctx, vehicleID, input)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(storage.ErrVehicleNotFound).Error())
	})

	t.Run("should return an error when the update fails", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		expectedErr := errors.New("failed to update cluster")
		expectedVehicle := getTestVehicles()[0]

		vehicleRepo.EXPECT().GetByID(
			ctx,
			vehicleID,
		).Return(expectedVehicle, nil)

		vehicleRepo.EXPECT().Update(
			ctx,
			vehicleID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, expectedErr)

		// when
		result, err := svc.Update(ctx, vehicleID, input)

		// then
		assert.Nil(t, result)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestVehicleService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("should successfully delete a vehicle", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		id := int32(1)

		vehicleRepo.EXPECT().GetByID(ctx, id).Return(getTestVehicles()[0], nil)
		vehicleRepo.EXPECT().Delete(ctx, id).Return(nil)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error if vehicle not found", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		id := int32(1)
		expectedErr := storage.ErrEntityNotFound
		vehicleRepo.EXPECT().GetByID(ctx, id).Return(nil, expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})

	t.Run("should return error if deleting vehicle fails", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		id := int32(4)
		expectedErr := errors.New("failed to delete")

		vehicleRepo.EXPECT().GetByID(ctx, id).Return(getTestVehicles()[0], nil)
		vehicleRepo.EXPECT().Delete(ctx, id).Return(expectedErr)

		// when
		err := svc.Delete(ctx, id)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, handleError(expectedErr).Error())
	})
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		vehicleRepo := storageMock.NewMockVehicleRepository(t)
		svc := NewVehicleService(vehicleRepo)

		// when
		ready := svc.Ready()

		// then
		assert.True(t, ready)
	})

	t.Run("should return false if the service is not ready", func(t *testing.T) {
		svc := NewVehicleService(nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}

func getTestVehicles() []*entities.Vehicle {
	now := time.Now()

	return []*entities.Vehicle{
		{
			ID:            1,
			CreatedAt:     now,
			UpdatedAt:     now,
			NumberPlate:   "FL TBZ 123",
			Description:   "Test description",
			Status:        entities.VehicleStatusActive,
			Type:          entities.VehicleTypeTrailer,
			WaterCapacity: 2000.5,
		},
		{
			ID:            2,
			CreatedAt:     now,
			UpdatedAt:     now,
			NumberPlate:   "FL TBZ 3456",
			Description:   "Test description",
			Status:        entities.VehicleStatusNotAvailable,
			Type:          entities.VehicleTypeTransporter,
			WaterCapacity: 1000.5,
		},
	}
}
