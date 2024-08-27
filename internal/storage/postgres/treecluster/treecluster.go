package treecluster

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster/mapper"
)

type TreeClusterRepository struct {
	querier sqlc.Querier
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

func NewTreeClusterRepository(querier sqlc.Querier, mappers TreeClusterMappers) storage.TreeClusterRepository {
	return &TreeClusterRepository{
		querier:            querier,
		TreeClusterMappers: mappers,
	}
}

func (r *TreeClusterRepository) GetAll(ctx context.Context) ([]*entities.TreeCluster, error) {
	rows, err := r.querier.GetAllTreeClusters(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlList(rows), nil
}

func (r *TreeClusterRepository) GetByID(ctx context.Context, id int32) (*entities.TreeCluster, error) {
	row, err := r.querier.GetTreeClusterByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSql(row), nil
}

func (r *TreeClusterRepository) GetSensorByTreeClusterID(ctx context.Context, id int32) (*entities.Sensor, error) {
	row, err := r.querier.GetSensorByTreeClusterID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.sensorMapper.FromSql(row), nil
}

func (r *TreeClusterRepository) Create(ctx context.Context, tc *entities.TreeCluster) (*entities.TreeCluster, error) {
	entity := sqlc.CreateTreeClusterParams{
		Region:    tc.Region,
		Address:   tc.Address,
		Latitude:  tc.Latitude,
		Longitude: tc.Longitude,
	}

	id, err := r.querier.CreateTreeCluster(ctx, &entity)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *TreeClusterRepository) Update(ctx context.Context, tc *entities.TreeCluster) (*entities.TreeCluster, error) {
	entity := sqlc.UpdateTreeClusterParams{
		ID:        tc.ID,
		Region:    tc.Region,
		Address:   tc.Address,
		Latitude:  tc.Latitude,
		Longitude: tc.Longitude,
	}

	err := r.querier.UpdateTreeCluster(ctx, &entity)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, tc.ID)
}

func (r *TreeClusterRepository) Archive(ctx context.Context, id int32) error {
	return r.querier.ArchiveTreeCluster(ctx, id)
}

func (r *TreeClusterRepository) Delete(ctx context.Context, id int32) error {
	return r.querier.DeleteTreeCluster(ctx, id)
}
