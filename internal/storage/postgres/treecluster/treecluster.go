package treecluster

import (
	"context"
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
	regionMapper sensorMapper.InternalRegionRepoMapper
}

func NewTreeClusterRepositoryMappers(
	tcMapper sensorMapper.InternalTreeClusterRepoMapper,
	sMapper sensorMapper.InternalSensorRepoMapper,
	rMapper sensorMapper.InternalRegionRepoMapper,
) TreeClusterMappers {
	return TreeClusterMappers{
		mapper:       tcMapper,
		sensorMapper: sMapper,
		regionMapper: rMapper,
	}
}

func NewTreeClusterRepository(s *store.Store, mappers TreeClusterMappers) storage.TreeClusterRepository {
	s.SetEntityType(store.TreeCluster)
	return &TreeClusterRepository{
		store:              s,
		TreeClusterMappers: mappers,
	}
}

func WithRegion(region *entities.Region) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Region = region
	}
}

func WithAddress(address string) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Address = address
	}
}

func WithDescription(description string) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Description = description
	}
}

func WithLatitude(latitude float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Latitude = latitude
	}
}

func WithLongitude(longitude float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Longitude = longitude
	}
}

func WithMoistureLevel(moistureLevel float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.MoistureLevel = moistureLevel
	}
}

func WithWateringStatus(wateringStatus entities.TreeClusterWateringStatus) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.WateringStatus = wateringStatus
	}
}

func WithSoilCondition(soilCondition entities.TreeSoilCondition) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.SoilCondition = soilCondition
	}
}

func WithLastWatered(lastWatered time.Time) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.LastWatered = &lastWatered
	}
}

func WithArchived(archived bool) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Archived = archived
	}
}

func WithTrees(trees []*entities.Tree) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Trees = trees
	}
}

func (r *TreeClusterRepository) Archive(ctx context.Context, id int32) error {
	return r.store.ArchiveTreeCluster(ctx, id)
}

func (r *TreeClusterRepository) Delete(ctx context.Context, id int32) error {
	return r.store.DeleteTreeCluster(ctx, id)
}
