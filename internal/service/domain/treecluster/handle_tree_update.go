package treecluster

import (
	"context"
	"log/slog"
	"slices"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	svcUtils "github.com/green-ecolution/green-ecolution-backend/internal/service/domain/utils"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func (s *TreeClusterService) HandleNewSensorData(ctx context.Context, event *entities.EventNewSensorData) error {
	slog.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")
	tree, err := s.treeRepo.GetBySensorID(ctx, event.New.SensorID)
	if err != nil {
		slog.Error("failed to get tree by sensor id", "sensor_id", event.New.SensorID, "err", err)
		return nil
	}

	sensorData, err := s.treeClusterRepo.GetAllLatestSensorDataByClusterID(ctx, tree.TreeCluster.ID)
	if err != nil {
		slog.Error("failed to get latest sensor data", "cluster_id", tree.TreeCluster.ID, "err", err)
		return nil
	}

	var wateringStatus entities.WateringStatus
	if len(sensorData) == 0 {
		return nil
	} else if len(sensorData) == 1 {
		wateringStatus = tree.WateringStatus
	} else {
		sensorIDs := utils.Map(sensorData, func(data *entities.SensorData) string {
			return data.SensorID
		})

		trees, err := s.treeRepo.GetBySensorIDs(ctx, sensorIDs...)
		if err != nil {
			slog.Error("failed to get trees by sensor id", "sensor_ids", sensorIDs, "err", err)
			return nil
		}

		slices.SortFunc(trees, func(a, b *entities.Tree) int {
			return int(a.PlantingYear - b.PlantingYear)
		})

		youngestTree := trees[0]

		var w30CentibarAvg, w60CentibarAvg, w90CentibarAvg int
		for _, data := range sensorData {
			w30, w60, w90, err := svcUtils.CheckAndSortWatermarks(data.Data.Watermarks)
			if err != nil {
				slog.Error("sensor data watermarks are malformed", "watermarks", data.Data.Watermarks)
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

func (s *TreeClusterService) HandleUpdateTree(ctx context.Context, event *entities.EventUpdateTree) error {
	slog.Debug("handle event", "event", event.Type(), "service", "TreeClusterService")

	if event.Prev.TreeCluster == nil && event.New.TreeCluster == nil {
		return nil
	}

	if s.isNoUpdateNeeded(event) {
		return nil
	}

	if err := s.handleTreeClusterUpdate(ctx, event.Prev.TreeCluster); err != nil {
		return err
	}

	if event.Prev.TreeCluster.ID != event.New.TreeCluster.ID {
		if err := s.handleTreeClusterUpdate(ctx, event.New.TreeCluster); err != nil {
			return err
		}
	}

	return nil
}

func (s *TreeClusterService) isNoUpdateNeeded(event *entities.EventUpdateTree) bool {
	treePosSame := event.Prev.Latitude == event.New.Latitude && event.Prev.Longitude == event.New.Longitude
	tcSame := event.Prev.TreeCluster != nil && event.New.TreeCluster != nil && event.Prev.TreeCluster.ID == event.New.TreeCluster.ID
	return treePosSame && tcSame
}

func (s *TreeClusterService) handleTreeClusterUpdate(ctx context.Context, tc *entities.TreeCluster) error {
	if tc == nil {
		return nil
	}

	updateFn := func(tc *entities.TreeCluster) (bool, error) {
		lat, long, region, err := s.getUpdatedLatLong(ctx, tc)
		if err != nil {
			return false, err
		}
		tc.Latitude = lat
		tc.Longitude = long
		tc.Region = region
		return true, nil
	}

	if err := s.treeClusterRepo.Update(ctx, tc.ID, updateFn); err == nil {
		return s.publishUpdateEvent(ctx, tc)
	}
	return nil
}
