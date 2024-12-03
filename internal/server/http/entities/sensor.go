package entities

import (
	"time"
)

type SensorStatus string // @Name SensorStatus

const (
	SensorStatusOnline  SensorStatus = "online"
	SensorStatusOffline SensorStatus = "offline"
	SensorStatusUnknown SensorStatus = "unknown"
)

type SensorResponse struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Status    SensorStatus `json:"status"`
	Latitude  float64      `json:"latitude"`
	Longitude float64      `json:"longitude"`
} // @Name Sensor

type SensorListResponse struct {
	Data       []*SensorResponse `json:"data"`
	Pagination Pagination        `json:"pagination"`
} // @Name SensorList

type SensorDataResponse struct {
	ID               string  `json:"id"`
	BatteryLevel     float64 `json:"battery_level"`
	Temperature      float64 `json:"temperature"`
	Humidity         float64 `json:"humidity"`
	TrunkMoisture    float64 `json:"trunk_moisture"`
	SoilWaterTension float64 `json:"soil_water_tension"`
	Depth            float64 `json:"depth"`
} // @Name SensorData

type SensorDataListResponse struct {
	Data       []*SensorDataResponse `json:"data"`
	Pagination Pagination            `json:"pagination"`
} // @Name SensorDataList
