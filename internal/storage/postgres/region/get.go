package region

import (
	"context"
	"errors"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/jackc/pgx/v5"
)

func (r *RegionRepository) GetAll(ctx context.Context) ([]*entities.Region, error) {
	rows, err := r.store.GetAllRegions(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *RegionRepository) GetByID(ctx context.Context, id int32) (*entities.Region, error) {
	row, err := r.store.GetRegionById(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

func (r *RegionRepository) GetByName(ctx context.Context, plate string) (*entities.Region, error) {
	row, err := r.store.GetRegionByName(ctx, plate)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

func (r *RegionRepository) GetByPoint(ctx context.Context, latitude, longitude float64) (*entities.Region, error) {
	p := fmt.Sprintf("POINT(%f %f)", longitude, latitude)
	region, err := r.store.GetRegionByPoint(ctx, p)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(region), nil
}
