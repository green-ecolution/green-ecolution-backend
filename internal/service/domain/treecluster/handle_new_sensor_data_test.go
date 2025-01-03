package treecluster

import (
	"context"
	"testing"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	storageMock "github.com/green-ecolution/green-ecolution-backend/internal/storage/_mock"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestTreeClusterService_HandleNewSensorData(t *testing.T) {
	t.Run("should update watering status on new sensor data event", func(t *testing.T) {
		clusterRepo := storageMock.NewMockTreeClusterRepository(t)
		treeRepo := storageMock.NewMockTreeRepository(t)
		regionRepo := storageMock.NewMockRegionRepository(t)
		eventManager := worker.NewEventManager(entities.EventTypeNewSensorData)
		svc := NewTreeClusterService(clusterRepo, treeRepo, regionRepo, eventManager)

		// event
		_, ch, _ := eventManager.Subscribe(entities.EventTypeNewSensorData)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go eventManager.Run(ctx)

		sensorDataEvent := entities.SensorData{
			SensorID: "sensor-1",
			Data: &entities.MqttPayload{
				Watermarks: []entities.Watermark{
					{Centibar: 30, Depth: 30},
					{Centibar: 40, Depth: 60},
					{Centibar: 50, Depth: 90},
				},
			},
		}

		tree := entities.Tree{
			ID: 1,
			TreeCluster: &entities.TreeCluster{
				ID:             1,
				WateringStatus: entities.WateringStatusUnknown,
			},
			PlantingYear: 2024,
		}

		treeWithSensorID1 := entities.Tree{
			ID: 2,
			TreeCluster: &entities.TreeCluster{
				ID:             1,
				WateringStatus: entities.WateringStatusUnknown,
			},
			Sensor: &entities.Sensor{
				ID: "sensor-1",
			},
			PlantingYear: 2022,
		}

		treeWithSensorID2 := entities.Tree{
			ID: 3,
			TreeCluster: &entities.TreeCluster{
				ID:             1,
				WateringStatus: entities.WateringStatusUnknown,
			},
			Sensor: &entities.Sensor{
				ID: "sensor-2",
			},
		}

		allLatestSensorData := []*entities.SensorData{
			{
				SensorID: "sensor-1",
				Data: &entities.MqttPayload{
					Watermarks: []entities.Watermark{
						{Centibar: 30, Depth: 30},
						{Centibar: 40, Depth: 60},
						{Centibar: 50, Depth: 90},
					},
				},
			},
			{
				SensorID: "sensor-2",
				Data: &entities.MqttPayload{
					Watermarks: []entities.Watermark{
						{Centibar: 30, Depth: 30},
						{Centibar: 40, Depth: 60},
						{Centibar: 50, Depth: 90},
					},
				},
			},
		}

		event := entities.NewEventSensorData(&sensorDataEvent)

		treeRepo.EXPECT().GetBySensorID(mock.Anything, "sensor-1").Return(&tree, nil)
		clusterRepo.EXPECT().GetAllLatestSensorDataByClusterID(mock.Anything, 1).Return(allLatestSensorData, nil)
		treeRepo.EXPECT().GetBySensorIDs(mock.Anything, "sensor-1", "sensor-2").Return([]*entities.Tree{&treeWithSensorID1, &treeWithSensorID2}, nil)

		// when
		err := svc.HandleNewSensorData(context.Background(), &event)

		// then
		assert.NoError(t, err)
		select {
		case recievedEvent := <-ch:
			_ = recievedEvent.(entities.EventNewSensorData)
			// TODO: Check if watering status has updated
		case <-time.After(1 * time.Second):
			t.Fatal("event was not received")
		}
	})
}
