package store

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	sqlc.Querier
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool, querier sqlc.Querier) *Store {
	if db == nil {
		slog.Error("Database connection is nil")
		panic("Database connection is nil")
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

// ErrorType defines the classification of errors.
type ErrorType string

const (
	DatabaseError   ErrorType = "DatabaseError"
	NotFoundError   ErrorType = "NotFoundError"
	UNexpectedError ErrorType = "UnexpectedError"
)

// ClassifiedError represents an error with additional contextual information.
type ClassifiedError struct {
	Type      ErrorType
	Massage   string
	Original  error
	File      string
	Line      int
	Timestamp string
	Code      string
}

func (e *ClassifiedError) Error() string {
	return fmt.Sprintf("%s %s (at %s:%d)%s", e.Type, e.Massage, e.File, e.Line, e.Code)
}

// HandleError classifies and logs errors with additional context.
func (s *Store) HandleError(err error, contextMs ...string) error {
	if err == nil {
		return nil
	}

	// Get the context message
	contextMsg := "No context provided"
	if len(contextMs) > 0 {
		contextMsg = contextMs[0]
	}

	// Capture file and line number
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	} else {
		baseMarker := "internal/"
		if idx := strings.Index(file, baseMarker); idx != -1 {
			file = file[idx:]
		}
	}

	// Generate the timestamp for when the error occurred.
	timestamp := time.Now().Format(time.RFC3339)
	var classifiedError *ClassifiedError

	// Handle specific error types.
	if errors.Is(err, pgx.ErrNoRows) {
		// NotFoundError
		classifiedError = &ClassifiedError{
			Type:      NotFoundError,
			Massage:   contextMsg,
			Original:  storage.ErrEntityNotFound,
			File:      file,
			Line:      line,
			Timestamp: timestamp,
		}
	} else if pgErr, ok := err.(*pgconn.PgError); ok {
		// DatabaseError
		classifiedError = &ClassifiedError{
			Type:     DatabaseError,
			Massage:  contextMsg,
			Original: pgErr,
			File:     file,
			Line:     line,
			Code:     pgErr.Code,
		}
	} else {
		// UNexpectedError
		classifiedError = &ClassifiedError{
			Type:      UNexpectedError,
			Massage:   contextMsg,
			Original:  err,
			File:      file,
			Line:      line,
			Timestamp: timestamp,
		}
	}
	return classifiedError
}

func (s *Store) WithTx(ctx context.Context, fn func(*Store) error) error {
	if fn == nil {
		slog.Error("Transaction function is nil")
		return errors.New("transaction function is nil")
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Failed to start transaction", "error", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	qtx := sqlc.New(tx)
	err = fn(NewStore(s.db, qtx))
	if err == nil {
		slog.Info("Transaction cmmitted successfully")
		slog.Debug("Committing transaction")
		return tx.Commit(ctx)
	}

	slog.Debug("Rolling back transaction")
	slog.Warn("Transaction rollback initiated due to error", "error", err)
	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil {
		slog.Error("Transaction rollback failed", "rollbackError", rollbackErr, "originalError", err)
		return errors.Join(err, rollbackErr)
	}

	return fmt.Errorf("transaction failed: %w", err)
}

func (s *Store) Close() {
	slog.Info("closing database connection")
	s.db.Close()
}
