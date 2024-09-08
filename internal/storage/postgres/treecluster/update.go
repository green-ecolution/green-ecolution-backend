package treecluster

import (
	"context"
	"errors"
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func (r *TreeClusterRepository) Update(ctx context.Context, id int32, tcFn ...entities.EntityFunc[entities.TreeCluster]) (*entities.TreeCluster, error) {
	tc, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	for _, fn := range tcFn {
		fn(tc)
	}

	err = r.updateEntity(ctx, tc)
	if err != nil {
		return nil, err
	}


	return tc, nil
}

func (r *TreeClusterRepository) UpdateGeometry(ctx context.Context, id int32, latitude, longitude float64) error {
	// TODO: implement
	slog.Info("Not implemented yet", "function", "UpdateGeometry", "context", ctx, "id", id, "latitude", latitude, "longitude", longitude)

	return errors.New("not implemented")
}

func (r *TreeClusterRepository) updateEntity(ctx context.Context, tc *entities.TreeCluster) error {
	args := sqlc.UpdateTreeClusterParams{
		ID:             tc.ID,
		Region:         tc.Region,
		Address:        tc.Address,
		Description:    tc.Description,
		Latitude:       tc.Latitude,
		Longitude:      tc.Longitude,
		MoistureLevel:  tc.MoistureLevel,
		WateringStatus: sqlc.TreeClusterWateringStatus(tc.WateringStatus),
		SoilCondition:  sqlc.TreeSoilCondition(tc.SoilCondition),
		LastWatered:    utils.TimeToPgTimestamp(tc.LastWatered),
		Archived:       tc.Archived,
	}

	return r.store.UpdateTreeCluster(ctx, &args)
}

