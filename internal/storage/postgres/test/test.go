package test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
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

func SetupPostgresContainer() (func(), *sql.DB, error) {
	slog.Info("Setting up postgres container")
	ctx := context.Background()

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
			WithStartupTimeout(time.Second * 30).
			WithPollInterval(time.Microsecond * 100).
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

	dbUrl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbName)
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		slog.Error("Error connecting to PostgreSQL", "error", err)
		return closeFunc, db, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("Error pinging PostgreSQL", "error", err)
		return closeFunc, db, err
	}

	execMigration(db)

	return closeFunc, db, nil
}

func execMigration(db *sql.DB) error {
	slog.Info("Executing migration")
	// Execute migration with make migrate/up

  rootDir := utils.RootDir()
	migrationPath := fmt.Sprintf("%s/internal/storage/postgres/migrations/", rootDir)
	seedPath := fmt.Sprintf("%s/internal/storage/postgres/test/seed/", rootDir)

	if err := goose.Up(db, migrationPath); err != nil {
		slog.Error("Error executing migration", "error", err)
		return err
	}

	if err := goose.Up(db, seedPath, goose.WithNoVersioning()); err != nil {
		slog.Error("Error executing seed", "error", err)
		return err
	}

	return nil
}

func dbUrl() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbName)
}
