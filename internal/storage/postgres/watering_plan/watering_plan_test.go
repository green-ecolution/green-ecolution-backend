package watering_plan

import (
	"context"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
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
	)
	defer suite.Terminate(ctx)

	code = m.Run()
}
