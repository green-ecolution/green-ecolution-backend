package sensor

import (
	"context"
	"encoding/json"
	"log/slog"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/green-ecolution/green-ecolution-backend/internal/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/mapper/generated"
	sensorResponse "github.com/green-ecolution/green-ecolution-backend/internal/service/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sensorRepo "github.com/green-ecolution/green-ecolution-backend/internal/storage/entities/sensor"
)

type MqttService struct {
	sensorRepo  storage.SensorRepository
	mapper      mapper.MqttMapper
	isConnected bool
}

func NewMqttService(sensorRepository storage.SensorRepository) *MqttService {
	return &MqttService{sensorRepo: sensorRepository, mapper: &generated.MqttMapperImpl{}}
}

func (s *MqttService) HandleMessage(_ MQTT.Client, msg MQTT.Message) {
	jsonStr := string(msg.Payload())
  slog.Debug("Received message", "message", jsonStr)

	var sensorData sensorResponse.MqttPayloadResponse
	if err := json.Unmarshal([]byte(jsonStr), &sensorData); err != nil {
    slog.Error("Error unmarshalling sensor data", "error", err)
		return
	}

	payloadEntity := s.mapper.ToEntity(
		s.mapper.FromResponse(&sensorData),
	)
  slog.Debug("Mapped entity", "entity", payloadEntity)

	entity := &sensorRepo.MqttEntity{
		Data:   *payloadEntity,
		TreeID: "6686f54fd32cf640e8ae6eb1",
	}

	if _, err := s.sensorRepo.Insert(context.Background(), entity); err != nil {
    slog.Error("Error upserting sensor data", "error", err)
		return
	}
}

func (s *MqttService) SetConnected(ready bool) {
	s.isConnected = ready
}

func (s *MqttService) Ready() bool {
	return s.isConnected
}
