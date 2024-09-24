package region

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	store "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type RegionRepository struct {
	store *store.Store
	RegionMappers
}

type RegionMappers struct {
	mapper mapper.InternalRegionRepoMapper
}

func NewRegionMappers(rMapper mapper.InternalRegionRepoMapper) RegionMappers {
	return RegionMappers{
		mapper: rMapper,
	}
}

func NewRegionRepository(s *store.Store, mappers RegionMappers) storage.RegionRepository {
	s.SetEntityType(store.Vehicle)
	return &RegionRepository{
		store:         s,
		RegionMappers: mappers,
	}
}

func WithName(name string) entities.EntityFunc[entities.Region] {
	return func(v *entities.Region) {
    v.Name = name
	}
}

func (r *RegionRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteRegion(ctx, id)
}
