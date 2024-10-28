package region

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultRegion() *entities.Region {
	return &entities.Region{
		Name: "",
	}
}

func (r *RegionRepository) Create(ctx context.Context, vFn ...entities.EntityFunc[entities.Region]) (*entities.Region, error) {
	entity := defaultRegion()
	for _, fn := range vFn {
		fn(entity)
	}

	if entity.Name == "" {
		return nil, errors.New("name is required")
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	entity.ID = *id
	return r.GetByID(ctx, *id)
}

func (r *RegionRepository) createEntity(ctx context.Context, entity *entities.Region) (*int32, error) {
	args := sqlc.CreateRegionParams{
		Name: entity.Name,
	}

	id, err := r.store.CreateRegion(ctx, &args)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return &id, nil
}
