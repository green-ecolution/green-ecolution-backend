package mapper_test

import (
	"testing"
	"time"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestFlowerbedMapper_FromSql(t *testing.T) {
	flowerbedMapper := &generated.InternalFlowerbedRepoMapperImpl{}

	t.Run("should convert from sql to entity", func(t *testing.T) {
		// given
		src := allTestFlowerbeds[0]

		// when
		got := flowerbedMapper.FromSql(src)

		// then
		assert.NotNil(t, got)
		assert.Equal(t, src.ID, got.ID)
		assert.Equal(t, src.CreatedAt.Time, got.CreatedAt)
		assert.Equal(t, src.Size, got.Size)
		assert.Equal(t, src.Description, got.Description)
		assert.Equal(t, src.NumberOfPlants, got.NumberOfPlants)
		assert.Equal(t, src.MoistureLevel, got.MoistureLevel)
		assert.Equal(t, src.Address, got.Address)
		assert.Equal(t, src.Archived, got.Archived)
		assert.Equal(t, src.Latitude, got.Latitude)
		assert.Equal(t, src.Longitude, got.Longitude)
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src *sqlc.Flowerbed = nil

		// when
		got := flowerbedMapper.FromSql(src)

		// then
		assert.Nil(t, got)
	})
}

func TestFlowerbedMapper_FromSqlList(t *testing.T) {
	flowerbedMapper := &generated.InternalFlowerbedRepoMapperImpl{}

	t.Run("should convert from sql slice to entity slice", func(t *testing.T) {
		// given
		src := allTestFlowerbeds

		// when
		got := flowerbedMapper.FromSqlList(src)

		// then
		assert.NotNil(t, got)
		assert.Len(t, got, 2)

		for i, src := range src {
			assert.Equal(t, src.ID, got[i].ID)
			assert.Equal(t, src.CreatedAt.Time, got[i].CreatedAt)
			assert.Equal(t, src.UpdatedAt.Time, got[i].UpdatedAt)
			assert.Equal(t, src.Size, got[i].Size)
			assert.Equal(t, src.Description, got[i].Description)
			assert.Equal(t, src.NumberOfPlants, got[i].NumberOfPlants)
			assert.Equal(t, src.MoistureLevel, got[i].MoistureLevel)
			assert.Equal(t, src.Address, got[i].Address)
			assert.Equal(t, src.Archived, got[i].Archived)
			assert.Equal(t, src.Latitude, got[i].Latitude)
			assert.Equal(t, src.Longitude, got[i].Longitude)
		}
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src []*sqlc.Flowerbed = nil

		// when
		got := flowerbedMapper.FromSqlList(src)

		// then
		assert.Nil(t, got)
	})
}

var allTestFlowerbeds = []*sqlc.Flowerbed{
	{
		ID:             1,
		CreatedAt:      pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
		Size:           10.5,
		Description:    "A beautiful flowerbed",
		NumberOfPlants: 20,
		MoistureLevel:  30.0,
		Address:        "123 Garden Lane",
		Archived:       false,
		Latitude:       40.7128,
		Longitude:      -74.0060,
	},
	{
		ID:             2,
		CreatedAt:      pgtype.Timestamp{Time: time.Now().Add(1 * time.Hour)},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now().Add(1 * time.Hour)},
		Size:           15.0,
		Description:    "Another flowerbed",
		NumberOfPlants: 10,
		MoistureLevel:  20.0,
		Address:        "456 Garden Avenue",
		Archived:       true,
		Latitude:       41.1234,
		Longitude:      -73.9876,
	},
}