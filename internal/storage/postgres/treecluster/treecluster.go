package treecluster

import (
	"context"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper"
	. "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

type TreeClusterRepository struct {
	store *Store
	TreeClusterMappers
}

type TreeClusterMappers struct {
	mapper       mapper.InternalTreeClusterRepoMapper
	sensorMapper sensorMapper.InternalSensorRepoMapper
}

func NewTreeClusterRepositoryMappers(tcMapper mapper.InternalTreeClusterRepoMapper, sMapper sensorMapper.InternalSensorRepoMapper) TreeClusterMappers {
	return TreeClusterMappers{
		mapper:       tcMapper,
		sensorMapper: sMapper,
	}
}

func NewTreeClusterRepository(store *Store, mappers TreeClusterMappers) storage.TreeClusterRepository {
	store.SetEntityType(TreeCluster)
	return &TreeClusterRepository{
		store:              store,
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
	var desc string // empty string or description
	if tc.Description != nil {
		desc = *tc.Description
	}

	entity := sqlc.CreateTreeClusterParams{
		Region:      tc.Region,
		Address:     tc.Address,
		Description: desc,
		Latitude:    tc.Latitude,
		Longitude:   tc.Longitude,
	}

	id, err := r.store.CreateTreeCluster(ctx, &entity)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	if tc.MoistureLevel != nil {
		r.UpdateMoistureLevel(ctx, id, *tc.MoistureLevel)
	}

	if tc.WateringStatus != nil || *tc.WateringStatus != entities.TreeClusterWateringStatusUnknown {
		r.UpdateWateringStatus(ctx, id, *tc.WateringStatus)
	}

	if tc.SoilCondition != nil || *tc.SoilCondition != entities.TreeSoilConditionUnknown {
		r.UpdateSoilCondition(ctx, id, *tc.SoilCondition)
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

func UpdateGeometry(ctx context.Context, id int32, latitude float64, longitude float64) error {
	// TODO: implement

	return nil
}

func (r *TreeClusterRepository) Update(ctx context.Context, tc *entities.UpdateTreeCluster) (*entities.TreeCluster, error) {
	prev, err := r.GetByID(ctx, tc.ID)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	entity := sqlc.UpdateTreeClusterParams{
		ID:          tc.ID,
		Region:      utils.CompareAndUpdate(prev.Region, tc.Region),
		Address:     utils.CompareAndUpdate(prev.Address, tc.Address),
		Description: utils.CompareAndUpdate(prev.Description, tc.Description),
		Latitude:    utils.CompareAndUpdate(prev.Latitude, tc.Latitude),
		Longitude:   utils.CompareAndUpdate(prev.Longitude, tc.Longitude),
	}

	err = r.store.UpdateTreeCluster(ctx, &entity)
	if err != nil {
		return nil, err
	}

	if tc.MoistureLevel != nil && prev.MoistureLevel != *tc.MoistureLevel {
		r.UpdateMoistureLevel(ctx, tc.ID, *tc.MoistureLevel)
	}

	if tc.WateringStatus != nil && prev.WateringStatus != *tc.WateringStatus {
		r.UpdateWateringStatus(ctx, tc.ID, *tc.WateringStatus)
	}

	if tc.SoilCondition != nil && prev.SoilCondition != *tc.SoilCondition {
		r.UpdateSoilCondition(ctx, tc.ID, *tc.SoilCondition)
	}

	if tc.LastWatered != nil && prev.LastWatered != *tc.LastWatered {
		r.UpdateLastWatered(ctx, tc.ID, *tc.LastWatered)
	}

	return r.GetByID(ctx, tc.ID)
}

func (r *TreeClusterRepository) Archive(ctx context.Context, id int32) error {
	return r.store.ArchiveTreeCluster(ctx, id)
}

func (r *TreeClusterRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteTreeCluster(ctx, id)
}
