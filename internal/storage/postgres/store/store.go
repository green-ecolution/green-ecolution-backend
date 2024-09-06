package store

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type EntityType string

const (
	Image  EntityType = "image"
	Sensor EntityType = "sensor"
)

type Store struct {
	*sqlc.Queries
	db         *pgx.Conn
	entityType EntityType
}

func NewStore(db *pgx.Conn, entityType EntityType) *Store {
	return &Store{
		Queries:    sqlc.New(db),
		db:         db,
		entityType: entityType,
	}
}

func (s *Store) HandleError(err error) error {
	if err == nil {
		return nil
	}

	slog.Error("An Error occured in database operation", "error", err, "entityType", s.entityType)
	switch err {
	case pgx.ErrNoRows:
		switch s.entityType {
		case Image:
			slog.Error("Image not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrImageNotFound
		case Sensor:
			slog.Error("Sensor not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrSensorNotFound
		default:
			slog.Error("Entity not found", "error", err, "stack", errors.WithStack(err))
			return storage.ErrEntityNotFound
		}
	case pgx.ErrTooManyRows:
		slog.Error("Recieve more rows then expected", "error", err, "stack", errors.WithStack(err))
		return storage.ErrToManyRows
	case pgx.ErrTxClosed:
		slog.Error("Connection is closed", "error", err, "stack", errors.WithStack(err))
		return storage.ErrTxClosed
	case pgx.ErrTxCommitRollback:
		slog.Error("Transaction cannot commit or rollback", "error", err, "stack", errors.WithStack(err))
		return storage.ErrTxCommitRollback
	case sql.ErrConnDone:
		slog.Error("Connection is closed", "error", err, "stack", errors.WithStack(err))
		return storage.ErrConnectionClosed

	default:
		slog.Error("Unknown error", "error", err, "stack", errors.WithStack(err))
		return errors.Wrap(err, "unknown error in postgres store")
	}

}

func (s *Store) WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := sqlc.New(tx)
	err = fn(q)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) Close() {
	s.db.Close(context.Background())
}
