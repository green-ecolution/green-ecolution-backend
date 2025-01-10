package tree

import (
	"context"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (r *TreeRepository) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	rows, err := r.store.GetAllTrees(ctx)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	row, err := r.store.GetTreeByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, t); err != nil {
		return nil, r.store.HandleError(err)
	}

	return t, nil
}

func (r *TreeRepository) GetBySensorID(ctx context.Context, id string) (*entities.Tree, error) {
	_, err := r.store.GetSensorByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrSensorNotFound
		} else {
			return nil, r.store.HandleError(err)
		}
	}

	row, err := r.store.GetTreeBySensorID(ctx, &id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, t); err != nil {
		return nil, r.store.HandleError(err)
	}

	return t, nil
}

func (r *TreeRepository) GetBySensorIDs(ctx context.Context, ids ...string) ([]*entities.Tree, error) {
	rows, err := r.store.GetTreesBySensorIDs(ctx, ids)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetTreesByIDs(ctx context.Context, ids []int32) ([]*entities.Tree, error) {
	rows, err := r.store.GetTreesByIDs(ctx, ids)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error) {
	_, err := r.store.GetTreeClusterByID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	rows, err := r.store.GetTreesByTreeClusterID(ctx, &id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, r.store.HandleError(err)
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByCoordinates(ctx context.Context, latitude, longitude float64) (*entities.Tree, error) {
	params := sqlc.GetTreeByCoordinatesParams{
		Latitude:  latitude,
		Longitude: longitude,
	}
	row, err := r.store.GetTreeByCoordinates(ctx, &params)
	if err != nil {
		return nil, r.store.HandleError(err)
	}
	tree := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, tree); err != nil {
		return nil, r.store.HandleError(err)
	}
	return tree, nil
}

func (r *TreeRepository) GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error) {
	rows, err := r.store.GetAllImagesByTreeID(ctx, id)
	if err != nil {
		return nil, r.store.HandleError(err)
	}

	return r.iMapper.FromSqlList(rows), nil
}

func (r *TreeRepository) GetSensorByTreeID(ctx context.Context, treeID int32) (*entities.Sensor, error) {
	row, err := r.store.GetSensorByTreeID(ctx, treeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrSensorNotFound
		} else {
			return nil, r.store.HandleError(err)
		}
	}

	data := r.sMapper.FromSql(row)
	if err := r.store.MapSensorFields(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TreeRepository) GetAllLatestSensorDataByTreeID(ctx context.Context, treeID int32) ([]*entities.SensorData, error) {
	rows, err := r.store.GetAllLatestSensorDataByTreeID(ctx, treeID)
	if err != nil {
		return nil, r.store.HandleError(err)
	}
	domainData, err := r.sMapper.FromSqlSensorDataList(rows)
	if err != nil {
		return nil, errors.Wrap(err, "failed to map sensor data")
	}

	return domainData, nil
}

func (r *TreeRepository) getTreeClusterByTreeID(ctx context.Context, treeID int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByTreeID(ctx, treeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrTreeClusterNotFound
		} else {
			return nil, r.store.HandleError(err)
		}
	}

	return r.tcMapper.FromSql(row), nil
}

// Map sensor, images and tree cluster entity to domain flowerbed
func (r *TreeRepository) mapFields(ctx context.Context, t *entities.Tree) error {
	if err := mapImages(ctx, r, t); err != nil {
		return r.store.HandleError(err)
	}

	if err := mapSensor(ctx, r, t); err != nil {
		return r.store.HandleError(err)
	}

	_ = mapTreeCluster(ctx, r, t)

	return nil
}

func mapImages(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	images, err := r.GetAllImagesByID(ctx, t.ID)
	if err != nil {
		return r.store.HandleError(err)
	}
	t.Images = images
	return nil
}

func mapSensor(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	sensor, err := r.GetSensorByTreeID(ctx, t.ID)
	if err != nil {
		if errors.Is(err, storage.ErrSensorNotFound) {
			// If sensor is not found, set sensor to nil
			t.Sensor = nil
			return nil
		}
		return r.store.HandleError(err)
	}
	t.Sensor = sensor
	return nil
}

func mapTreeCluster(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	treeCluster, err := r.getTreeClusterByTreeID(ctx, t.ID)
	if err != nil {
		return r.store.HandleError(err)
	}
	t.TreeCluster = treeCluster
	return nil
}

func (r *TreeRepository) FindNearestTree(ctx context.Context, latitude, longitude float64) (*entities.Tree, error) {
	params := &sqlc.FindNearestTreeParams{
		StMakepoint:   latitude,
		StMakepoint_2: longitude,
	}

	nearestTree, err := r.store.FindNearestTree(ctx, params)
	if err != nil {
		return nil, err
	}

	tree := r.mapper.FromSql(nearestTree)
	if err := r.mapFields(ctx, tree); err != nil {
		return nil, r.store.HandleError(err)
	}
	return tree, nil
}
