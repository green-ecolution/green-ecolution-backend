package entities

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt/entities/sensor"
)

type SensorStatus string // @Name SensorStatus

const (
	SensorStatusOnline  SensorStatus = "online"
	SensorStatusOffline SensorStatus = "offline"
	SensorStatusUnknown SensorStatus = "unknown"
)

type SensorResponse struct {
	ID         string              `json:"id"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	Status     SensorStatus        `json:"status"`
	LatestData *SensorDataResponse `json:"latest_data"`
	Latitude   float64             `json:"latitude"`
	Longitude  float64             `json:"longitude"`
} // @Name Sensor

type SensorListResponse struct {
	Data       []*SensorResponse `json:"data"`
	Pagination Pagination        `json:"pagination"`
} // @Name SensorList

type SensorDataResponse struct {
	ID        int32                       `json:"id"`
	SensorID  string                      `json:"sensor_id"`
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
	Data      *sensor.MqttPayloadResponse `json:"data"`
} // @Name SensorData

type SensorDataListResponse struct {
	Data       []*SensorDataResponse `json:"data"`
	Pagination Pagination            `json:"pagination"`
} // @Name SensorDataList
