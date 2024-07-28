package tree

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	treeResponse "github.com/green-ecolution/green-ecolution-backend/internal/service/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sensorRepo "github.com/green-ecolution/green-ecolution-backend/internal/storage/entities/sensor"
	treeRepo "github.com/green-ecolution/green-ecolution-backend/internal/storage/entities/tree"
	"github.com/jinzhu/copier"
)

type TreeService struct {
	treeRepo     storage.TreeRepository
	sensorRepo   storage.SensorRepository
	treeMapper   mapper.TreeMapper
	sensorMapper mapper.MqttMapper
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository) *TreeService {
	return &TreeService{
		treeRepo:     repoTree,
		sensorRepo:   repoSensor,
		treeMapper:   &generated.TreeMapperImpl{},
		sensorMapper: &generated.MqttMapperImpl{},
	}
}

func (s *TreeService) fetchSensorData(ctx context.Context, treeID string) ([]*sensor.MqttPayload, error) {
	data, err := s.sensorRepo.GetAllByTreeID(ctx, treeID)
	if err != nil {
		return nil, err
	}

	return s.sensorMapper.FromEntityList(data), nil
}

func (s *TreeService) GetTreeByIDResponse(ctx context.Context, id string, withSensorData bool) (*treeResponse.TreeSensorDataResponse, error) {
	treeEntity, err := s.treeRepo.Get(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	treeData := s.treeMapper.FromEntity(treeEntity)
	var sensorData []*sensor.MqttPayload

	if withSensorData {
		data, err := s.fetchSensorData(ctx, id)
		if err != nil {
			return nil, handleError(err)
		}
		sensorData = data
	}

	response := treeResponse.TreeSensorDataResponse{
		Tree:       s.treeMapper.ToResponse(treeData),
		SensorData: s.sensorMapper.ToResponseList(sensorData),
	}
	return &response, nil
}

func (s *TreeService) GetAllTreesResponse(ctx context.Context, withSensorData bool) ([]treeResponse.TreeSensorDataResponse, error) {
	treeEntities, err := s.treeRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	treeData := s.treeMapper.FromEntityList(treeEntities)

	response := make([]treeResponse.TreeSensorDataResponse, len(treeData))
	var (
		sensorData map[string][]*sensor.MqttPayload
		wg         sync.WaitGroup
	)

	sensorData = make(map[string][]*sensor.MqttPayload)

	if withSensorData {
		wg.Add(len(treeEntities))
		for i := range treeEntities {
			go func(treeID string) {
				defer wg.Done()
				data, err := s.fetchSensorData(ctx, treeID)
				if err != nil {
					log.Printf("Error fetching sensor data for tree %s: %v", treeID, err)
					return
				}
				sensorData[treeID] = data
			}(treeData[i].ID)
		}
		wg.Wait()
	}

	for i, t := range treeData {
		response[i].Tree = s.treeMapper.ToResponse(t)
		if withSensorData {
			response[i].SensorData = s.sensorMapper.ToResponseList(sensorData[t.ID])
		}
	}

	return response, nil
}

func (s *TreeService) InsertTree(ctx context.Context, data *tree.Tree) error {
	entity := s.treeMapper.ToEntity(data)
	err := s.treeRepo.Insert(ctx, entity)
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil
}

func (s *TreeService) GetTreePredictionResponse(ctx context.Context, id string, withSensorData bool) (*treeResponse.TreeSensorPredictionResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	var treeEntity *treeRepo.TreeEntity
	var treeEntityError error
	var lastSensorEntity *sensorRepo.MqttEntity
	var lastSensorEntityError error

	go func() {
		defer wg.Done()
		treeEntity, treeEntityError = s.treeRepo.Get(ctx, id)
	}()

	go func() {
		defer wg.Done()
		lastSensorEntity, lastSensorEntityError = s.sensorRepo.GetLastByTreeID(ctx, id)
	}()

	wg.Wait()

	err := errors.Join(treeEntityError, lastSensorEntityError)
	if err != nil {
		return nil, handleError(err)
	}

	humidity := lastSensorEntity.Data.UplinkMessage.DecodedPayload.Humidity

	var mappedTree tree.Tree
	err = copier.Copy(&mappedTree, treeEntity)
	if err != nil {
		return nil, handleError(err)
	}

	prediction := &tree.SensorPrediction{
		SensorID: lastSensorEntity.Data.EndDeviceIDs.DeviceID,
		Humidity: humidity,
		Health:   getHealth(humidity),
		Tree:     &mappedTree,
	}

	predictionResponse := &treeResponse.SensorPredictionResponse{
		SensorID: prediction.SensorID,
		Humidity: prediction.Humidity,
		Health:   prediction.Health,
		Tree:     s.treeMapper.ToResponse(prediction.Tree),
	}

	var rawSensorData []*sensor.MqttPayload
	if withSensorData {
		rawSensorData, err = s.fetchSensorData(ctx, id)
		if err != nil {
			return nil, handleError(err)
		}
	}

	response := &treeResponse.TreeSensorPredictionResponse{
		SensorPrediction: predictionResponse,
		Tree:             s.treeMapper.ToResponse(&mappedTree),
		SensorData:       s.sensorMapper.ToResponseList(rawSensorData),
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
