package tree

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestTreeRepository_Update(t *testing.T) {
	t.Run("should update tree successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)
		date := time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC)

		newSpecies := "Updated Species"
		newNumber := "UpdatedNumber"
		newLatitude := 55.123456
		newLongitude := 10.654321
		newPlantingYear := int32(2025)
		newDescription := "Updated description"
		newWateringStatus := entities.WateringStatusGood
		newLastWateredValue := &date
		newProvider := "foo-provider"

		// when
		updatedTree, err := r.Update(context.Background(), treeID, func(tree *entities.Tree) (bool, error) {
			tree.Species = newSpecies
			tree.Number = newNumber
			tree.Latitude = newLatitude
			tree.Longitude = newLongitude
			tree.PlantingYear = newPlantingYear
			tree.Provider = newProvider
			tree.Description = newDescription
			tree.WateringStatus = newWateringStatus
			tree.LastWatered = newLastWateredValue
			return true, nil
		})

		// then
		assert.NoError(t, err)
		assert.NotNil(t, updatedTree)
		assert.Equal(t, newSpecies, updatedTree.Species, "Species should match")
		assert.Equal(t, newNumber, updatedTree.Number, "Tree Number should match")
		assert.Equal(t, newLatitude, updatedTree.Latitude, "Latitude should match")
		assert.Equal(t, newLongitude, updatedTree.Longitude, "Longitude should match")
		assert.Equal(t, newPlantingYear, updatedTree.PlantingYear, "Planting Year should match")
		assert.Equal(t, newProvider, updatedTree.Provider, "Provider should match")
		assert.Equal(t, newDescription, updatedTree.Description, "Description should match")
		assert.Equal(t, newWateringStatus, updatedTree.WateringStatus, "Watering Status should match")
		assert.Equal(t, newLastWateredValue, updatedTree.LastWatered, "Last watered should match")
	})

	t.Run("should return error when tree is not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		updatedTree, err := r.Update(context.Background(), int32(99), func(tree *entities.Tree) (bool, error) {
			tree.Species = "Non-existent species"
			return true, nil
		})

		// then
		assert.Error(t, err)
		assert.Nil(t, updatedTree)
	})

	t.Run("should return error the id is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		updatedTree, err := r.Update(context.Background(), int32(-1), func(tree *entities.Tree) (bool, error) {
			tree.Species = "species"
			return true, nil
		})

		// then
		assert.Error(t, err)
		assert.Nil(t, updatedTree)
	})

	t.Run("should return error the id is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		updatedTree, err := r.Update(context.Background(), int32(0), func(tree *entities.Tree) (bool, error) {
			tree.Species = "species"
			return true, nil
		})

		// then
		assert.Error(t, err)
		assert.Nil(t, updatedTree)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		updatedTree, err := r.Update(ctx, int32(1), func(tree *entities.Tree) (bool, error) {
			tree.Species = "Canceled context species"
			return true, nil
		})

		// then
		assert.Error(t, err)
		assert.Nil(t, updatedTree)
	})

	t.Run("should update tree and link new images successfully when tree has other images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)
		newSpecies := "Updated Species"
		newDescription := "Updated description"

		sqlImages, err := suite.Store.GetAllImages(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		images := mappers.iMapper.FromSqlList(sqlImages)

		var updateImages []*entities.Image
		for _, image := range images {
			if image.ID == 3 || image.ID == 4 {
				updateImages = append(updateImages, image)
			}
		}
		assert.NotEmpty(t, updateImages)
		assert.Len(t, updateImages, 2)

		tree, getErr := r.GetByID(context.Background(), treeID)
		assert.NoError(t, getErr)

		// when
		updatedTree, updateErr := r.Update(
			context.Background(),
			treeID,
			func(tree *entities.Tree) (bool, error) {
				tree.Species = newSpecies
				tree.Description = newDescription
				tree.Images = updateImages
				return true, nil
			},
		)

		sqlImages, err = suite.Store.GetAllImagesByTreeID(context.Background(), treeID)
		if err != nil {
			t.Fatal(err)
		}
		images = mappers.iMapper.FromSqlList(sqlImages)

		// then
		assert.NoError(t, updateErr)
		assert.NotNil(t, updatedTree)
		assert.Equal(t, tree.ID, updatedTree.ID)
		assert.Equal(t, newSpecies, updatedTree.Species, "Species should match updated value")
		assert.Equal(t, newDescription, updatedTree.Description, "Description should match updated value")
		assert.NotEqual(t, updatedTree.Images, len(tree.Images), "Image count should not be equal to the original tree")
		for i, img := range updatedTree.Images {
			assert.Equal(t, images[i].ID, img.ID, "Image ID should match")
			assert.Equal(t, images[i].URL, img.URL, "Image URL should match")
			assert.Equal(t, images[i].MimeType, img.MimeType, "Image MimeType should match")
			assert.Equal(t, images[i].Filename, img.Filename, "Image Filename should match")
		}
	})

	t.Run("should update tree and link new images successfully when tree has no other images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(2)
		newSpecies := "Updated Species"
		newDescription := "Updated description"

		sqlImages, err := suite.Store.GetAllImages(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		images := mappers.iMapper.FromSqlList(sqlImages)
		assert.NotEmpty(t, images)

		tree, err := r.GetByID(context.Background(), treeID)
		assert.NoError(t, err)
		assert.Empty(t, tree.Images)

		// when
		updatedTree, err := r.Update(
			context.Background(),
			treeID,
			func(tree *entities.Tree) (bool, error) {
				tree.Species = newSpecies
				tree.Description = newDescription
				tree.Images = images
				return true, nil
			},
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, updatedTree)
		assert.Equal(t, tree.ID, updatedTree.ID)
		assert.Equal(t, newSpecies, updatedTree.Species, "Species should match updated value")
		assert.Equal(t, newDescription, updatedTree.Description, "Description should match updated value")
		assert.NotEqual(t, updatedTree.Images, len(tree.Images), "Image count should not be equal to the original tree")
		for i, img := range updatedTree.Images {
			assert.Equal(t, images[i].ID, img.ID, "Image ID should match")
			assert.Equal(t, images[i].URL, img.URL, "Image URL should match")
			assert.Equal(t, images[i].MimeType, img.MimeType, "Image MimeType should match")
			assert.Equal(t, images[i].Filename, img.Filename, "Image Filename should match")
		}
	})
}

//
//func TestTreeRepository_UpdateTreeClusterID(t *testing.T) {
//	t.Run("should update the assignment of tree cluster ids successfully ", func(t *testing.T) {
//		// given
//		suite.ResetDB(t)
//		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		treeIDs := []int32{3, 4}
//		newTreeClusterID := int32(2)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(context.Background(), treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.NoError(t, err)
//		for _, id := range treeIDs {
//			updatedTree, err := r.GetByID(context.Background(), id)
//			assert.NoError(t, err)
//			assert.NotNil(t, updatedTree)
//			assert.NotNil(t, updatedTree.TreeCluster)
//			assert.Equal(t, *expectedTreeClusterID, updatedTree.TreeCluster.ID)
//		}
//	})
//
//	t.Run("should return error when trees are not found", func(t *testing.T) {
//		// given
//		suite.ResetDB(t)
//		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		treeIDs := []int32{99, 199}
//		newTreeClusterID := int32(2)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(context.Background(), treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.Error(t, err)
//	})
//
//	t.Run("should return error when a tree id is negative", func(t *testing.T) {
//		// given
//		suite.ResetDB(t)
//		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		treeIDs := []int32{1, -1}
//		newTreeClusterID := int32(2)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(context.Background(), treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.Error(t, err)
//	})
//
//	t.Run("should return error when a tree id is zero", func(t *testing.T) {
//		// given
//		suite.ResetDB(t)
//		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		treeIDs := []int32{1, -1}
//		newTreeClusterID := int32(2)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(context.Background(), treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.Error(t, err)
//	})
//
//	t.Run("should return error when tree cluster is not found", func(t *testing.T) {
//		// given
//		suite.ResetDB(t)
//		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		treeIDs := []int32{3, 4}
//		newTreeClusterID := int32(99)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(context.Background(), treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.Error(t, err)
//	})
//
//	t.Run("should return error when tree cluster id negative", func(t *testing.T) {
//		// given
//		suite.ResetDB(t)
//		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		treeIDs := []int32{3, 4}
//		newTreeClusterID := int32(-1)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(context.Background(), treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.Error(t, err)
//	})
//
//	t.Run("should return error when tree cluster id is zero", func(t *testing.T) {
//		// given
//		suite.ResetDB(t)
//		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		treeIDs := []int32{3, 4}
//		newTreeClusterID := int32(0)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(context.Background(), treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.Error(t, err)
//	})
//
//	t.Run("should return error when context is canceled", func(t *testing.T) {
//		// given
//		r := TreeRepository{store: suite.Store, TreeMappers: mappers}
//		ctx, cancel := context.WithCancel(context.Background())
//		cancel()
//
//		treeIDs := []int32{3, 4}
//		newTreeClusterID := int32(2)
//		expectedTreeClusterID := &newTreeClusterID
//
//		// when
//		err := r.updateTreeClusterID(ctx, treeIDs, expectedTreeClusterID)
//
//		// then
//		assert.Error(t, err)
//	})
//}
