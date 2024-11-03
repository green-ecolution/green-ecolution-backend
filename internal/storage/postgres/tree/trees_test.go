package tree

import (
	"context"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"os"
	"testing"
)

var (
	mappers TreeMappers
	suite   *testutils.PostgresTestSuite
)

func TestMain(m *testing.M) {
	code := 1
	ctx := context.Background()
	defer func() { os.Exit(code) }()
	suite = testutils.SetupPostgresTestSuite(ctx)
	mappers = NewTreeRepositoryMappers(
		&generated.InternalTreeRepoMapperImpl{},
		&generated.InternalImageRepoMapperImpl{},
		&generated.InternalSensorRepoMapperImpl{},
		&generated.InternalTreeClusterRepoMapperImpl{},
	)
	defer suite.Terminate(ctx)
	code = m.Run()
}
