package flowerbed

import (
	"context"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	mappers FlowerbedMappers
	suite   *testutils.PostgresTestSuite
)

func TestMain(m *testing.M) {
	code := 1
	ctx := context.Background()
	defer func() { os.Exit(code) }()
	suite = testutils.SetupPostgresTestSuite(ctx)
	mappers = NewFlowerbedMappers(
		&generated.InternalFlowerbedRepoMapperImpl{},
		&generated.InternalImageRepoMapperImpl{},
		&generated.InternalSensorRepoMapperImpl{},
		&generated.InternalRegionRepoMapperImpl{},
	)
	defer suite.Terminate(ctx)

	code = m.Run()
}

func TestRegionRepository_Delete(t *testing.T) {
	t.Run("should delete flowerbed", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 2)

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when flowerbed has linked images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when flowerbed not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		// when
		err := r.Delete(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.Delete(ctx, 1)

		// then
		assert.Error(t, err)
	})
}

func TestRegionRepository_DeleteAndUnlinkImages(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")

	t.Run("should unlink images and delete flowerbed", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.DeleteAndUnlinkImages(context.Background(), 1)

		// then
		assert.NoError(t, err)

		_, err = r.GetByID(context.Background(), 1)
		assert.Error(t, err)
	})

	t.Run("should return error when flowerbed has no images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// then
		err := r.DeleteAndUnlinkImages(context.Background(), 2)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when flowerbed not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// then
		err := r.DeleteAndUnlinkImages(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// then
		err := r.DeleteAndUnlinkImages(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.DeleteAndUnlinkImages(ctx, 1)

		// then
		assert.Error(t, err)
	})
}

func TestRegionRepository_UnlinkImage(t *testing.T) {
	t.Run("should unlink image from flowerbed", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 1, 1)

		// then
		assert.NoError(t, err)

		images, err := r.GetAllImagesByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotContains(t, images, &entities.Image{ID: 1})
	})

	t.Run("should return error when flowerbed_id is not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 99, 1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when image_id is not found", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.UnlinkImage(context.Background(), 1, 99)

		// then
		assert.Error(t, err)
	})
}

func TestRegionRepository_UnlinkAllImages(t *testing.T) {
	t.Run("should unlink all images from flowerbed", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.UnlinkAllImages(context.Background(), 1)

		// then
		assert.NoError(t, err)

		images, err := r.GetAllImagesByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Empty(t, images)
	})

	t.Run("should return error when unlinking all images fails", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.UnlinkAllImages(context.Background(), 99)

		// then
		assert.Error(t, err)
	})
}

func TestTreeClusterRepository_Archived(t *testing.T) {
	t.Run("should archive flowerbed", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.Archive(context.Background(), 1)
		got, errGot := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NoError(t, errGot)
		assert.NotNil(t, got)
		assert.True(t, got.Archived)
	})

	t.Run("should return error when flowerbed with non-existing id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.Archive(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		err := r.Archive(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should not return error when archive flowerbed twice", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)
		err := r.Archive(context.Background(), 1)
		assert.NoError(t, err)

		// when
		gotBefore, errGotBefore := r.GetByID(context.Background(), 1)
		err = r.Archive(context.Background(), 1)
		gotAfter, errGotAfter := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NoError(t, errGotBefore)
		assert.NoError(t, errGotAfter)
		assert.NotNil(t, gotBefore)
		assert.NotNil(t, gotAfter)
		assert.True(t, gotBefore.Archived)
		assert.True(t, gotAfter.Archived)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.Archive(ctx, 1)

		// then
		assert.Error(t, err)
	})
}