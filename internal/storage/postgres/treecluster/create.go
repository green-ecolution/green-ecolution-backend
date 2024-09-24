package treecluster

import (
	"context"
	"fmt"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultTreeCluster() *entities.TreeCluster {
	return &entities.TreeCluster{
		Region:         &entities.Region{},
		Address:        "",
		Description:    "",
		MoistureLevel:  0,
		Latitude:       54.7752933631787,
		Longitude:      9.451569031678723,
		WateringStatus: entities.TreeClusterWateringStatusUnknown,
		SoilCondition:  entities.TreeSoilConditionUnknown,
		Archived:       false,
		LastWatered:    nil,
		Trees:          make([]*entities.Tree, 0),
	}
}

func (r *TreeClusterRepository) Create(ctx context.Context, tcFn ...entities.EntityFunc[entities.TreeCluster]) (*entities.TreeCluster, error) {
	entity := defaultTreeCluster()
	for _, fn := range tcFn {
		fn(entity)
	}

	if entity.Region == nil || entity.Region.ID == 0 {
		region, err := r.getRegion(ctx, entity)
		if err != nil {
			return nil, err
		}
		entity.Region = region
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}
	entity.ID = id

	return r.GetByID(ctx, id)
}

func (r *TreeClusterRepository) createEntity(ctx context.Context, entity *entities.TreeCluster) (int32, error) {
	args := sqlc.CreateTreeClusterParams{
		RegionID:       &entity.Region.ID,
		Address:        entity.Address,
		Description:    entity.Description,
		Latitude:       entity.Latitude,
		Longitude:      entity.Longitude,
		MoistureLevel:  entity.MoistureLevel,
		WateringStatus: sqlc.TreeClusterWateringStatus(entity.WateringStatus),
		SoilCondition:  sqlc.TreeSoilCondition(entity.SoilCondition),
	}

	return r.store.CreateTreeCluster(ctx, &args)
}

func (r *TreeClusterRepository) getRegion(ctx context.Context, entity *entities.TreeCluster) (*entities.Region, error) {
	p := fmt.Sprintf("POINT(%f %f)", entity.Latitude, entity.Longitude)
	region, err := r.store.GetRegionByPoint(ctx, p)
	if err != nil {
		return nil, err
	}

	return r.regionMapper.FromSql(region), nil
}
