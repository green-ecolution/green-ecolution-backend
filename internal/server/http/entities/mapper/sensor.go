package mapper

import (
	"encoding/json"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend MapSensorStatus MapSensorStatusReq
type SensorHTTPMapper interface {
	FromResponse(src *domain.Sensor) *entities.SensorResponse

	// goverter:ignore Data
	FromCreateRequest(*entities.SensorCreateRequest) *domain.SensorCreate
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

func MapSensorStatusReq(src entities.SensorStatus) domain.SensorStatus {
	return domain.SensorStatus(src)
}
