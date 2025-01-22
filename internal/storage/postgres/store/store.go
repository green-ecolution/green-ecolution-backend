package store

import (
	"context"
	"errors"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	sqlc.Querier
	db *pgxpool.Pool
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

// TODO: Improve error handling
func (s *Store) HandleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return storage.ErrEntityNotFound
	}

	return err
}

func (s *Store) WithTx(ctx context.Context, fn func(*Store) error) error {
	log := logger.GetLogger(ctx)
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
		log.Debug("committing transaction")
		return tx.Commit(ctx)
	}

	log.Debug("rolling back transaction")
	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		log.Error("error rolling back transaction", "error", rollbackErr)
		return errors.Join(err, rollbackErr)
	}

	return err
}

func (s *Store) Close() {
	s.db.Close()
}
