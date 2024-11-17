package store_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pgContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	code := 1
	ctx := context.Background()
	defer func() { os.Exit(code) }()
	pgContainer = testutils.SetupPostgres(ctx)
	defer pgContainer.Terminate(ctx)

	code = m.Run()
}

func TestStore_NewStore(t *testing.T) {
	t.Run("should create new store", func(t *testing.T) {
		// given
		pool := poolConn(t)

		// when
		s, err := store.NewStore(pool, sqlc.New(pool))

		// then
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})

	t.Run("should return error when querier is nil", func(t *testing.T) {
		// given
		pool := poolConn(t)

		// when
		s, err := store.NewStore(pool, nil)

		// then
		assert.Error(t, err)
		assert.Nil(t, s)
	})

	t.Run("should return error when pool is nil", func(t *testing.T) {
		// given
		var pool *pgxpool.Pool

		// when
		s, err := store.NewStore(pool, sqlc.New(pool))

		// then
		assert.Error(t, err)
		assert.Nil(t, s)
	})

	t.Run("should return error when pool is nil and querier is nil", func(t *testing.T) {
		// when
		s, err := store.NewStore(nil, nil)

		// then
		assert.Error(t, err)
		assert.Nil(t, s)
	})
}

func TestStore_DB(t *testing.T) {
	t.Run("should return db", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s, _ := store.NewStore(pool, sqlc.New(pool))

		// when
		db := s.DB()

		// then
		assert.NotNil(t, db)
	})

	t.Run("should return nil when db is nil", func(t *testing.T) {
		// given
		s := &store.Store{}

		// when
		db := s.DB()

		// then
		assert.Nil(t, db)
	})
}

func TestStore_WithTx(t *testing.T) {
	execMigration(t)
	t.Run("should execute function with transaction", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s, _ := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), func(q *sqlc.Queries) error {
			return nil
		})

		// then
		assert.NoError(t, err)
	})

	t.Run("should return error when function returns error", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s, _ := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), func(q *sqlc.Queries) error {
			return pgx.ErrNoRows
		})

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when function is nil", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s, _ := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), nil)

		// then
		assert.Error(t, err)
	})

	t.Run("should return error when context is canceled", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s, _ := store.NewStore(pool, sqlc.New(pool))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		err := s.WithTx(ctx, func(q *sqlc.Queries) error {
			return nil
		})

		// then
		assert.Error(t, err)
	})

	t.Run("should commit transaction", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s, _ := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), func(q *sqlc.Queries) error {
			q.CreateSensor(context.Background(), sqlc.SensorStatusOnline)
			return nil
		})

		// then
		assert.NoError(t, err)

		// validate
		got, err := s.GetSensorByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, int32(1), got.ID)
		assert.Equal(t, sqlc.SensorStatusOnline, got.Status)

		// cleanup
		_ = s.DeleteSensor(context.Background(), 1)
	})

	t.Run("should rollback transaction", func(t *testing.T) {
		// given
		pool := poolConn(t)
		s, _ := store.NewStore(pool, sqlc.New(pool))

		// when
		err := s.WithTx(context.Background(), func(q *sqlc.Queries) error {
			q.CreateSensor(context.Background(), sqlc.SensorStatusOnline)
			return assert.AnError
		})

		// then
		assert.Error(t, err)

		// validate
		got, _ := s.GetSensorByID(context.Background(), 1)
		assert.Empty(t, got)
	})
}

func execMigration(t testing.TB) {
	dbURL := dbURL(t)
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		t.Fatalf("Error while connecting to PostgreSQL: %s", err)
	}

	testutils.ExecMigration(db)
	defer db.Close()
}

func dbURL(t testing.TB) string {
	ctx := context.Background()
	dbURL, err := pgContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Could not get connection string: %s", err)
	}

	return dbURL
}

func poolConn(t testing.TB) *pgxpool.Pool {
	dbURL := dbURL(t)

	pgxConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		t.Fatalf("Error while parsing PostgreSQL connection string: %s", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		t.Fatalf("Error while connecting to PostgreSQL: %s", err)
	}

	return pool
}
