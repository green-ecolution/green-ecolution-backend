package tree

import (
	"context"

	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/pkg/errors"
)

func (r *TreeRepository) GetAll(ctx context.Context) ([]*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllTrees(ctx)
	if err != nil {
		log.Debug("failed to get trees in db", "error", err)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (r *TreeRepository) GetAllByProvider(ctx context.Context, provider string) ([]*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllTreesByProvider(ctx, &provider)
	if err != nil {
		log.Debug("failed to get trees in db", "error", err)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByID(ctx context.Context, id int32) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	row, err := r.store.GetTreeByID(ctx, id)
	if err != nil {
		log.Debug("failed to get tree by id in db", "error", err, "tree_id", id)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}

	t := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TreeRepository) GetBySensorID(ctx context.Context, id string) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	row, err := r.store.GetTreeBySensorID(ctx, &id)
	if err != nil {
		log.Debug("failed to get tree by sensor id in db", "error", err, "sensor_id", id)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}

	t := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TreeRepository) GetBySensorIDs(ctx context.Context, ids ...string) ([]*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetTreesBySensorIDs(ctx, ids)
	if err != nil {
		log.Debug("failed to get trees by multiple sensor ids in db", "error", err, "sensor_ids", ids)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (r *TreeRepository) GetTreesByIDs(ctx context.Context, ids []int32) ([]*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetTreesByIDs(ctx, ids)
	if err != nil {
		log.Debug("failed to get trees by ids in db", "error", err, "tree_ids", ids)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetTreesByTreeClusterID(ctx, &id)
	if err != nil {
		log.Debug("failed to get tree by cluster id in db", "error", err, "cluster_id", id)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}

	t := r.mapper.FromSqlList(rows)
	for _, tree := range t {
		if err := r.mapFields(ctx, tree); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (r *TreeRepository) GetByCoordinates(ctx context.Context, latitude, longitude float64) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	params := sqlc.GetTreeByCoordinatesParams{
		Latitude:  latitude,
		Longitude: longitude,
	}
	row, err := r.store.GetTreeByCoordinates(ctx, &params)
	if err != nil {
		log.Debug("failed to get tree by coordinates in db", "error", err, "latitude", latitude, "longitude", longitude)
		return nil, r.store.MapError(err, sqlc.Tree{})
	}
	tree := r.mapper.FromSql(row)
	if err := r.mapFields(ctx, tree); err != nil {
		return nil, err
	}
	return tree, nil
}

func (r *TreeRepository) GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error) {
	log := logger.GetLogger(ctx)
	rows, err := r.store.GetAllImagesByTreeID(ctx, id)
	if err != nil {
		log.Debug("failed to get images from tree id in db", "error", err, "tree_id", id)
		return nil, r.store.MapError(err, sqlc.Image{})
	}

	return r.iMapper.FromSqlList(rows), nil
}

func (r *TreeRepository) GetSensorByTreeID(ctx context.Context, treeID int32) (*entities.Sensor, error) {
	log := logger.GetLogger(ctx)
	row, err := r.store.GetSensorByTreeID(ctx, treeID)
	if err != nil {
		log.Debug("failed to get sensor by tree id", "error", err, "tree_id", treeID)
		return nil, r.store.MapError(err, sqlc.Sensor{})
	}

	data := r.sMapper.FromSql(row)
	if err := r.store.MapSensorFields(ctx, data); err != nil { // TODO: handle error
		return nil, err
	}

	return data, nil
}

func (r *TreeRepository) getTreeClusterByTreeID(ctx context.Context, treeID int32) (*entities.TreeCluster, error) {
	row, err := r.store.GetTreeClusterByTreeID(ctx, treeID)
	if err != nil {
		return nil, err
	}

	return r.tcMapper.FromSql(row), nil
}

// Map sensor, images and tree cluster entity to domain flowerbed
func (r *TreeRepository) mapFields(ctx context.Context, t *entities.Tree) error {
	if err := mapImages(ctx, r, t); err != nil {
		return err
	}

	if err := mapSensor(ctx, r, t); err != nil {
		return err
	}

	_ = mapTreeCluster(ctx, r, t)

	return nil
}

func mapImages(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	images, err := r.GetAllImagesByID(ctx, t.ID)
	if err != nil {
		return err
	}
	t.Images = images
	return nil
}

func mapSensor(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	sensor, err := r.GetSensorByTreeID(ctx, t.ID)
	if err != nil {
		var entityNotFoundErr storage.ErrEntityNotFound
		if errors.As(err, &entityNotFoundErr) {
			// If sensor is not found, set sensor to nil
			t.Sensor = nil
			return nil
		}
		return err
	}
	t.Sensor = sensor
	return nil
}

func mapTreeCluster(ctx context.Context, r *TreeRepository, t *entities.Tree) error {
	treeCluster, err := r.getTreeClusterByTreeID(ctx, t.ID)
	if err != nil {
		return err
	}
	t.TreeCluster = treeCluster
	return nil
}

func (r *TreeRepository) FindNearestTree(ctx context.Context, latitude, longitude float64) (*entities.Tree, error) {
	log := logger.GetLogger(ctx)
	params := &sqlc.FindNearestTreeParams{
		StMakepoint:   latitude,
		StMakepoint_2: longitude,
	}

	nearestTree, err := r.store.FindNearestTree(ctx, params)
	if err != nil {
		log.Debug("failed to find nearest tree on given coordinates", "error", err, "latitude", latitude, "longitude", longitude)
		return nil, err
	}

	tree := r.mapper.FromSql(nearestTree)
	if err := r.mapFields(ctx, tree); err != nil {
		return nil, err
	}
	return tree, nil
}
