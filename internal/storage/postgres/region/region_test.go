package region

import (
	"context"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	store "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
)

type regionFields struct {
	store         *store.Store
	RegionMappers RegionMappers
}

var (
	defaultFields regionFields
	suite         *testutils.PostgresTestSuite
)

func defaultRegionMappers() RegionMappers {
	return NewRegionMappers(&generated.InternalRegionRepoMapperImpl{})
}

func TestMain(m *testing.M) {
  code := 1
  ctx := context.Background()
  defer func() { os.Exit(code) }()
  suite = testutils.SetupPostgresTestSuite(ctx)
  defaultFields = regionFields{
    store:         suite.Store,
    RegionMappers: defaultRegionMappers(),
  }
  defer suite.Container.Terminate(ctx)

  code = m.Run()
}

