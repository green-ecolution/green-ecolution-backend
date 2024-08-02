package sensor

import (
	"context"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type MqttService struct {
	sensorRepo  storage.SensorRepository
	isConnected bool
}

func NewMqttService(sensorRepository storage.SensorRepository) *MqttService {
	return &MqttService{sensorRepo: sensorRepository}
}

func (s *MqttService) HandleMessage(ctx context.Context, payload *domain.MqttPayload) (*domain.MqttPayload, error) {
	data, err := s.sensorRepo.Insert(ctx, payload)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *MqttService) SetConnected(ready bool) {
	s.isConnected = ready
}

func (s *MqttService) Ready() bool {
	return s.isConnected
}
