package tree

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	service "github.com/green-ecolution/green-ecolution-backend/internal/service/_mock"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTreeService_ImportTree(t *testing.T) {
	ctx := context.Background()

	treeToImport := &entities.TreeImport{
		Latitude:     54.801539,
		Longitude:    9.446741,
		PlantingYear: 2023,
		Species:      "Oak",
		Number:       "T001",
	}

	t.Run("should create a new tree when no matching tree exists", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)
		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		expectedTree := getTestTrees()[0]
		treeRepo.EXPECT().GetByCoordinates(ctx, treeToImport.Latitude, treeToImport.Longitude).Return(nil, nil)
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
		err := svc.ImportTree(ctx, []*entities.TreeImport{treeToImport})

		// Then
		assert.NoError(t, err)
	})

	t.Run("should update an existing tree when matching coordinates and planting year are found", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		existingTree := getTestTrees()[0]
		updatedTree := getTestTrees()[0]
		updatedTree.Description = "Updated description"

		treeRepo.EXPECT().GetByCoordinates(ctx, treeToImport.Latitude, treeToImport.Longitude).Return(existingTree, nil)
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
		err := svc.ImportTree(ctx, []*entities.TreeImport{treeToImport})

		// Then
		assert.NoError(t, err)
	})
	t.Run("should delete and recreate a tree when planting years differ", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		// Define existing tree and tree import data
		existingTree := getTestTrees()[0]
		treeToImport.PlantingYear = 2024

		treeRepo.EXPECT().GetByCoordinates(ctx, treeToImport.Latitude, treeToImport.Longitude).Return(existingTree, nil)
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
			mock.Anything).Return(getTestTrees()[1], nil)

		// Call ImportTree
		err := svc.ImportTree(ctx, []*entities.TreeImport{treeToImport})

		// Assert that no error occurred
		assert.NoError(t, err)
	})
	t.Run("should return error if processing delete queue fails", func(t *testing.T) {
		// Given
		treeRepo := storageMock.NewMockTreeRepository(t)
		sensorRepo := storageMock.NewMockSensorRepository(t)
		imageRepo := storageMock.NewMockImageRepository(t)
		treeClusterRepo := storageMock.NewMockTreeClusterRepository(t)
		locator := service.NewMockGeoClusterLocator(t)

		svc := NewTreeService(treeRepo, sensorRepo, imageRepo, treeClusterRepo, locator)

		existingTree := getTestTrees()[0]
		expectedErr := errors.New("error deleting tree")

		treeRepo.EXPECT().GetByCoordinates(ctx, treeToImport.Latitude, treeToImport.Longitude).Return(existingTree, nil)
		treeRepo.EXPECT().GetByID(ctx, existingTree.ID).Return(existingTree, nil)
		treeRepo.EXPECT().Delete(ctx, existingTree.ID).Return(expectedErr)

		// When
		err := svc.ImportTree(ctx, []*entities.TreeImport{treeToImport})

		// then
		assert.Error(t, err)
		assert.ErrorContains(t, err, expectedErr.Error())
	})
}

func getTestTrees() []*entities.Tree {
	now := time.Now()

	return []*entities.Tree{
		{
			ID:           1,
			CreatedAt:    now,
			UpdatedAt:    now,
			Species:      "Oak",
			Number:       "T001",
			Latitude:     9.446741,
			Longitude:    54.801539,
			Description:  "A mature oak tree",
			PlantingYear: 2023,
			Readonly:     true,
		},
		{
			ID:           2,
			CreatedAt:    now,
			UpdatedAt:    now,
			Species:      "Pine",
			Number:       "T002",
			Latitude:     9.446700,
			Longitude:    54.801510,
			Description:  "A young pine tree",
			PlantingYear: 2023,
			Readonly:     true,
		},
	}
}
