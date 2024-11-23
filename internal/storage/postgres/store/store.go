package store

import (
	"context"
	"errors"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EntityType string

const (
	Sensor      EntityType = "sensor"
	Image       EntityType = "image"
	Flowerbed   EntityType = "flowerbed"
	TreeCluster EntityType = "treecluster"
	Tree        EntityType = "tree"
	Vehicle     EntityType = "vehicle"
	Region      EntityType = "region"
	WateringPlan EntityType = "watering plan"
)

type Store struct {
	sqlc.Querier
	db         *pgxpool.Pool
	entityType EntityType
}

func NewStore(db *pgxpool.Pool, querier sqlc.Querier) *Store {
	if db == nil {
		slog.Error("db is nil")
		panic("db is nil")
	}

	if querier == nil {
		slog.Error("querier is nil")
		panic("querier is nil")
	}

	return &Store{
		Querier: querier,
		db:      db,
	}
}

func (s *Store) DB() *pgxpool.Pool {
	return s.db
}

// TODO: Remove
func (s *Store) SetEntityType(entityType EntityType) {
	s.entityType = entityType
}

// TODO: Improve error handling
func (s *Store) HandleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrEntityNotFound
	}

	slog.Error("An Error occurred in database operation", "error", err)
	return err
}

func (s *Store) WithTx(ctx context.Context, fn func(*Store) error) error {
	if fn == nil {
		return errors.New("txFn is nil")
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	qtx := sqlc.New(tx)
	err = fn(NewStore(s.db, qtx))
	if err == nil {
		slog.Debug("Committing transaction")
		return tx.Commit(ctx)
	}

	slog.Debug("Rolling back transaction")
	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		slog.Error("Error rolling back transaction", "error", rollbackErr)
		return errors.Join(err, rollbackErr)
	}

	return err
}

func (s *Store) Close() {
	s.db.Close()
}
