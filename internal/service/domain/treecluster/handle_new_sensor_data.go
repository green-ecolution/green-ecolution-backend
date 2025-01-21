package treecluster

import (
	"context"
	"errors"
	"slices"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	svcUtils "github.com/green-ecolution/green-ecolution-backend/internal/service/domain/utils"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func (s *TreeClusterService) HandleNewSensorData(ctx context.Context, event *entities.EventNewSensorData) error {
	log := logger.GetLogger(ctx)
	log.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")
	tree, err := s.treeRepo.GetBySensorID(ctx, event.New.SensorID)
	if err != nil {
		// when error, it can be because the sensor has not linked tree or the tree does not exists
		if errors.Is(err, storage.ErrSensorNotFound) {
			log.Error("failed to get sensor by id", "sensor_id", event.New.SensorID, "err", err)
			return nil
		}
		log.Info("the sensor has no selected tree. This event will be ignored", "sensor_id", event.New.SensorID, "err", err)
		return nil
	}

	if tree.TreeCluster == nil {
		log.Info("this tree will has no linked tree cluster. This event will be ignored", "tree_id", tree.ID, "err", err)
		return nil
	}

	sensorData, err := s.treeClusterRepo.GetAllLatestSensorDataByClusterID(ctx, tree.TreeCluster.ID)
	if err != nil {
		log.Error("failed to get latest sensor data", "cluster_id", tree.TreeCluster.ID, "err", err)
		return nil
	}

	var wateringStatus entities.WateringStatus
	if len(sensorData) == 0 {
		// assertion - if there is no sensor data after receiving the event, the world is ending
		return nil
	} else if len(sensorData) == 1 {
		wateringStatus = tree.WateringStatus
	} else {
		sensorIDs := utils.Map(sensorData, func(data *entities.SensorData) string {
			return data.SensorID
		})

		trees, err := s.treeRepo.GetBySensorIDs(ctx, sensorIDs...)
		if err != nil {
			log.Error("failed to get trees by sensor id", "sensor_ids", sensorIDs, "err", err)
			return nil
		}

		slices.SortFunc(trees, func(a, b *entities.Tree) int {
			return int(b.PlantingYear - a.PlantingYear)
		})

		youngestTree := trees[0]

		var w30CentibarAvg, w60CentibarAvg, w90CentibarAvg int
		for _, data := range sensorData {
			w30, w60, w90, err := svcUtils.CheckAndSortWatermarks(data.Data.Watermarks)
			if err != nil {
				log.Error("sensor data watermarks are malformed", "watermarks", data.Data.Watermarks)
				return nil
			}

			w30CentibarAvg += w30.Centibar
			w60CentibarAvg += w60.Centibar
			w90CentibarAvg += w90.Centibar
		}

		watermarks := []entities.Watermark{
			{
				Centibar: w30CentibarAvg / len(sensorData),
				Depth:    30,
			},
			{
				Centibar: w60CentibarAvg / len(sensorData),
				Depth:    60,
			},
			{
				Centibar: w90CentibarAvg / len(sensorData),
				Depth:    90,
			},
		}

		wateringStatus = svcUtils.CalculateWateringStatus(youngestTree.PlantingYear, watermarks)
	}

	if wateringStatus == tree.TreeCluster.WateringStatus {
		return nil
	}

	updateFn := func(tc *entities.TreeCluster) (bool, error) {
		tc.WateringStatus = wateringStatus
		return true, nil
	}

	if err := s.treeClusterRepo.Update(ctx, tree.TreeCluster.ID, updateFn); err == nil {
		return s.publishUpdateEvent(ctx, tree.TreeCluster)
	}

	return nil
}
