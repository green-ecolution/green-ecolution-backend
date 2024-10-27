package testutils

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/pressly/goose/v3"
)

// ExecMigration executes the migration
func ExecMigration(db *sql.DB) error {
	slog.Info("Executing migration...")

	rootDir := utils.RootDir()
	migrationPath := fmt.Sprintf("%s/internal/storage/postgres/migrations/", rootDir)

	if err := goose.Up(db, migrationPath); err != nil {
		slog.Error("Error executing migration", "error", err)
		return err
	}

	return nil
}

// ResetDB resets the database
func ResetDB(db *sql.DB) error {
	slog.Info("Resetting database...")

	rootDir := utils.RootDir()
	migrationPath := fmt.Sprintf("%s/internal/storage/postgres/migrations/", rootDir)

	if err := goose.Reset(db, migrationPath); err != nil {
		slog.Error("Error executing migration", "error", err)
		return err
	}

	if err := goose.Up(db, migrationPath); err != nil {
		slog.Error("Error executing migration", "error", err)
		return err
	}

	return nil
}
