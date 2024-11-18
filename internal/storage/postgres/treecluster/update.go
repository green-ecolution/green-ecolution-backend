package treecluster

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func (r *TreeClusterRepository) Update(ctx context.Context, id int32, updateFn func(*entities.TreeCluster) (bool, error)) error {
	return r.store.WithTx(ctx, func(q *sqlc.Queries) error {
		cancel := r.store.SwitchQuerier(q)
		defer cancel()

		tc, err := r.GetByID(ctx, id)
		if err != nil {
			return err
		}

		if updateFn == nil {
			return errors.New("updateFn is nil")
		}

		updated, err := updateFn(tc)
		if err != nil {
			return err
		}

		if !updated {
			return nil
		}

		return r.updateEntity(ctx, tc)
	})
}

func (r *TreeClusterRepository) updateEntity(ctx context.Context, tc *entities.TreeCluster) error {
	var regionID *int32
	if tc.Region != nil {
		regionID = &tc.Region.ID
	}
	args := sqlc.UpdateTreeClusterParams{
		ID:             tc.ID,
		RegionID:       regionID,
		Address:        tc.Address,
		Description:    tc.Description,
		MoistureLevel:  tc.MoistureLevel,
		WateringStatus: sqlc.WateringStatus(tc.WateringStatus),
		SoilCondition:  sqlc.TreeSoilCondition(tc.SoilCondition),
		LastWatered:    utils.TimeToPgTimestamp(tc.LastWatered),
		Archived:       tc.Archived,
		Name:           tc.Name,
	}

	if err := r.store.UnlinkTreeClusterID(ctx, &tc.ID); err != nil {
		return err
	}

	if len(tc.Trees) > 0 {
		treeIDs := utils.Map(tc.Trees, func(t *entities.Tree) int32 {
			return t.ID
		})

		if err := r.LinkTreesToCluster(ctx, tc.ID, treeIDs); err != nil {
			return err
		}
	}

	if tc.Latitude == nil || tc.Longitude == nil {
		if err := r.store.RemoveTreeClusterLocation(ctx, tc.ID); err != nil {
			return err
		}
	} else {
		locationArgs := sqlc.SetTreeClusterLocationParams{
			ID:        tc.ID,
			Latitude:  tc.Latitude,
			Longitude: tc.Longitude,
		}
		if err := r.store.SetTreeClusterLocation(ctx, &locationArgs); err != nil {
			return err
		}
	}

	return r.store.UpdateTreeCluster(ctx, &args)
}
