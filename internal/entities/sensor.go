package entities

import "time"

type SensorStatus string

const (
	SensorStatusOnline  SensorStatus = "online"
	SensorStatusOffline SensorStatus = "offline"
	SensorStatusUnknown SensorStatus = "unknown"
)

type Sensor struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    SensorStatus
	Data      []*SensorData
}

type SensorData struct {
	ID        int32
	SensorID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *MqttPayload
}

type SensorCreate struct {
	Status SensorStatus
	Data   []*SensorData
}

type SensorUpdate struct {
	Status SensorStatus
	Data   []*SensorData
}
