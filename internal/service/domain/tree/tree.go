package tree

import (
	"context"
	"errors"
	"sync"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"golang.org/x/sync/errgroup"
)

type TreeService struct {
	treeRepo   storage.TreeRepository
	sensorRepo storage.SensorRepository
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository) *TreeService {
	return &TreeService{
		treeRepo:   repoTree,
		sensorRepo: repoSensor,
	}
}

func (s *TreeService) fetchSensorData(ctx context.Context, treeID string) ([]*sensor.MqttPayload, error) {
	data, err := s.sensorRepo.GetAllByTreeID(ctx, treeID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *TreeService) GetTreeByIDResponse(ctx context.Context, id string, withSensorData bool) (*tree.TreeSensorData, error) {
	treeData, err := s.treeRepo.Get(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	var sensorData []*sensor.MqttPayload
	if withSensorData {
		data, err := s.fetchSensorData(ctx, id)
		if err != nil {
			return nil, handleError(err)
		}
		sensorData = data
	}

	response := tree.TreeSensorData{
		Tree:       treeData,
		SensorData: sensorData,
	}
	return &response, nil
}

func (s *TreeService) GetAllTreesResponse(ctx context.Context, withSensorData bool) ([]*tree.TreeSensorData, error) {
	treeData, err := s.treeRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	response := make([]*tree.TreeSensorData, len(treeData))
	var (
		sensorData map[string][]*sensor.MqttPayload
		wg         sync.WaitGroup
	)

	sensorData = make(map[string][]*sensor.MqttPayload)

	if withSensorData {
		wg.Add(len(treeData))
		for i := range treeData {
			go func(treeID string) {
				defer wg.Done()
				data, err := s.fetchSensorData(ctx, treeID)
				if err != nil {
					return
				}
				sensorData[treeID] = data
			}(treeData[i].ID)
		}
		wg.Wait()
	}

	for i := range treeData {
		response[i] = &tree.TreeSensorData{
			Tree: treeData[i],
		}

		if withSensorData {
			response[i].SensorData = sensorData[treeData[i].ID]
		}
	}

	return response, nil
}

func (s *TreeService) InsertTree(ctx context.Context, data *tree.Tree) error {
	err := s.treeRepo.Insert(ctx, data)
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil
}

func (s *TreeService) GetTreePredictionResponse(ctx context.Context, id string, withSensorData bool) (*tree.TreeSensorPrediction, error) {
	wg, errgroupCtx := errgroup.WithContext(ctx)
	wg.SetLimit(2)

	var (
		treeData   *tree.Tree
		lastSensor *sensor.MqttPayload
		err        error
	)

	wg.Go(func() (treeError error) {
		treeData, treeError = s.treeRepo.Get(errgroupCtx, id)
		return
	})

	wg.Go(func() (lastSensorError error) {
		lastSensor, lastSensorError = s.sensorRepo.GetLastByTreeID(errgroupCtx, id)
		return
	})

	if err = wg.Wait(); err != nil {
		return nil, handleError(err)
	}

	humidity := lastSensor.UplinkMessage.DecodedPayload.Humidity

	prediction := &tree.SensorPrediction{
		SensorID: lastSensor.EndDeviceIDs.DeviceID,
		Humidity: humidity,
		Health:   getHealth(humidity),
		Tree:     treeData,
	}

	predictionResponse := &tree.SensorPrediction{
		SensorID: prediction.SensorID,
		Humidity: prediction.Humidity,
		Health:   prediction.Health,
		Tree:     prediction.Tree,
	}

	var rawSensorData []*sensor.MqttPayload
	if withSensorData {
		rawSensorData, err = s.fetchSensorData(ctx, id)
		if err != nil {
			return nil, handleError(err)
		}
	}

	response := &tree.TreeSensorPrediction{
		SensorPrediction: predictionResponse,
		Tree:             treeData,
		SensorData:       rawSensorData,
	}

	return response, nil
}

func getHealth(humidity int) tree.PredictedHealth {
	const (
		ThresholdBad      = 40
		ThresholdModerate = 70
	)
	if humidity < ThresholdBad {
		return tree.HealthBad
	} else if humidity < ThresholdModerate {
		return tree.HealthModerate
	}

	return tree.HealthGood
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrMongoDataNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}
