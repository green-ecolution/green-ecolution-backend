package treecluster

import (
	"context"
	"log/slog"
	"time"

	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
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

func WithRegion(region string) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating region", "region", region)
		tc.Region = region
	}
}

func WithAddress(address string) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating address", "address", address)
		tc.Address = address
	}
}

func WithDescription(description string) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating description", "description", description)
		tc.Description = description
	}
}

func WithLatitude(latitude float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating latitude", "latitude", latitude)
		tc.Latitude = latitude
	}
}

func WithLongitude(longitude float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating longitude", "longitude", longitude)
		tc.Longitude = longitude
	}
}

func WithMoistureLevel(moistureLevel float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating moistureLevel", "moistureLevel", moistureLevel)
		tc.MoistureLevel = moistureLevel
	}
}

func WithWateringStatus(wateringStatus entities.TreeClusterWateringStatus) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating wateringStatus", "wateringStatus", wateringStatus)
		tc.WateringStatus = wateringStatus
	}
}

func WithSoilCondition(soilCondition entities.TreeSoilCondition) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating soilCondition", "soilCondition", soilCondition)
		tc.SoilCondition = soilCondition
	}
}

func WithLastWatered(lastWatered time.Time) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating lastWatered", "lastWatered", lastWatered)
		tc.LastWatered = &lastWatered
	}
}

func WithArchived(archived bool) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating archived", "archived", archived)
		tc.Archived = archived
	}
}

func WithTrees(trees []*entities.Tree) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		slog.Debug("updating trees", "trees", trees)
		tc.Trees = trees
	}
}

func (r *TreeClusterRepository) Archive(ctx context.Context, id int32) error {
	return r.store.ArchiveTreeCluster(ctx, id)
}

func (r *TreeClusterRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteTreeCluster(ctx, id)
}
