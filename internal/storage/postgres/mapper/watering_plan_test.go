package mapper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestWateringPlanMapper_FromSql(t *testing.T) {
	wateringPlanMapper := &generated.InternalWateringPlanRepoMapperImpl{}

	t.Run("should convert from sql to entity", func(t *testing.T) {
		// given
		src := allTestWateringPlans[0]

		// when
		got := wateringPlanMapper.FromSql(src)

		// then
		assert.NotNil(t, got)
		assert.Equal(t, src.ID, got.ID)
		assert.Equal(t, src.CreatedAt.Time, got.CreatedAt)
		assert.Equal(t, src.UpdatedAt.Time, got.UpdatedAt)
		assert.Equal(t, src.Date.Time, got.Date)
		assert.Equal(t, src.Description, got.Description)
		assert.Equal(t, src.Distance, got.Distance)
		assert.Equal(t, src.TotalWaterRequired, got.TotalWaterRequired)
		assert.Equal(t, src.WateringPlanStatus, sqlc.WateringPlanStatus(got.WateringPlanStatus))
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src *sqlc.WateringPlan = nil

		// when
		got := wateringPlanMapper.FromSql(src)

		// then
		assert.Nil(t, got)
	})
}

func TestWateringPlanMapper_FromSqlList(t *testing.T) {
	wateringPlanMapper := &generated.InternalWateringPlanRepoMapperImpl{}

	t.Run("should convert from sql slice to entity slice", func(t *testing.T) {
		// given
		src := allTestWateringPlans

		// when
		got := wateringPlanMapper.FromSqlList(src)

		// then
		assert.NotNil(t, got)
		assert.Len(t, got, 5)

		for i, src := range src {
			assert.NotNil(t, got)
			assert.NotNil(t, got)
			assert.Equal(t, src.ID, got[i].ID)
			assert.Equal(t, src.CreatedAt.Time, got[i].CreatedAt)
			assert.Equal(t, src.UpdatedAt.Time, got[i].UpdatedAt)
			assert.Equal(t, src.Date.Time, got[i].Date)
			assert.Equal(t, src.Description, got[i].Description)
			assert.Equal(t, src.Distance, got[i].Distance)
			assert.Equal(t, src.TotalWaterRequired, got[i].TotalWaterRequired)
			assert.Equal(t, src.WateringPlanStatus, sqlc.WateringPlanStatus(got[i].WateringPlanStatus))
		}
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src []*sqlc.WateringPlan = nil

		// when
		got := wateringPlanMapper.FromSqlList(src)

		// then
		assert.Nil(t, got)
	})
}

var allTestWateringPlans = []*sqlc.WateringPlan{
	{
		ID: 1,
		Date: pgtype.Date{
			Time:  time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Description:        "New watering plan for the west side of the city",
		WateringPlanStatus: "planned",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID: 2,
		Date: pgtype.Date{
			Time:  time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Description:        "New watering plan for the east side of the city",
		WateringPlanStatus: "active",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID: 3,
		Date: pgtype.Date{
			Time:  time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Description:        "Very important watering plan due to no rainfall",
		WateringPlanStatus: "finished",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID: 4,
		Date: pgtype.Date{
			Time:  time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Description:        "New watering plan for the south side of the city",
		WateringPlanStatus: "not competed",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
	{
		ID: 5,
		Date: pgtype.Date{
			Time:  time.Date(2024, 9, 22, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Description:        "Canceled due to flood",
		WateringPlanStatus: "canceled",
		Distance:           utils.P(63.0),
		TotalWaterRequired: utils.P(6000.0),
	},
}

func TestMapWateringPlanStatus(t *testing.T) {
	tests := []struct {
		input    sqlc.WateringPlanStatus
		expected entities.WateringPlanStatus
	}{
		{input: sqlc.WateringPlanStatusActive, expected: entities.WateringPlanStatusActive},
		{input: sqlc.WateringPlanStatusCanceled, expected: entities.WateringPlanStatusCanceled},
		{input: sqlc.WateringPlanStatusFinished, expected: entities.WateringPlanStatusFinished},
		{input: sqlc.WateringPlanStatusNotcompeted, expected: entities.WateringPlanStatusNotCompeted},
		{input: sqlc.WateringPlanStatusPlanned, expected: entities.WateringPlanStatusPlanned},
		{input: sqlc.WateringPlanStatusUnknown, expected: entities.WateringPlanStatusUnknown},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("should return %v for input %v", test.expected, test.input), func(t *testing.T) {
			result := mapper.MapWateringPlanStatus(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}
