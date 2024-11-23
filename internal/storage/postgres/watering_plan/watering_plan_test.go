package wateringplan

import (
	"context"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/stretchr/testify/assert"
)

var (
	mappers WateringPlanMappers
	suite   *testutils.PostgresTestSuite
)

func TestMain(m *testing.M) {
	code := 1
	ctx := context.Background()
	defer func() { os.Exit(code) }()
	suite = testutils.SetupPostgresTestSuite(ctx)
	mappers = NewWateringPlanRepositoryMappers(
		&generated.InternalWateringPlanRepoMapperImpl{},
		&generated.InternalVehicleRepoMapperImpl{},
		&generated.InternalTreeClusterRepoMapperImpl{},
	)
	defer suite.Terminate(ctx)

	code = m.Run()
}

func TestWateringPlanRepository_Delete(t *testing.T) {
	suite.ResetDB(t)
	suite.InsertSeed(t, "internal/storage/postgres/seed/test/watering_plan")

	t.Run("should delete watering plan", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 1)
		got, errGot := r.GetByID(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.Error(t, errGot)
		assert.Nil(t, got)
	})

	t.Run("should return error when watering plan not found", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), 99)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when watering plan with negative id", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)

		// when
		err := r.Delete(context.Background(), -1)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when vehicle is canceled", func(t *testing.T) {
		// given
		r := NewWateringPlanRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := r.Delete(ctx, 1)

		// then
		assert.Error(t, err)
	})

	// TODO: Add test cases to check if the correct entities in the pivot tables are deleted
}
