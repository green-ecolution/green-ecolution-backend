package treecluster

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func defaultTreeCluster() *entities.TreeCluster {
	return &entities.TreeCluster{
		Region:         nil,
		Address:        "",
		Description:    "",
		MoistureLevel:  0,
		Latitude:       nil,
		Longitude:      nil,
		WateringStatus: entities.WateringStatusUnknown,
		SoilCondition:  entities.TreeSoilConditionUnknown,
		Archived:       false,
		LastWatered:    nil,
		Trees:          make([]*entities.Tree, 0),
		Name:           "",
	}
}

func (r *TreeClusterRepository) Create(ctx context.Context, createFn func(*entities.TreeCluster) (bool, error)) (*entities.TreeCluster, error) {
	if createFn == nil {
		return nil, errors.New("createFn is nil")
	}

	var createdTc *entities.TreeCluster
	err := r.store.WithTx(ctx, func(q *sqlc.Queries) error {
		cancel := r.store.SwitchQuerier(q)
		defer cancel()

		entity := defaultTreeCluster()
		created, err := createFn(entity)
		if err != nil {
			return err
		}

		if !created {
			return nil
		}

		if err := r.validateTreeClusterEntity(entity); err != nil {
			return err
		}

		id, err := r.createEntity(ctx, entity)
		if err != nil {
			return err
		}
		createdTc, err = r.GetByID(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTc, nil
}

func (r *TreeClusterRepository) createEntity(ctx context.Context, entity *entities.TreeCluster) (int32, error) {
	var region *int32
	if entity.Region != nil {
		region = &entity.Region.ID
	}

	args := sqlc.CreateTreeClusterParams{
		RegionID:       region,
		Address:        entity.Address,
		Description:    entity.Description,
		MoistureLevel:  entity.MoistureLevel,
		WateringStatus: sqlc.WateringStatus(entity.WateringStatus),
		SoilCondition:  sqlc.TreeSoilCondition(entity.SoilCondition),
		Name:           entity.Name,
	}

	id, err := r.store.CreateTreeCluster(ctx, &args)
	if err != nil {
		return -1, err
	}

	if len(entity.Trees) > 0 {
		treeIDs := utils.Map(entity.Trees, func(t *entities.Tree) int32 {
			return t.ID
		})

		err = r.LinkTreesToCluster(ctx, id, treeIDs)
		if err != nil {
			return -1, err
		}
	}

	if entity.Latitude != nil && entity.Longitude != nil {
		err = r.store.SetTreeClusterLocation(ctx, &sqlc.SetTreeClusterLocationParams{
			ID:        id,
			Latitude:  entity.Latitude,
			Longitude: entity.Longitude,
		})
		if err != nil {
			return -1, err
		}
	}

	return id, nil
}

func (r *TreeClusterRepository) validateTreeClusterEntity(tc *entities.TreeCluster) error {
	if tc == nil {
		return errors.New("tree cluster is nil")
	}

	if tc.Name == "" {
		return errors.New("tree cluster name is empty")
	}

	return nil
}

func (r *TreeClusterRepository) LinkTreesToCluster(ctx context.Context, treeClusterID int32, treeIDs []int32) error {
	args := &sqlc.LinkTreesToTreeClusterParams{
		Column1:       treeIDs,
		TreeClusterID: &treeClusterID,
	}
	return r.store.LinkTreesToTreeCluster(ctx, args)
}
