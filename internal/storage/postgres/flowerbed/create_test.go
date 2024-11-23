package flowerbed

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRegionRepository_Create(t *testing.T) {
	region := &entities.Region{ID: 1, Name: "Mürwik"}
	sensor := &entities.Sensor{ID: "sensor-1", Latitude: 54.82124518093376, Longitude: 9.485702120628517, Status: entities.SensorStatusOnline}

	input := entities.Flowerbed{
		Description:    "New description",
		Size:           40000,
		NumberOfPlants: 100000,
		MoistureLevel:  3.5,
		Address:        "126 Water street",
		Latitude:       utils.P(54.81269326939148),
		Longitude:      utils.P(54.81269326939148),
		Region:         region,
		Sensor:         sensor,
	}

	t.Run("should create flowerbed with region, longitude and latitude", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, input.Latitude, got.Latitude)
		assert.Equal(t, input.Longitude, got.Longitude)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Empty(t, got.Sensor)
		assert.Empty(t, got.Images)
		assert.Equal(t, "", got.Address)
		assert.Equal(t, 0.0, got.MoistureLevel)
		assert.Equal(t, int32(0), got.NumberOfPlants)
		assert.Equal(t, 0.0, got.Size)
		assert.False(t, got.Archived)

		// assert region
		assert.Equal(t, region.ID, got.Region.ID)
		assert.Equal(t, region.Name, got.Region.Name)
	})

	t.Run("should create flowerbed with all values set", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithDescription(input.Description),
			WithAddress(input.Address),
			WithMoistureLevel(input.MoistureLevel),
			WithNumberOfPlants(input.NumberOfPlants),
			WithSensorID(input.Sensor.ID),
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
			WithSize(input.Size),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, input.Description, got.Description)
		assert.Equal(t, input.Latitude, got.Latitude)
		assert.Equal(t, input.Longitude, got.Longitude)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Equal(t, input.Address, got.Address)
		assert.Empty(t, got.Images)
		assert.Equal(t, input.MoistureLevel, got.MoistureLevel)
		assert.Equal(t, input.NumberOfPlants, got.NumberOfPlants)
		assert.Equal(t, input.Size, got.Size)
		assert.False(t, got.Archived)

		// assert region
		assert.Equal(t, region.ID, got.Region.ID)
		assert.Equal(t, region.Name, got.Region.Name)

		// assert sensor
		assert.Equal(t, sensor.ID, got.Sensor.ID)
		assert.Equal(t, sensor.Status, got.Sensor.Status)
	})

	t.Run("should return error when flowerbed is invalid", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Create(context.Background())

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed latitude is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithLatitude(nil),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed longitude is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(nil),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed region is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Create(
			context.Background(),
			WithLatitude(input.Latitude),
			WithRegionID(0),
			WithLongitude(input.Longitude),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewFlowerbedRepository(suite.Store, mappers)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Create(
			ctx,
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestRegionRepository_CreateAndLinkImages(t *testing.T) {
	region := &entities.Region{ID: 1, Name: "Mürwik"}
	images := []*entities.Image{
		{
			ID:       1,
			URL:      "/test/url/to/image",
			Filename: utils.P("Screenshot"),
			MimeType: utils.P("png"),
		},
	}
	imageIDs := []int32{1}

	input := entities.Flowerbed{
		Latitude:  utils.P(54.81269326939148),
		Longitude: utils.P(54.81269326939148),
		Region:    region,
	}

	t.Run("should create flowerbed with all values set", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.CreateAndLinkImages(
			context.Background(),
			WithDescription(input.Description),
			WithImagesIDs(imageIDs),
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, input.Latitude, got.Latitude)
		assert.Equal(t, input.Longitude, got.Longitude)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Empty(t, got.Sensor)
		assert.Equal(t, "", got.Address)
		assert.Equal(t, 0.0, got.MoistureLevel)
		assert.Equal(t, int32(0), got.NumberOfPlants)
		assert.Equal(t, 0.0, got.Size)
		assert.False(t, got.Archived)

		// assert region
		assert.Equal(t, region.ID, got.Region.ID)
		assert.Equal(t, region.Name, got.Region.Name)

		// assert images
		for i, fb := range got.Images {
			assert.Equal(t, images[i].ID, fb.ID)
			assert.Equal(t, images[i].Filename, fb.Filename)
			assert.Equal(t, images[i].MimeType, fb.MimeType)
			assert.Equal(t, images[i].URL, fb.URL)
		}

		retrievedImages, err := r.GetAllImagesByID(context.Background(), int32(1))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(retrievedImages))
	})

	t.Run("should create flowerbed with no images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.CreateAndLinkImages(
			context.Background(),
			WithDescription(input.Description),
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "", got.Description)
		assert.Equal(t, input.Latitude, got.Latitude)
		assert.Equal(t, input.Longitude, got.Longitude)
		assert.NotZero(t, got.ID)
		assert.WithinDuration(t, got.CreatedAt, time.Now(), time.Second)
		assert.WithinDuration(t, got.UpdatedAt, time.Now(), time.Second)
		assert.Empty(t, got.Sensor)
		assert.Equal(t, "", got.Address)
		assert.Equal(t, 0.0, got.MoistureLevel)
		assert.Equal(t, int32(0), got.NumberOfPlants)
		assert.Equal(t, 0.0, got.Size)
		assert.False(t, got.Archived)

		// assert region
		assert.Equal(t, region.ID, got.Region.ID)
		assert.Equal(t, region.Name, got.Region.Name)

		// assert images
		assert.Empty(t, got.Images)
	})

	t.Run("should return error for invalid image ID", func(t *testing.T) {
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		got, err := r.CreateAndLinkImages(
			context.Background(),
			WithDescription(input.Description),
			WithImagesIDs([]int32{99}),
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
		)

		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/treecluster")
		r := NewFlowerbedRepository(suite.Store, mappers)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.CreateAndLinkImages(
			ctx,
			WithDescription(input.Description),
			WithImagesIDs([]int32{1}),
			WithLatitude(input.Latitude),
			WithRegionID(input.Region.ID),
			WithLongitude(input.Longitude),
		)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
