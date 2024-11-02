package mapper_test

import (
	"testing"
	"time"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestVehicleMapper_FromSql(t *testing.T) {
	verhicleMapper := &generated.InternalVehicleRepoMapperImpl{}

	t.Run("should convert from sql to entity", func(t *testing.T) {
		// given
		src := allTestVehicles[0]

		// when
		got := verhicleMapper.FromSql(src)

		// then
		assert.NotNil(t, got)
		assert.Equal(t, src.ID, got.ID)
		assert.Equal(t, src.CreatedAt.Time, got.CreatedAt)
		assert.Equal(t, src.UpdatedAt.Time, got.UpdatedAt)
		assert.Equal(t, src.NumberPlate, got.NumberPlate)
		assert.Equal(t, src.Description, got.Description)
		assert.Equal(t, src.WaterCapacity, got.WaterCapacity)
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src *sqlc.Vehicle = nil

		// when
		got := verhicleMapper.FromSql(src)

		// then
		assert.Nil(t, got)
	})
}

func TestVehicleMapper_FromSqlList(t *testing.T) {
	verhicleMapper := &generated.InternalVehicleRepoMapperImpl{}

	t.Run("should convert from sql slice to entity slice", func(t *testing.T) {
		// given
		src := allTestVehicles

		// when
		got := verhicleMapper.FromSqlList(src)

		// then
		assert.NotNil(t, got)
		assert.Len(t, got, 2)

		for i, src := range src {
			assert.NotNil(t, got)
			assert.Equal(t, src.ID, got[i].ID)
			assert.Equal(t, src.CreatedAt.Time, got[i].CreatedAt)
			assert.Equal(t, src.UpdatedAt.Time, got[i].UpdatedAt)
			assert.Equal(t, src.NumberPlate, got[i].NumberPlate)
			assert.Equal(t, src.Description, got[i].Description)
			assert.Equal(t, src.WaterCapacity, got[i].WaterCapacity)
		}
 })

 t.Run("should return nil for nil input", func(t *testing.T) {
     // given
     var src []*sqlc.Vehicle = nil

     // when
     got := verhicleMapper.FromSqlList(src)

     // then
     assert.Nil(t, got)
 })
}

var allTestVehicles = []*sqlc.Vehicle{
 {
	ID:             1,
	CreatedAt:      pgtype.Timestamp{Time: time.Now()},
	UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
	NumberPlate: "FL TZ 1234",
	Description: "This is a big car",
	WaterCapacity: 2000.10,
 },
 {
	ID:             2,
	CreatedAt:      pgtype.Timestamp{Time: time.Now()},
	UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
	NumberPlate: "FL TZ 1235",
	Description: "This is a small car",
	WaterCapacity: 1000,
 },
}
