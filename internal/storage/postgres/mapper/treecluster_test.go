package mapper_test

import (
	"testing"
	"time"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestTreeclusterMapper_FromSql(t *testing.T) {
	treeclusterMapper := &generated.InternalTreeClusterRepoMapperImpl{}

	t.Run("should convert from sql to entity", func(t *testing.T) {
		// given
		src := allTestTreecluster[0]

		// when
		got := treeclusterMapper.FromSql(src)

		// then
		assert.NotNil(t, got)
		assert.Equal(t, src.ID, got.ID)
		assert.Equal(t, src.CreatedAt.Time, got.CreatedAt)
		assert.Equal(t, src.UpdatedAt.Time, got.UpdatedAt)
		assert.Equal(t, src.WateringStatus, sqlc.WateringStatus(got.WateringStatus))
		assert.Equal(t, src.LastWatered.Time, *got.LastWatered)
		assert.Equal(t, src.MoistureLevel, got.MoistureLevel)
		assert.Equal(t, src.Address, got.Address)
		assert.Equal(t, src.Archived, got.Archived)
		assert.Equal(t, src.Latitude, got.Latitude)
		assert.Equal(t, src.Longitude, got.Longitude)
		assert.Equal(t, src.Name, got.Name)
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src *sqlc.TreeCluster = nil

		// when
		got := treeclusterMapper.FromSql(src)

		// then
		assert.Nil(t, got)
	})
}

func TestTreeclusterMapper_FromSqlList(t *testing.T) {
	treeclusterMapper := &generated.InternalTreeClusterRepoMapperImpl{}

	t.Run("should convert from sql slice to entity slice", func(t *testing.T) {
		// given
		src := allTestTreecluster

		// when
		got := treeclusterMapper.FromSqlList(src)

		// then
		assert.NotNil(t, got)
		assert.Len(t, got, 2)

		for i, src := range src {
			assert.Equal(t, src.ID, got[i].ID)
			assert.Equal(t, src.CreatedAt.Time, got[i].CreatedAt)
			assert.Equal(t, src.UpdatedAt.Time, got[i].UpdatedAt)
			assert.Equal(t, src.WateringStatus, sqlc.WateringStatus(got[i].WateringStatus))
			assert.Equal(t, src.LastWatered.Time, *got[i].LastWatered)
			assert.Equal(t, src.MoistureLevel, got[i].MoistureLevel)
			assert.Equal(t, src.Address, got[i].Address)
			assert.Equal(t, src.Archived, got[i].Archived)
			assert.Equal(t, src.Latitude, got[i].Latitude)
			assert.Equal(t, src.Longitude, got[i].Longitude)
			assert.Equal(t, src.Name, got[i].Name)
		}
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src []*sqlc.TreeCluster = nil

		// when
		got := treeclusterMapper.FromSqlList(src)

		// then
		assert.Nil(t, got)
	})
}

var allTestTreecluster = []*sqlc.TreeCluster{
	{
		ID:             1,
		CreatedAt:      pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
		WateringStatus: sqlc.WateringStatusGood,
		LastWatered:    pgtype.Timestamp{Time: time.Now()},
		MoistureLevel:  4.10,
		Address:        "123 Garden Lane",
		Description:    "Cluster with newly planted trees",
		Archived:       false,
		SoilCondition:  sqlc.TreeSoilConditionSandig,
		Latitude:       utils.Float64Ptr(41.1234),
		Longitude:      utils.Float64Ptr(-73.9876),
		Name:           "Treecluster 01",
	},
	{
		ID:             2,
		CreatedAt:      pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
		WateringStatus: sqlc.WateringStatusGood,
		LastWatered:    pgtype.Timestamp{Time: time.Now()},
		MoistureLevel:  4.10,
		Address:        "345 Garden Lane",
		Description:    "Cluster needs a lot of care",
		Archived:       false,
		SoilCondition:  sqlc.TreeSoilConditionTonig,
		Latitude:       utils.Float64Ptr(41.1234),
		Longitude:      utils.Float64Ptr(-73.9876),
		Name:           "Treecluster 02",
	},
}
