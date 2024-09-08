package entities

import "time"

type SensorStatus string

const (
	SensorStatusOnline  SensorStatus = "online"
	SensorStatusOffline SensorStatus = "offline"
	SensorStatusUnknown SensorStatus = "unknown"
)

type Sensor struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    SensorStatus
	Data []*SensorData
}

type SensorData struct {
	ID        int32
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *MqttPayload
}
