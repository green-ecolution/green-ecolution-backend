package treecluster

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
)

func defaultTreeCluster() *entities.TreeCluster {
	return &entities.TreeCluster{
		Region:         "",
		Address:        "",
		Description:    "",
		MoistureLevel:  0,
		Latitude:       54.7752933631787,
		Longitude:      9.451569031678723,
		WateringStatus: entities.TreeClusterWateringStatusUnknown,
		SoilCondition:  entities.TreeSoilConditionUnknown,
		Archived:       false,
		LastWatered:    nil,
	}
}

func (r *TreeClusterRepository) Create(ctx context.Context, tcFn ...entities.EntityFunc[entities.TreeCluster]) (*entities.TreeCluster, error) {
	entity := defaultTreeCluster()
	for _, fn := range tcFn {
		fn(entity)
	}

	id, err := r.createEntity(ctx, entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	entity.ID = *id
	r.handleCreateTree(ctx, id, entity)

	return r.GetByID(ctx, *id)
}

func (r *TreeClusterRepository) createEntity(ctx context.Context, entity *entities.TreeCluster) (*int32, error) {
	args := sqlc.CreateTreeClusterParams{
		Region:         entity.Region,
		Address:        entity.Address,
		Description:    entity.Description,
		Latitude:       entity.Latitude,
		Longitude:      entity.Longitude,
		MoistureLevel:  entity.MoistureLevel,
		WateringStatus: sqlc.TreeClusterWateringStatus(entity.WateringStatus),
		SoilCondition:  sqlc.TreeSoilCondition(entity.SoilCondition),
	}

	id, err := r.store.CreateTreeCluster(ctx, &args)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return &id, nil
}

func (r *TreeClusterRepository) handleCreateTree(ctx context.Context, tcID *int32, entity *entities.TreeCluster) error {
	for _, tree := range entity.Trees {
		if tree.ID <= 0 {
			treeID, err := r.createTree(ctx, &entity.ID, tree)
			if err != nil {
				return err
			}
			tree.ID = treeID
		}
	}

	return nil
}

func (r *TreeClusterRepository) createTree(ctx context.Context, tcID *int32, tree *entities.Tree) (int32, error) {

	entity := sqlc.CreateTreeParams{
		TreeClusterID:       tcID,
		Species:             tree.Species,
		Age:                 tree.Age,
		HeightAboveSeaLevel: tree.HeightAboveSeaLevel,
		SensorID:            &tree.Sensor.ID,
		PlantingYear:        tree.PlantingYear,
		Latitude:            tree.Latitude,
		Longitude:           tree.Longitude,
	}

	return r.store.CreateTree(ctx, &entity)
}
