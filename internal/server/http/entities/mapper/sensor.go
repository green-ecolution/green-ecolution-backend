package mapper

import (
	"encoding/json"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend MapSensorStatus
type SensorHTTPMapper interface {
  // goverter:ignore Type
	FromResponse(src *domain.Sensor) *entities.SensorResponse
}

func MapSensorData(src []byte) (*domain.MqttPayload, error) {
	var payload domain.MqttPayload
	err := json.Unmarshal(src, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func MapSensorStatus(src domain.SensorStatus) entities.SensorStatus {
	return entities.SensorStatus(src)
}
