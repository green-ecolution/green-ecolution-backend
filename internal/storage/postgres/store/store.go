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
)

type Store struct {
	sqlc.Querier
	db         *pgxpool.Pool
	entityType EntityType
}

func NewStore(db *pgxpool.Pool, querier sqlc.Querier) (*Store, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	if querier == nil {
		return nil, errors.New("querier is nil")
	}

	return &Store{
		Querier: querier,
		db:      db,
	}, nil
}

func (s *Store) SwitchQuerier(querier sqlc.Querier) func() {
  originalQuerier := s.Querier
  s.Querier = querier

  return func() {
    s.Querier = originalQuerier
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
		switch s.entityType {
		case Sensor:
			return storage.ErrSensorNotFound
		case Image:
			return storage.ErrImageNotFound
		case Flowerbed:
			return storage.ErrFlowerbedNotFound
		case Vehicle:
			return storage.ErrVehicleNotFound
		case TreeCluster:
			return storage.ErrTreeClusterNotFound
		case Tree:
			return storage.ErrTreeNotFound
		default:
			return storage.ErrEntityNotFound
		}
	}

	slog.Error("An Error occurred in database operation", "error", err)
	return err
}

func (s *Store) WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
  if fn == nil {
    return errors.New("txFn is nil")
  }

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

  qtx := sqlc.New(tx)
	err = fn(qtx)
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

