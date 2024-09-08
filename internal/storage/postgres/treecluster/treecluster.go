package treecluster

import (
	"context"
	"errors"
	"log/slog"
	"time"

	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

type TreeClusterRepository struct {
	store *store.Store
	TreeClusterMappers
}

type TreeClusterMappers struct {
	mapper       sensorMapper.InternalTreeClusterRepoMapper
	sensorMapper sensorMapper.InternalSensorRepoMapper
}

func NewTreeClusterRepositoryMappers(tcMapper sensorMapper.InternalTreeClusterRepoMapper, sMapper sensorMapper.InternalSensorRepoMapper) TreeClusterMappers {
	return TreeClusterMappers{
		mapper:       tcMapper,
		sensorMapper: sMapper,
	}
}

func NewTreeClusterRepository(s *store.Store, mappers TreeClusterMappers) storage.TreeClusterRepository {
	s.SetEntityType(store.TreeCluster)
	return &TreeClusterRepository{
		store:              s,
		TreeClusterMappers: mappers,
	}
}

func (r *TreeClusterRepository) GetAll(ctx context.Context) ([]*entities.TreeCluster, error) {
	rows, err := r.store.GetAllTreeClusters(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *TreeClusterRepository) GetByID(ctx context.Context, id int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.mapper.FromSql(row), nil
}

func (r *TreeClusterRepository) GetSensorByTreeClusterID(ctx context.Context, id int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByTreeClusterID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.sensorMapper.FromSql(row), nil
}

func (r *TreeClusterRepository) Create(ctx context.Context, tc *entities.CreateTreeCluster) (*entities.TreeCluster, error) {
	entity := sqlc.CreateTreeClusterParams{
		Region:        tc.Region,
		Address:       tc.Address,
		Description:   tc.Description,
		MoistureLevel: tc.MoistureLevel,
		Latitude:      tc.Latitude,
		Longitude:     tc.Longitude,
	}

	id, err := r.store.CreateTreeCluster(ctx, &entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	if tc.WateringStatus != nil && *tc.WateringStatus != entities.TreeClusterWateringStatusUnknown {
		if err := r.UpdateWateringStatus(ctx, id, *tc.WateringStatus); err != nil {
			return nil, err
		}
	}

	if tc.SoilCondition != nil && *tc.SoilCondition != entities.TreeSoilConditionUnknown {
		if err := r.UpdateSoilCondition(ctx, id, *tc.SoilCondition); err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, id)
}

func (r *TreeClusterRepository) UpdateSoilCondition(ctx context.Context, id int32, soilCondition entities.TreeSoilCondition) error {
	args := sqlc.UpdateTreeClusterSoilConditionParams{
		ID:            id,
		SoilCondition: sqlc.TreeSoilCondition(soilCondition),
	}
	err := r.store.UpdateTreeClusterSoilCondition(ctx, &args)
	if err != nil {
		return r.store.HandleError(err)
	}

	return nil
}

func (r *TreeClusterRepository) UpdateWateringStatus(ctx context.Context, id int32, wateringStatus entities.TreeClusterWateringStatus) error {
	args := sqlc.UpdateTreeClusterWateringStatusParams{
		ID:             id,
		WateringStatus: sqlc.TreeClusterWateringStatus(wateringStatus),
	}
	err := r.store.UpdateTreeClusterWateringStatus(ctx, &args)
	if err != nil {
		return r.store.HandleError(err)
	}

	return nil
}

func (r *TreeClusterRepository) UpdateMoistureLevel(ctx context.Context, id int32, moistureLevel float64) error {
	args := sqlc.UpdateTreeClusterMoistureLevelParams{
		ID:            id,
		MoistureLevel: moistureLevel,
	}
	err := r.store.UpdateTreeClusterMoistureLevel(ctx, &args)
	if err != nil {
		return r.store.HandleError(err)
	}

	return nil
}

func (r *TreeClusterRepository) UpdateLastWatered(ctx context.Context, id int32, lastWatered time.Time) error {
	args := sqlc.UpdateTreeClusterLastWateredParams{
		ID:          id,
		LastWatered: utils.TimeToPgTimestamp(lastWatered),
	}
	err := r.store.UpdateTreeClusterLastWatered(ctx, &args)
	if err != nil {
		return r.store.HandleError(err)
	}

	return nil
}

func (r *TreeClusterRepository) UpdateGeometry(ctx context.Context, id int32, latitude, longitude float64) error {
	// TODO: implement
	slog.Info("Not implemented yet", "function", "UpdateGeometry", "context", ctx, "id", id, "latitude", latitude, "longitude", longitude)

	return errors.New("not implemented")
}

func (r *TreeClusterRepository) Update(ctx context.Context, tc *entities.UpdateTreeCluster) (*entities.TreeCluster, error) {
	prev, err := r.GetByID(ctx, tc.ID)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	if err := r.updateEntity(ctx, prev, tc); err != nil {
		return nil, err
	}

	return r.updateAttributes(ctx, prev, tc)
}

func (r *TreeClusterRepository) updateEntity(ctx context.Context, prev *entities.TreeCluster, tc *entities.UpdateTreeCluster) error {
	entity := sqlc.UpdateTreeClusterParams{
		ID:          tc.ID,
		Region:      utils.CompareAndUpdate(prev.Region, tc.Region),
		Address:     utils.CompareAndUpdate(prev.Address, tc.Address),
		Description: utils.CompareAndUpdate(prev.Description, tc.Description),
		Latitude:    utils.CompareAndUpdate(prev.Latitude, tc.Latitude),
		Longitude:   utils.CompareAndUpdate(prev.Longitude, tc.Longitude),
	}

	return r.store.UpdateTreeCluster(ctx, &entity)
}

func (r *TreeClusterRepository) updateAttributes(ctx context.Context, prev *entities.TreeCluster, tc *entities.UpdateTreeCluster) (*entities.TreeCluster, error) {
	updateFuncs := []func(context.Context, *entities.TreeCluster, *entities.UpdateTreeCluster) error{
		r.UpdateMoistureLevelIfChanged,
		r.UpdateWateringStatusIfChanged,
		r.UpdateSoilConditionIfChanged,
		r.UpdateLastWateredIfChanged,
		r.ArchiveIfChanged,
	}

	for _, fn := range updateFuncs {
		if err := fn(ctx, prev, tc); err != nil {
			return nil, err
		}
	}

	return r.GetByID(ctx, tc.ID)
}

func (r *TreeClusterRepository) UpdateMoistureLevelIfChanged(ctx context.Context, prev *entities.TreeCluster, tc *entities.UpdateTreeCluster) error {
	if tc.MoistureLevel == nil {
		return nil
	}

	return r.UpdateMoistureLevel(ctx, prev.ID, *tc.MoistureLevel)
}

func (r *TreeClusterRepository) UpdateWateringStatusIfChanged(ctx context.Context, prev *entities.TreeCluster, tc *entities.UpdateTreeCluster) error {
	if tc.WateringStatus == nil {
		return nil
	}

	return r.UpdateWateringStatus(ctx, prev.ID, *tc.WateringStatus)
}

func (r *TreeClusterRepository) UpdateSoilConditionIfChanged(ctx context.Context, prev *entities.TreeCluster, tc *entities.UpdateTreeCluster) error {
	if tc.SoilCondition == nil {
		return nil
	}

	return r.UpdateSoilCondition(ctx, prev.ID, *tc.SoilCondition)
}

func (r *TreeClusterRepository) UpdateLastWateredIfChanged(ctx context.Context, prev *entities.TreeCluster, tc *entities.UpdateTreeCluster) error {
	if tc.LastWatered == nil {
		return nil
	}

	return r.UpdateLastWatered(ctx, prev.ID, *tc.LastWatered)
}

func (r *TreeClusterRepository) ArchiveIfChanged(ctx context.Context, prev *entities.TreeCluster, tc *entities.UpdateTreeCluster) error {
	if tc.Archived == nil {
		return nil
	}

	if *tc.Archived {
		return r.Archive(ctx, prev.ID)
	}

	return nil
}

func (r *TreeClusterRepository) Archive(ctx context.Context, id int32) error {
	return r.store.ArchiveTreeCluster(ctx, id)
}

func (r *TreeClusterRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteTreeCluster(ctx, id)
}
