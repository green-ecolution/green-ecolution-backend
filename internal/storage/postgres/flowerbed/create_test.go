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

	input := entities.Flowerbed{
		Description: "New description",
		Size: 40.000,
		NumberOfPlants: 100.000,
		MoistureLevel: 3.5,
		Address: "126 Water street",
		Latitude: utils.P(54.81269326939148), 
		Longitude: utils.P(54.81269326939148),
		Region: region,
	}

	t.Run("should create flowerbed with description, region, longitude and latitude", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Create(
			context.Background(), 
			WithDescription(input.Description),
			WithLatitude(input.Latitude),
			WithRegion(input.Region),
			WithLongitude(input.Longitude),
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
		assert.Empty(t, got.Sensor)
		assert.Empty(t, got.Images)
		assert.Equal(t, "", got.Address)
		assert.Equal(t, 0.0, got.MoistureLevel)
		assert.Equal(t, int32(0), got.NumberOfPlants)
		assert.Equal(t, 0,0, got.Size)
		assert.False(t, got.Archived)

		// assert region
		assert.Equal(t, region.ID, got.Region.ID)
		assert.Equal(t, region.Name, got.Region.Name)
	})
}
