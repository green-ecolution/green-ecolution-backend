package tree

import (
	"context"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	mappers TreeMappers
	suite   *testutils.PostgresTestSuite
)

func TestMain(m *testing.M) {
	code := 1
	ctx := context.Background()
	defer func() { os.Exit(code) }()
	suite = testutils.SetupPostgresTestSuite(ctx)
	mappers = NewTreeRepositoryMappers(
		&generated.InternalTreeRepoMapperImpl{},
		&generated.InternalImageRepoMapperImpl{},
		&generated.InternalSensorRepoMapperImpl{},
		&generated.InternalTreeClusterRepoMapperImpl{},
	)
	defer suite.Terminate(ctx)
	code = m.Run()
}

func TestTreeRepository_Delete(t *testing.T) {
	t.Run("should delete tree successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)

		// when
		err := r.Delete(context.Background(), treeID)

		// then
		assert.NoError(t, err)

		deletedTree, err := r.GetByID(context.Background(), treeID)
		assert.Error(t, err)
		assert.Nil(t, deletedTree)
	})

	t.Run("should delete tree and linked images successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)
		sqlImages, err := suite.Store.GetAllImagesByTreeID(context.Background(), treeID)
		if err != nil {
			t.Fatal(err)
		}
		images := mappers.iMapper.FromSqlList(sqlImages)
		assert.NotEmpty(t, images)

		// when
		err = r.Delete(context.Background(), treeID)

		// then
		assert.NoError(t, err)

		deletedTree, err := r.GetByID(context.Background(), treeID)
		assert.Error(t, err)
		assert.Nil(t, deletedTree)

		for _, img := range images {
			_, err = suite.Store.GetImageByID(context.Background(), img.ID)
			assert.Error(t, err)
		}
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if tree ID is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if tree ID is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 0)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.Delete(ctx, 1)

		// then
		assert.Error(t, err)
	})
}

func TestTreeRepository_DeleteAndUnlinkImages(t *testing.T) {
	t.Run("should delete tree and unlink all images successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)
		sqlImages, err := suite.Store.GetAllImagesByTreeID(context.Background(), treeID)
		if err != nil {
			t.Fatal(err)
		}
		images := mappers.iMapper.FromSqlList(sqlImages)
		assert.NotEmpty(t, images)

		// when
		err = r.DeleteAndUnlinkImages(context.Background(), treeID)
		assert.NoError(t, err)

		deletedTree, err := r.GetByID(context.Background(), treeID)

		// then
		assert.Error(t, err)
		assert.Nil(t, deletedTree)

		for _, img := range images {
			_, err = suite.Store.GetImageByID(context.Background(), img.ID)
			assert.NoError(t, err)
		}
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.DeleteAndUnlinkImages(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if tree ID is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.DeleteAndUnlinkImages(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if tree ID is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.DeleteAndUnlinkImages(context.Background(), 0)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error by canceling the context", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.DeleteAndUnlinkImages(ctx, 1)

		// then
		assert.Error(t, err)
	})
}

func TestTreeRepository_UnlinkImage(t *testing.T) {
	t.Run("should unlink a specific image from tree successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeID := int32(1)
		imageID := int32(1)

		isContained, err := isContainedInTreeImages(imageID, treeID)
		assert.NoError(t, err)
		assert.True(t, isContained)

		// when
		err = r.UnlinkImage(context.Background(), treeID, imageID)

		// then
		assert.NoError(t, err)
		isContained, err = isContainedInTreeImages(imageID, treeID)
		assert.NoError(t, err)
		assert.False(t, isContained)
	})

	t.Run("should return error if tree not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 99, 1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if image not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 1, 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if tree is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), -1, 1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if tree is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 0, 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if image is negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 1, -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error if image is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 1, 0)

		// then
		assert.Error(t, err)
	})
}

func TestTreeRepository_UnlinkTreeClusterID(t *testing.T) {
	t.Run("should unlink tree cluster ID from a tree successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeClusterID := int32(1)
		trees, err := suite.Store.GetTreesByTreeClusterID(context.Background(), &treeClusterID)
		assert.NoError(t, err)
		assert.NotEmpty(t, trees)

		// when
		err = r.UnlinkTreeClusterID(context.Background(), treeClusterID)

		// then
		assert.NoError(t, err)
		trees, err = suite.Store.GetTreesByTreeClusterID(context.Background(), &treeClusterID)
		assert.NoError(t, err)
		assert.Empty(t, trees)
	})

	t.Run("should return empty trees when tree cluster is not set in any tree", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		treeClusterID := int32(6)
		trees, err := suite.Store.GetTreesByTreeClusterID(context.Background(), &treeClusterID)
		assert.NoError(t, err)
		assert.Empty(t, trees)

		// when
		err = r.UnlinkTreeClusterID(context.Background(), treeClusterID)

		// then
		assert.NoError(t, err)
		trees, err = suite.Store.GetTreesByTreeClusterID(context.Background(), &treeClusterID)
		assert.NoError(t, err)
		assert.Empty(t, trees)
	})

	t.Run("should return error when tree cluster not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkTreeClusterID(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when tree cluster id negative", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkTreeClusterID(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when tree cluster id is zero", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		err := r.UnlinkTreeClusterID(context.Background(), 0)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewTreeRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.UnlinkTreeClusterID(ctx, 1)

		// then
		assert.Error(t, err)
	})
}

func contains(id int32, list []int32) bool {
	found := false
	for _, image := range list {
		if image == id {
			found = true
			break
		}
	}
	return found
}

func isContainedInTreeImages(imageID, treeID int32) (bool, error) {
	sqlImages, err := suite.Store.GetAllImagesByTreeID(context.Background(), treeID)
	if err != nil {
		return false, err
	}
	images := mappers.iMapper.FromSqlList(sqlImages)
	imageIDs := make([]int32, len(images))
	for i, image := range images {
		imageIDs[i] = image.ID
	}
	return contains(imageID, imageIDs), nil
}
