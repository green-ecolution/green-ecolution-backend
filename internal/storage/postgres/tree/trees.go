package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree/mapper"
)

type TreeRepository struct {
	querier sqlc.Querier
	TreeMappers
}

type TreeMappers struct {
	mapper    mapper.InternalTreeRepoMapper
	imgMapper imgMapper.InternalImageRepoMapper
}

func NewTreeRepositoryMappers(treeMapper mapper.InternalTreeRepoMapper, imageMapper imgMapper.InternalImageRepoMapper) TreeMappers {
	return TreeMappers{
		mapper:    treeMapper,
		imgMapper: imageMapper,
	}
}

func NewTreeRepository(querier sqlc.Querier, mappers TreeMappers) storage.TreeRepository {
	return &TreeRepository{
		querier:     querier,
		TreeMappers: mappers,
	}
}

func (r *TreeRepository) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	rows, err := r.querier.GetAllTrees(ctx)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlTreeList(rows), nil
}

func (r *TreeRepository) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	row, err := r.querier.GetTreeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlTree(row), nil
}

func (r *TreeRepository) GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error) {
	rows, err := r.querier.GetTreesByTreeClusterID(ctx, &id)
	if err != nil {
		return nil, err
	}

	return r.mapper.FromSqlTreeList(rows), nil
}

func (r *TreeRepository) GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error) {
	rows, err := r.querier.GetAllImagesByTreeID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.imgMapper.FromSqlList(rows), nil
}

func (r *TreeRepository) Create(ctx context.Context, tree *entities.Tree) (*entities.Tree, error) {
	entity := sqlc.CreateTreeParams{
		TreeClusterID:       &tree.TreeCluster.ID,
		Species:             tree.Species,
		Age:                 tree.Age,
		HeightAboveSeaLevel: tree.HeightAboveSeaLevel,
		SensorID:            &tree.Sensor.ID,
		PlantingYear:        tree.PlantingYear,
		Latitude:            tree.Latitude,
		Longitude:           tree.Longitude,
	}

	id, err := r.querier.CreateTree(ctx, &entity)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *TreeRepository) Update(ctx context.Context, tree *entities.Tree) (*entities.Tree, error) {
	entity := sqlc.UpdateTreeParams{
		ID:                  tree.ID,
		Species:             tree.Species,
		Age:                 tree.Age,
		HeightAboveSeaLevel: tree.HeightAboveSeaLevel,
		SensorID:            &tree.Sensor.ID,
		PlantingYear:        tree.PlantingYear,
		Latitude:            tree.Latitude,
		Longitude:           tree.Longitude,
	}

	err := r.querier.UpdateTree(ctx, &entity)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, tree.ID)
}

func (r *TreeRepository) Delete(ctx context.Context, id int32) error {
	return r.querier.DeleteTree(ctx, id)
}
