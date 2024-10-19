package treecluster

import (
	"context"
	"errors"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestGeoClusterLocator_UpdateCluster(t *testing.T) {
	t.Run("should update cluster location when successful", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)

		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{{ID: 1}, {ID: 2}},
		}
		expectedLat, expectedLong := 54.801539, 9.446741
		expectedRegion := &entities.Region{ID: 10}

		clusterRepo.EXPECT().GetByID(context.Background(), clusterID).Return(expectedCluster, nil)
		treeRepo.EXPECT().GetCenterPoint(context.Background(), []int32{1, 2}).Return(expectedLat, expectedLong, nil)
		regionRepo.EXPECT().GetByPoint(context.Background(), expectedLat, expectedLong).Return(expectedRegion, nil)
		clusterRepo.EXPECT().Update(context.Background(), clusterID, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		// when
		err := locator.UpdateCluster(context.Background(), &clusterID)

		// then
		assert.NoError(t, err)
	})

	t.Run("should do nothing when clusterID is nil", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
	
		// when
		err := locator.UpdateCluster(context.Background(), nil)
	
		// then
		assert.NoError(t, err)
		clusterRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("should return error when GetByID fails", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
	
		clusterID := int32(1)
		expectedError := storage.ErrTreeClusterNotFound
	
		clusterRepo.EXPECT().GetByID(context.Background(), clusterID).Return(nil, expectedError)
	
		// when
		err := locator.UpdateCluster(context.Background(), &clusterID)
	
		// then
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should return error when GetCenterPoint fails", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
	
		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{{ID: 1}, {ID: 2}},
		}
		expectedError := errors.New("empty geometry")
	
		clusterRepo.EXPECT().GetByID(context.Background(), clusterID).Return(expectedCluster, nil)
		treeRepo.EXPECT().GetCenterPoint(context.Background(), []int32{1, 2}).Return(0, 0, expectedError)
	
		// when
		err := locator.UpdateCluster(context.Background(), &clusterID)
	
		// then
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
 	
	t.Run("should return error when Update fails", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
	
		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{{ID: 1}, {ID: 2}},
		}
		expectedLat, expectedLong := 54.801539, 9.446741
		expectedRegion := &entities.Region{ID: 10}
		expectedError := errors.New("Update failed")
	
		clusterRepo.EXPECT().GetByID(context.Background(), clusterID).Return(expectedCluster, nil)
		treeRepo.EXPECT().GetCenterPoint(context.Background(), []int32{1, 2}).Return(expectedLat, expectedLong, nil)
		regionRepo.EXPECT().GetByPoint(context.Background(), expectedLat, expectedLong).Return(expectedRegion, nil)
		clusterRepo.EXPECT().Update(context.Background(), clusterID, mock.Anything, mock.Anything, mock.Anything).Return(nil, expectedError)
	
		// when
		err := locator.UpdateCluster(context.Background(), &clusterID)
	
		// then
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("should remove cluster coordinates when cluster has no trees", func(t *testing.T) {
		// given
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		locator := NewLocationUpdate(clusterRepo, treeRepo, regionRepo)
	
		clusterID := int32(1)
		expectedCluster := &entities.TreeCluster{
			ID:    clusterID,
			Trees: []*entities.Tree{},
		}
	
		clusterRepo.EXPECT().GetByID(context.Background(), clusterID).Return(expectedCluster, nil)
		clusterRepo.EXPECT().Update(context.Background(), clusterID, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	
		// when
		err := locator.UpdateCluster(context.Background(), &clusterID)
	
		// then
		assert.NoError(t, err)
	})
}
