package postgreSQL

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgreSQL/Schema"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgreSQL/sensorSQL"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgreSQL/treeSQL"
	_ "github.com/lib/pq"
	"log"
)

func NewPostgresDB(ctx context.Context, cfg config.DatabaseConfig) (*storage.Repository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("cannot create client: %w", err)
	}

	log.Println("Trying to connect to PostgreSQL...")

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	// Ping the database to check if it's available
	if err := db.PingContext(ctx); err != nil {
		log.Println(err)
		return nil, fmt.Errorf("cannot ping client: %w", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	sensorRepo := sensorSQL.NewSensorRepository(db)
	treeRepo := treeSQL.NewTreeRepository(db)
	schemaRepo := Schema.NewSchemaRepository(db)

	return &storage.Repository{
		Sensor: sensorRepo,
		Tree:   treeRepo,
		Schema: schemaRepo,
	}, nil
}
