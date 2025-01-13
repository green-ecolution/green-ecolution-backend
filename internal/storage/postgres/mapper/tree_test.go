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

func TestTreeMapper_FromSql(t *testing.T) {
	treeMapper := &generated.InternalTreeRepoMapperImpl{}

	t.Run("should convert from sql to entity", func(t *testing.T) {
		// given
		src := allTestTrees[0]

		// when
		got := treeMapper.FromSql(src)

		// then
		assert.NotNil(t, got)
		assert.Equal(t, src.ID, got.ID)
		assert.Equal(t, src.CreatedAt.Time, got.CreatedAt)
		assert.Equal(t, src.UpdatedAt.Time, got.UpdatedAt)
		assert.Equal(t, src.PlantingYear, got.PlantingYear)
		assert.Equal(t, src.Species, got.Species)
		assert.Equal(t, src.Number, got.Number)
		assert.Equal(t, src.Latitude, got.Latitude)
		assert.Equal(t, src.Longitude, got.Longitude)
		assert.Equal(t, src.WateringStatus, sqlc.WateringStatus(got.WateringStatus))
		assert.Equal(t, src.Readonly, got.Readonly)
		assert.Equal(t, *src.Description, got.Description)
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src *sqlc.Tree = nil

		// when
		got := treeMapper.FromSql(src)

		// then
		assert.Nil(t, got)
	})
}

func TestTreeMapper_FromSqlList(t *testing.T) {
	treeMapper := &generated.InternalTreeRepoMapperImpl{}

	t.Run("should convert from sql slice to entity slice", func(t *testing.T) {
		// given
		src := allTestTrees

		// when
		got := treeMapper.FromSqlList(src)

		// then
		assert.NotNil(t, got)
		assert.Len(t, got, 2)

		for i, src := range src {
			assert.NotNil(t, got)
			assert.Equal(t, src.ID, got[i].ID)
			assert.Equal(t, src.CreatedAt.Time, got[i].CreatedAt)
			assert.Equal(t, src.UpdatedAt.Time, got[i].UpdatedAt)
			assert.Equal(t, src.PlantingYear, got[i].PlantingYear)
			assert.Equal(t, src.Species, got[i].Species)
			assert.Equal(t, src.Number, got[i].Number)
			assert.Equal(t, src.Latitude, got[i].Latitude)
			assert.Equal(t, src.Longitude, got[i].Longitude)
			assert.Equal(t, src.WateringStatus, sqlc.WateringStatus(got[i].WateringStatus))
			assert.Equal(t, src.Readonly, got[i].Readonly)
			assert.Equal(t, *src.Description, got[i].Description)
		}
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src []*sqlc.Tree = nil

		// when
		got := treeMapper.FromSqlList(src)

		// then
		assert.Nil(t, got)
	})
}

var allTestTrees = []*sqlc.Tree{
	{
		ID:             1,
		CreatedAt:      pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
		PlantingYear:   2024,
		Species:        "Oak",
		Latitude:       52.5200,
		Longitude:      13.4050,
		WateringStatus: sqlc.WateringStatusGood,
		Readonly:       true,
		Description:    utils.P("Newly planted tree"),
		Number:         "P 1234",
	},
	{
		ID:             2,
		CreatedAt:      pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
		PlantingYear:   2024,
		Species:        "Maple",
		Latitude:       52.5200,
		Longitude:      13.4050,
		WateringStatus: sqlc.WateringStatusModerate,
		Readonly:       true,
		Description:    utils.P("Also newly planted tree"),
		Number:         "P 2345",
	},
}
