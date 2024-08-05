package sensor

import (
	"context"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type MqttService struct {
	sensorRepo  storage.SensorRepository
	isConnected bool
}

func NewMqttService(sensorRepository storage.SensorRepository) *MqttService {
	return &MqttService{sensorRepo: sensorRepository}
}

func (s *MqttService) HandleMessage(_ context.Context, _ *domain.MqttPayload) (*domain.MqttPayload, error) {
	// TODO: Implement the business logic of HandleMessage

	return nil, nil
}

func (s *MqttService) SetConnected(ready bool) {
	s.isConnected = ready
}

func (s *MqttService) Ready() bool {
	return s.isConnected
}
