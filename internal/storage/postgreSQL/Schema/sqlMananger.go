package Schema

import (
	"database/sql"
	"log"
)

type SchemaRepository struct {
	db *sql.DB
}

func NewSchemaRepository(db *sql.DB) *SchemaRepository {
	return &SchemaRepository{db: db}
}
func (r *SchemaRepository) SetupSensorTable() {
	_, err := r.db.Exec(`
        CREATE TABLE IF NOT EXISTS sensors (
            id UUID PRIMARY KEY,
            tree_id TEXT,
            data JSONB
        );
        CREATE TABLE IF NOT EXISTS trees (
            id SERIAL PRIMARY KEY,
            species TEXT,
            tree_num INT,
            age INT,
            latitude FLOAT,
            longitude FLOAT,
            address TEXT,
            additional_info TEXT
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func (r *SchemaRepository) SetupTreeTable() {

	_, err := r.db.Exec(`
	CREATE TABLE IF NOT EXISTS trees (
		id SERIAL PRIMARY KEY,
		species TEXT,
		tree_num INT,
		age INT,
		latitude FLOAT,
		longitude FLOAT,
		address TEXT,
		additional_info TEXT
	);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

}

func (r *SchemaRepository) TeardownDatabase() {
	_, err := r.db.Exec(`
        DROP TABLE IF EXISTS sensors;
        DROP TABLE IF EXISTS trees;
    `)
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}
}
