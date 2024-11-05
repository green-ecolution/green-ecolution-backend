package tree

import (
	"context"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTreeRepository_Create(t *testing.T) {
	t.Run("should create a tree with default values", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "", got.Species)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Nil(t, got.TreeCluster)
		assert.Nil(t, got.Sensor)
		assert.Equal(t, "", got.Number)
		assert.Equal(t, int32(0), got.PlantingYear)
		assert.Equal(t, float64(0), got.Latitude)
		assert.Equal(t, float64(0), got.Longitude)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, entities.WateringStatusUnknown, got.WateringStatus)
	})

	t.Run("should create a tree with all values set", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background(),
			WithSpecies("Oak"),
			WithTreeNumber("T001"),
			WithPlantingYear(2023),
			WithLatitude(54.801539),
			WithLongitude(9.446741),
			WithDescription("A newly planted oak tree"),
			WithWateringStatus(entities.WateringStatusGood),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Nil(t, got.TreeCluster)
		assert.Nil(t, got.Sensor)
		assert.Empty(t, got.Images)
		assert.Equal(t, "Oak", got.Species)
		assert.Equal(t, "T001", got.Number)
		assert.Equal(t, int32(2023), got.PlantingYear)
		assert.Equal(t, 54.801539, got.Latitude)
		assert.Equal(t, 9.446741, got.Longitude)
		assert.Equal(t, "A newly planted oak tree", got.Description)
		assert.Equal(t, entities.WateringStatusGood, got.WateringStatus)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Create(ctx, WithSpecies("Oak"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestTreeRepository_CreateAndLinkImages(t *testing.T) {
	t.Run("should create tree and link images successfully", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/tree")
		r := NewTreeRepository(suite.Store, mappers)
		sqlImages, err := suite.Store.GetAllImages(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		images := mappers.iMapper.FromSqlList(sqlImages)

		// when
		tree, createErr := r.CreateAndLinkImages(context.Background(),
			WithSpecies("Oak"),
			WithTreeNumber("T001"),
			WithLatitude(54.801539),
			WithLongitude(9.446741),
			WithPlantingYear(2023),
			WithDescription("Test tree with images"),
			WithImages(images),
		)

		// then
		assert.NoError(t, createErr)
		assert.NotNil(t, tree)
		assert.Equal(t, "Oak", tree.Species)
		assert.Equal(t, "T001", tree.Number)
		assert.Equal(t, 54.801539, tree.Latitude)
		assert.Equal(t, 9.446741, tree.Longitude)
		assert.Equal(t, int32(2023), tree.PlantingYear)
		assert.Equal(t, "Test tree with images", tree.Description)
		assert.NotEmpty(t, tree.Images)
		for i, img := range tree.Images {
			assert.Equal(t, images[i].ID, img.ID)
			assert.Equal(t, images[i].URL, img.URL)
			assert.Equal(t, *images[i].Filename, *img.Filename)
			assert.Equal(t, *images[i].MimeType, *img.MimeType)
		}
	})
}
