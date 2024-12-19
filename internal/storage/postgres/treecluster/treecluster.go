package treecluster

import (
	"context"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

type TreeClusterRepository struct {
	store *store.Store
	TreeClusterMappers
}

type TreeClusterMappers struct {
	mapper       mapper.InternalTreeClusterRepoMapper
	sensorMapper mapper.InternalSensorRepoMapper
	regionMapper mapper.InternalRegionRepoMapper
	treeMapper   mapper.InternalTreeRepoMapper
}

func NewTreeClusterRepositoryMappers(
	tcMapper mapper.InternalTreeClusterRepoMapper,
	sMapper mapper.InternalSensorRepoMapper,
	rMapper mapper.InternalRegionRepoMapper,
	tMapper mapper.InternalTreeRepoMapper,
) TreeClusterMappers {
	return TreeClusterMappers{
		mapper:       tcMapper,
		sensorMapper: sMapper,
		regionMapper: rMapper,
		treeMapper:   tMapper,
	}
}

func NewTreeClusterRepository(s *store.Store, mappers TreeClusterMappers) storage.TreeClusterRepository {
	return &TreeClusterRepository{
		store:              s,
		TreeClusterMappers: mappers,
	}
}

func WithName(name string) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Name = name
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

func WithLatitude(latitude *float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Latitude = latitude
	}
}

func WithLongitude(longitude *float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.Longitude = longitude
	}
}

func WithMoistureLevel(moistureLevel float64) entities.EntityFunc[entities.TreeCluster] {
	return func(tc *entities.TreeCluster) {
		tc.MoistureLevel = moistureLevel
	}
}

func WithWateringStatus(wateringStatus entities.WateringStatus) entities.EntityFunc[entities.TreeCluster] {
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
	_, err := r.store.ArchiveTreeCluster(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *TreeClusterRepository) Delete(ctx context.Context, id int32) error {
	_, err := r.store.DeleteTreeCluster(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
