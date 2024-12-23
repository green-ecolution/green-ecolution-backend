package tree_test

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTreeService_ImportTree(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a new tree when no matching tree exists", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, eventManager)

		expectedTree := TestTreesList[0]
		treeRepo.EXPECT().GetByCoordinates(ctx, TestTreeImport.Latitude, TestTreeImport.Longitude).Return(nil, nil)
		treeRepo.EXPECT().Create(ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(expectedTree, nil)

		// When
		err := svc.ImportTree(ctx, []*entities.TreeImport{TestTreeImport})

		// Then
		assert.NoError(t, err)
	})

	t.Run("should update an existing tree when matching coordinates and planting year are found", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, eventManager)

		existingTree := TestTreesList[0]
		updatedTree := TestTreesList[0]
		updatedTree.Description = "Updated description"

		treeRepo.EXPECT().GetByCoordinates(ctx, TestTreeImport.Latitude, TestTreeImport.Longitude).Return(existingTree, nil)
		treeRepo.EXPECT().GetByID(ctx, existingTree.ID).Return(existingTree, nil)
		treeRepo.EXPECT().Update(ctx, existingTree.ID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(updatedTree, nil)

		// When
		err := svc.ImportTree(ctx, []*entities.TreeImport{TestTreeImport})

		// Then
		assert.NoError(t, err)
	})
	t.Run("should delete and recreate a tree when planting years differ", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, eventManager)

		// Define existing tree and tree import data
		existingTree := TestTreesList[0]
		TestTreeImport.PlantingYear = 2024

		treeRepo.EXPECT().GetByCoordinates(ctx, TestTreeImport.Latitude, TestTreeImport.Longitude).Return(existingTree, nil)
		treeRepo.EXPECT().GetByID(ctx, existingTree.ID).Return(existingTree, nil)
		treeRepo.EXPECT().Delete(ctx, existingTree.ID).Return(nil)
		treeRepo.EXPECT().Create(ctx,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).Return(TestTreesList[1], nil)

		// Call ImportTree
		err := svc.ImportTree(ctx, []*entities.TreeImport{TestTreeImport})

		// Assert that no error occurred
		assert.NoError(t, err)
	})
	t.Run("should return error if processing delete queue fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, eventManager)

		existingTree := TestTreesList[0]
		expectedErr := errors.New("error deleting tree")

		treeRepo.EXPECT().GetByCoordinates(ctx, TestTreeImport.Latitude, TestTreeImport.Longitude).Return(existingTree, nil)
		treeRepo.EXPECT().GetByID(ctx, existingTree.ID).Return(existingTree, nil)
		treeRepo.EXPECT().Delete(ctx, existingTree.ID).Return(expectedErr)

		// When
		err := svc.ImportTree(ctx, []*entities.TreeImport{TestTreeImport})

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, "500: error deleting tree")
	})
}
