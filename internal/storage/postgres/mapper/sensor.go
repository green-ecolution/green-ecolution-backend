package mapper

import (
	"encoding/json"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	mqtt "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/entities"
)

// goverter:converter
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTime
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:PgTimestampToTimePtr
// goverter:extend github.com/green-ecolution/green-ecolution-backend/internal/utils:TimeToTime
// goverter:extend MapSensorStatus
type InternalSensorRepoMapper interface {
  // goverter:ignore Data
	FromSql(src *sqlc.Sensor) *entities.Sensor
	FromSqlList(src []*sqlc.Sensor) []*entities.Sensor

	// goverter:ignore Data
	FromSqlSensorData(src *sqlc.SensorDatum) *entities.SensorData

	FromDomainSensorData(src *entities.MqttPayload) *mqtt.MqttPayload
}

func MapSensorData(src []byte) (*entities.MqttPayload, error) {
	var payload entities.MqttPayload
	err := json.Unmarshal(src, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func MapSensorStatus(src sqlc.SensorStatus) entities.SensorStatus {
	return entities.SensorStatus(src)
}
