package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbUsername = "postgres"
	dbPassword = "postgres"
	dbName     = "postgres"
)

var (
	dbHost string
	dbPort int
)

func SetupPostgresContainer() (shutdown func(), url *string, err error) {
	slog.Info("Setting up postgres container")
	ctx := context.Background()

	startupTimeout := time.Second * 30
	pollInterval := time.Microsecond * 100

	req := testcontainers.ContainerRequest{
		Image:        "postgis/postgis",
		ExposedPorts: []string{"5432/tcp", "55432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     dbUsername,
			"POSTGRES_PASSWORD": dbPassword,
		},
		ShmSize:    128 * 1024 * 1024,
		AutoRemove: true,
		Cmd:        []string{"postgres", "-c", "fsync=off"},
		WaitingFor: wait.ForSQL(nat.Port("5432/tcp"), "pgx", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUsername, dbPassword, host, port.Port(), dbName)
		}).
			WithStartupTimeout(startupTimeout).
			WithPollInterval(pollInterval).
			WithQuery("SELECT 1"),
	}

	psqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		slog.Error("Error creating container", "error", err)
		panic(err)
	}

	closeFunc := func() {
		slog.Info("Closing container")
		if err = psqlC.Terminate(ctx); err != nil {
			slog.Error("Error terminating container", "error", err)
			panic(err)
		}
	}

	dbHost, err = psqlC.Host(ctx)
	if err != nil {
		slog.Error("Error getting host", "error", err)
		closeFunc()
		return nil, nil, err
	}

	p, err := psqlC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		slog.Error("Error getting port", "error", err)
		closeFunc()
		panic(err)
	}
	dbPort = p.Int()

	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbName)
	pgxConn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return closeFunc, &dbURL, err
	}

	if err := pgxConn.Ping(ctx); err != nil {
		slog.Error("Error pinging PostgreSQL", "error", err)
		return closeFunc, &dbURL, err
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		slog.Error("Error connecting to PostgreSQL", "error", err)
		return closeFunc, &dbURL, err
	}

	if err = execMigration(db); err != nil {
		slog.Error("Error executing migration", "error", err)
		panic(err)
	}

	return closeFunc, &dbURL, nil
}

func execMigration(db *sql.DB) error {
	slog.Info("Executing migration")

	rootDir := utils.RootDir()
	migrationPath := fmt.Sprintf("%s/internal/storage/postgres/migrations/", rootDir)

	if err := goose.Up(db, migrationPath); err != nil {
		slog.Error("Error executing migration", "error", err)
		return err
	}

	return nil
}

// WithTx Run tests with a transaction. This function will rollback the transaction after the test is done.
func WithTx(_ *testing.T, fn func(db *pgx.Conn)) {
	ctx := context.Background()
	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbName)
	db, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		slog.Error("Error connecting to PostgreSQL", "error", err)
		panic(err)
	}

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		slog.Error("Error connecting to PostgreSQL", "error", err)
		panic(err)
	}

	defer func(db *pgx.Conn, ctx context.Context) {
		err := db.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(db, ctx)
	fn(db)

	if !db.IsClosed() {
		if err := tx.Rollback(ctx); err != nil {
			slog.Error("Error rolling back transaction", "error", err)
			panic(err)
		}
	}
}
