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
		for i, tc := range got {
			assert.Equal(t, allTestFlowerbeds[i].ID, tc.ID)
			assert.Equal(t, allTestFlowerbeds[i].Description, tc.Description)

			// assert region
			assert.NotNil(t, tc.Region)
			assert.Equal(t, allTestFlowerbeds[i].RegionID, tc.Region.ID)

			// assert sensor
			if allTestFlowerbeds[i].SensorID == -1 {
				assert.Nil(t, tc.Sensor)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, tc.Sensor)
				assert.Equal(t, allTestFlowerbeds[i].SensorID, tc.Sensor.ID)
			}

			// assert images
			if len(allTestFlowerbeds[i].Images) == 0 {
				assert.Equal(t, []*entities.Image{}, tc.Images)
				assert.NoError(t, err)
			} else {
				assert.NotNil(t, tc.Images)
				assert.Equal(t, allTestFlowerbeds[i].Images[0].ID, tc.Images[0].ID)
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

type testFlowerbeds struct {
	ID       int32
	Description     string
	RegionID int32
	SensorID int32
	Images []*entities.Image
}

var allTestFlowerbeds = []*testFlowerbeds{
	{
		ID:       1,
		Description:     "Big flowerbed nearby the sea",
		RegionID: 1,
		SensorID:  2,
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
		ID:       2,
		Description:     "Small flowerbed in the park",
		RegionID: 3,
		SensorID:  -1, // no sensor
		Images: []*entities.Image{},
	},
}
