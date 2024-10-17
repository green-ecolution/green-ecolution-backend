package region

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestRegionService_GetAll(t *testing.T) {
	t.Run("should return all regions", func(t *testing.T) {
		// given
		repo := storageMock.NewMockRegionRepository(t)
		svc := NewRegionService(repo)

		expectedRegions := []*entities.Region{
			{ID: 1, Name: "Region A"},
			{ID: 2, Name: "Region B"},
		}

		// when
		repo.EXPECT().GetAll(context.Background()).Return(expectedRegions, nil)
		regions, err := svc.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedRegions, regions)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// given
		repo := storageMock.NewMockRegionRepository(t)
		svc := NewRegionService(repo)

		repo.EXPECT().GetAll(context.Background()).Return(nil, storage.ErrRegionNotFound)
		regions, err := svc.GetAll(context.Background())

		// then
		assert.Error(t, err)
		assert.Equal(t, storage.ErrRegionNotFound, err)
		assert.Nil(t, regions)
	})
}

func TestRegionService_GetByID(t *testing.T) {
	t.Run("should return region when found", func(t *testing.T) {
		// given
		repo := storageMock.NewMockRegionRepository(t)
		svc := NewRegionService(repo)

		expectedRegion := &entities.Region{ID: 1, Name: "Region A"}

		// when
		repo.EXPECT().GetByID(context.Background(), int32(1)).Return(expectedRegion, nil)
		region, err := svc.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedRegion, region)
	})

	t.Run("should return error when region not found", func(t *testing.T) {
		// given
		repo := storageMock.NewMockRegionRepository(t)
		svc := NewRegionService(repo)

		// when
		repo.EXPECT().GetByID(context.Background(), int32(3)).Return(nil, storage.ErrRegionNotFound)
		region, err := svc.GetByID(context.Background(), 3)

		// then
		assert.Error(t, err)
		assert.Nil(t, region)
	})
}

func TestReady(t *testing.T) {
	t.Run("should return true if the service is ready", func(t *testing.T) {
		// given
		repo := storageMock.NewMockRegionRepository(t)
		svc := NewRegionService(repo)

		// when
		ready := svc.Ready()

		// then
		assert.True(t, ready)
	})

	t.Run("should return false if the service is not ready", func(t *testing.T) {
		// given
		svc := NewRegionService(nil)

		// when
		ready := svc.Ready()

		// then
		assert.False(t, ready)
	})
}

