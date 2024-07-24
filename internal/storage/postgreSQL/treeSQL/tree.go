package treeSQL

import (
	"context"
	"database/sql"
	"errors"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgreSQL/entities/tree"
)

type TreeRepository struct {
	db *sql.DB
}

func NewTreeRepository(db *sql.DB) *TreeRepository {
	return &TreeRepository{db: db}
}
func (r *TreeRepository) Insert(ctx context.Context, data *tree.TreeEntity) error {
	query := `
		INSERT INTO trees (species, tree_num, age, latitude, longitude, address, additional_info)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query, data.Species, data.TreeNum, data.Age, data.Location.Latitude, data.Location.Longitude, data.Location.Address, data.Location.AdditionalInfo)
	if err != nil {
		return storage.ErrCannotUpsertData
	}

	return nil
}
func (r *TreeRepository) Get(ctx context.Context, id string) (*tree.TreeEntity, error) {
	var data *tree.TreeEntity
	query := `
		SELECT id, species, tree_num, age, latitude, longitude, address, additional_info
		FROM trees
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&data.ID, &data.Species, &data.TreeNum, &data.Age,
		&data.Location.Latitude, &data.Location.Longitude,
		&data.Location.Address, &data.Location.AdditionalInfo,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrDataNotFound
		}
		return nil, err
	}

	return data, nil
}

func (r *TreeRepository) GetAll(ctx context.Context) ([]*tree.TreeEntity, error) {
	query := `
		SELECT id, species, tree_num, age, latitude, longitude, address, additional_info
		FROM trees
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, storage.ErrDataNotFound
	}
	defer rows.Close()

	var trees []*tree.TreeEntity

	for rows.Next() {
		entity := new(tree.TreeEntity)
		if err := rows.Scan(
			&entity.ID, &entity.Species, &entity.TreeNum, &entity.Age,
			&entity.Location.Latitude, &entity.Location.Longitude,
			&entity.Location.Address, &entity.Location.AdditionalInfo,
		); err != nil {
			return nil, err
		}

		trees = append(trees, entity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trees, nil
}
