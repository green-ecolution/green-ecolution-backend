package treecluster

import (
	"context"
	"errors"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
)

func TestGeoClusterLocator_UpdateCluster(t *testing.T) {
	t.Run("should update cluster location when successful", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewGeoLocation(treeRepo, regionRepo)

		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{{ID: 1}, {ID: 2}},
		}
		expectedLat, expectedLong := 54.801539, 9.446741
		expectedRegion := &entities.Region{ID: 10}

		treeRepo.EXPECT().GetCenterPoint(context.Background(), []int32{1, 2}).Return(expectedLat, expectedLong, nil)
		regionRepo.EXPECT().GetByPoint(context.Background(), expectedLat, expectedLong).Return(expectedRegion, nil)

		// when
		err := locator.UpdateCluster(context.Background(), expectedCluster)

		// then
		assert.NoError(t, err)

	})

	t.Run("should do nothing when clusterID is nil", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewGeoLocation(treeRepo, regionRepo)

		// when
		err := locator.UpdateCluster(context.Background(), nil)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when GetCenterPoint fails", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewGeoLocation(treeRepo, regionRepo)

		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{{ID: 1}, {ID: 2}},
		}
		expectedError := errors.New("400: empty geometry")

		treeRepo.EXPECT().GetCenterPoint(context.Background(), []int32{1, 2}).Return(0, 0, expectedError)

		// when
		err := locator.UpdateCluster(context.Background(), expectedCluster)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "400: empty geometry")
	})

	t.Run("should return error when getRegionByID fails", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewGeoLocation(treeRepo, regionRepo)

		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{{ID: 1}, {ID: 2}},
		}
		expectedLat, expectedLong := 54.801539, 9.446741
		expectedError := errors.New("500: get region failed")

		treeRepo.EXPECT().GetCenterPoint(context.Background(), []int32{1, 2}).Return(expectedLat, expectedLong, nil)
		regionRepo.EXPECT().GetByPoint(context.Background(), expectedLat, expectedLong).Return(nil, expectedError)

		// when
		err := locator.UpdateCluster(context.Background(), expectedCluster)

		// then
		assert.Error(t, err)
		assert.EqualError(t, err, "500: get region failed")
	})

	t.Run("should remove cluster coordinates when cluster has no trees", func(t *testing.T) {
		// given
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewGeoLocation(treeRepo, regionRepo)

		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{},
		}

		// when
		err := locator.UpdateCluster(context.Background(), expectedCluster)

		// then
		assert.NoError(t, err)
	})
}
