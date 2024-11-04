package flowerbed

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestFlowerbedRepository_GetAll(t *testing.T) {
	t.Run("should return all flowerbeds", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		assert.Len(t, got, len(allTestFlowerbeds))
		for i, fb := range got {
			assert.Equal(t, allTestFlowerbeds[i].ID, fb.ID)
			assert.Equal(t, allTestFlowerbeds[i].Description, fb.Description)

			// assert region
			assert.NotNil(t, fb.Region)
			assert.Equal(t, allTestFlowerbeds[i].RegionID, fb.Region.ID)

			// assert sensor
			if allTestFlowerbeds[i].SensorID == -1 {
				assert.Nil(t, fb.Sensor)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, fb.Sensor)
				assert.Equal(t, allTestFlowerbeds[i].SensorID, fb.Sensor.ID)
			}

			// assert images
			if len(allTestFlowerbeds[i].Images) == 0 {
				assert.Equal(t, []*entities.Image{}, fb.Images)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, fb.Images)
				assert.Equal(t, len(allTestFlowerbeds[i].Images), len(fb.Images))
				assert.Equal(t, allTestFlowerbeds[i].Images[0].ID, fb.Images[0].ID)
			}
		}
	})

	t.Run("should return empty slice when db is empty", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		_, err := r.GetAll(ctx)

		// then
		assert.Error(t, err)
	})
}

func TestFlowerbedRepository_GetByID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")

	t.Run("should return flowerbed by id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 1)

		expectedF := allTestFlowerbeds[0]

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, expectedF.ID, got.ID)
		assert.Equal(t, expectedF.Description, got.Description)

		// assert region
		assert.NotNil(t, got.Region)
		assert.Equal(t, allTestFlowerbeds[0].RegionID, got.Region.ID)

		// assert sensor
		if expectedF.SensorID == -1 {
			assert.Nil(t, got.Sensor)
			assert.NoError(t, err)
		} else {
			assert.NotNil(t, got.Sensor)
			assert.Equal(t, expectedF.SensorID, got.Sensor.ID)
		}

		// assert images
		if len(expectedF.Images) == 0 {
			assert.Equal(t, []*entities.Image{}, got.Images)
			assert.NoError(t, err)
		} else {
			assert.NotNil(t, got.Images)
			assert.Equal(t, len(expectedF.Images), len(got.Images))
			assert.Equal(t, expectedF.Images[0].ID, got.Images[0].ID)
		}
	})

	t.Run("should return error when flowerbed with non-existing id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetByID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with zero id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestFlowerbedRepository_GetAllImagesByID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")

	t.Run("should return all linked images to flowerbed", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetAllImagesByID(context.Background(), 1)

		expectedImages := allTestFlowerbeds[0].Images

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, len(expectedImages), len(got))
		for i, fb := range got {
			assert.Equal(t, expectedImages[i].ID, fb.ID)
			assert.Equal(t, expectedImages[i].Filename, fb.Filename)
			assert.Equal(t, expectedImages[i].MimeType, fb.MimeType)
			assert.Equal(t, expectedImages[i].URL, fb.URL)
		}
	})

	t.Run("should return empty slice when no images are found", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetAllImagesByID(context.Background(), 2)

		// then
		assert.NoError(t, err)
		assert.Equal(t, 0, len(got))
	})

	t.Run("should return error when flowerbed with non-existing id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetAllImagesByID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetAllImagesByID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetAllImagesByID(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetAllImagesByID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestFlowerbedRepository_GetSensorByFlowerbedID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")

	t.Run("should return linked sensor to flowerbed", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByFlowerbedID(context.Background(), 1)

		expectedSensor := allTestFlowerbeds[0].SensorID

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotNil(t, got)
		assert.Equal(t, expectedSensor, got.ID)
	})

	t.Run("should return error when no sensor is found", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByFlowerbedID(context.Background(), 2)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with non-existing id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByFlowerbedID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByFlowerbedID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetSensorByFlowerbedID(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetSensorByFlowerbedID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestFlowerbedRepository_GetRegionByFlowerbedID(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")

	t.Run("should return linked region to flowerbed", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetRegionByFlowerbedID(context.Background(), 1)

		expectedRegion := allTestFlowerbeds[0].RegionID

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.NotNil(t, got)
		assert.Equal(t, expectedRegion, got.ID)
	})

	t.Run("should return error when flowerbed with non-existing id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetRegionByFlowerbedID(context.Background(), 99)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetRegionByFlowerbedID(context.Background(), -1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.GetRegionByFlowerbedID(context.Background(), 0)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.GetRegionByFlowerbedID(ctx, 1)

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

type testFlowerbeds struct {
	ID          int32
	Description string
	RegionID    int32
	SensorID    int32
	Images      []*entities.Image
}

var allTestFlowerbeds = []*testFlowerbeds{
	{
		ID:          1,
		Description: "Big flowerbed nearby the sea",
		RegionID:    1,
		SensorID:    2,
		Images: []*entities.Image{
			{
				ID:       1,
				URL:      "/test/url/to/image",
				Filename: utils.P("Screenshot"),
				MimeType: utils.P("png"),
			},
		},
	},
	{
		ID:          2,
		Description: "Small flowerbed in the park",
		RegionID:    3,
		SensorID:    -1, // no sensor
		Images:      []*entities.Image{},
	},
}
