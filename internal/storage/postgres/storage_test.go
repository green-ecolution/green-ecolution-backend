package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
  _ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/twpayne/go-geos"
	pgxgeos "github.com/twpayne/pgx-geos"
)

type PostgresTestSuite struct {
	container *postgres.PostgresContainer
	snapshot  string
	dbURL     string
	store     *store.Store
}

var (
	TestSuite PostgresTestSuite
)

func TestMain(m *testing.M) {
	code := 1
	defer func() { os.Exit(code) }()

	TestSuite = PostgresTestSuite{
		snapshot: "integration-test",
	}

	ctx := context.Background()
	container, err := testutils.StartPostgres(ctx)
	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}
	TestSuite.container = container.Container
	TestSuite.dbURL = container.URL

	// Migrate database
	db, err := sql.Open("pgx", container.URL)
	if err != nil {
		log.Fatalf("Could not connect to postgres: %s", err)
	}

	if err = testutils.ExecMigration(db); err != nil {
		log.Fatalf("Could not execute migration: %s", err)
	}
  err = db.Close()
  if err != nil {
    log.Fatalf("Could not close database connection: %s", err)
  }

	// Take a snapshot
	if err = container.Container.Snapshot(ctx, postgres.WithSnapshotName(TestSuite.snapshot)); err != nil {
		log.Fatalf("Could not create snapshot: %s", err)
	}

	pgxConfig, err := pgxpool.ParseConfig(container.URL)
	if err != nil {
		log.Fatalf("Error while parsing PostgreSQL connection string: %s", err)
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxgeos.Register(ctx, conn, geos.NewContext())
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatalf("Error while connecting to PostgreSQL: %s", err)
	}
	TestSuite.store = store.NewStore(pool)

	code = m.Run()
}

// TakeSnapshot takes a snapshot of the current database state
func TakeSnapshot(t *testing.T, snapshot string) {
	t.Log("Taking snapshot...")
	if err := TestSuite.container.Snapshot(context.Background(), postgres.WithSnapshotName(snapshot)); err != nil {
		t.Fatalf("Could not create snapshot: %s", err)
	}
}

// ResetDB resets the database to the state of the snapshot
func ResetDB(t *testing.T) {
	t.Log("Resetting database...")
	if err := TestSuite.container.Restore(context.Background(), postgres.WithSnapshotName(TestSuite.snapshot)); err != nil {
		t.Fatalf("Could not restore snapshot: %s", err)
	}
}

// ResetDBByName resets the database to the state of the snapshot with the given name
func ResetDBByName(t *testing.T, snapshot string) {
	t.Log("Resetting database...")
	if err := TestSuite.container.Restore(context.Background(), postgres.WithSnapshotName(snapshot)); err != nil {
		t.Fatalf("Could not restore snapshot: %s", err)
	}
}
