package treecluster

import (
	"context"

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
