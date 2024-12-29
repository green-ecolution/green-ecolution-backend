package store

import (
	"context"

	"github.com/pkg/errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/jackc/pgx/v5"
)

var (
	regionMapper = generated.InternalRegionRepoMapperImpl{}
	treeMapper   = generated.InternalTreeRepoMapperImpl{}
	sensorMapper   = generated.InternalSensorRepoMapperImpl{}
)

// This function is required as soon as you want to add data to the tree cluster object
// from the database, e.g. the linked region or the linked trees.
// As this function is required in different repositories, it has been outsourced.
func (s *Store) MapClusterFields(ctx context.Context, tc *entities.TreeCluster) error {
	if err := s.mapRegion(ctx, tc); err != nil {
		return err
	}

	if err := s.mapTrees(ctx, tc); err != nil {
		return err
	}

	return nil
}

// This function is required as soon as you want to add the data to the sensor object
func (s *Store) MapSensorFields(ctx context.Context, sn *entities.Sensor) error {
	var err error

	sn.LatestData, err = s.GetLatestSensorDataBySensorID(ctx, sn.ID)
	if err != nil && !errors.Is(err, storage.ErrEntityNotFound) {
		return s.HandleError(err)
	}

	return nil
}

// This function provides the latest data from a specific sensor
func (s *Store) GetLatestSensorDataBySensorID(ctx context.Context, id string) (*entities.SensorData, error) {
	row, err := s.GetLatestSensorDataByID(ctx, id)
	if err != nil {
		return nil, s.HandleError(err)
	}

	domainData := sensorMapper.FromSqlSensorData(row)
	data, err := mapper.MapSensorData(row.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to map sensor data")
	}
	domainData.Data = data

	return domainData, nil
}


func (s *Store) mapRegion(ctx context.Context, tc *entities.TreeCluster) error {
	region, err := s.getRegionByTreeClusterID(ctx, tc.ID)
	if err != nil {
		// If region is not found, we can still return the tree cluster
		if !errors.Is(err, storage.ErrRegionNotFound) {
			return err
		}
	}
	tc.Region = region

	return nil
}

func (s *Store) mapTrees(ctx context.Context, tc *entities.TreeCluster) error {
	trees, err := s.getLinkedTreesByTreeClusterID(ctx, tc.ID)
	if err != nil {
		return err
	}
	tc.Trees = trees

	return nil
}

func (s *Store) getRegionByTreeClusterID(ctx context.Context, id int32) (*entities.Region, error) {
	row, err := s.GetRegionByTreeClusterID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrRegionNotFound
		}
		return nil, s.HandleError(err)
	}

	return regionMapper.FromSql(row), nil
}

func (s *Store) getLinkedTreesByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error) {
	rows, err := s.GetLinkedTreesByTreeClusterID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entities.Tree{}, nil
		}
		return nil, s.HandleError(err)
	}

	return treeMapper.FromSqlList(rows), nil
}
