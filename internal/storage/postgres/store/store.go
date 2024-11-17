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
)

type Store struct {
	*sqlc.Queries
	db         *pgxpool.Pool
	entityType EntityType
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: sqlc.New(db),
		db:      db,
	}
}

func (s *Store) SetEntityType(entityType EntityType) {
	s.entityType = entityType
}

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

func (s *Store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
  if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit(ctx)
	}

	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) CheckSensorExists(ctx context.Context, sensorID *int32) error {
	if sensorID != nil {
		_, err := s.GetSensorByID(ctx, *sensorID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrSensorNotFound
			} else {
				slog.Error("Error getting sensor by id", "error", err)
				return s.HandleError(err)
			}
		}
	}

	return nil
}

func (s *Store) CheckImageExists(ctx context.Context, imageID *int32) error {
	if imageID != nil {
		_, err := s.GetImageByID(ctx, *imageID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return storage.ErrImageNotFound
			} else {
				slog.Error("Error getting image by id", "error", err)
				return s.HandleError(err)
			}
		}
	}

	return nil
}
