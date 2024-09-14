package sensor

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/pagination"
)

type SensorStatus string // @Name SensorStatus

const (
	SensorStatusOnline  SensorStatus = "online"
	SensorStatusOffline SensorStatus = "offline"
	SensorStatusUnknown SensorStatus = "unknown"
)

type SensorResponse struct {
	ID        int32        `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Status    SensorStatus `json:"status"`
	Type      string       `json:"type"`
} // @Name Sensor

type SensorDataResponse struct {
	Data []*struct {
		ID               int32   `json:"id"`
		BatteryLevel     float64 `json:"battery_level"`
		Temperature      float64 `json:"temperature"`
		Humidity         float64 `json:"humidity"`
		TrunkMoisture    float64 `json:"trunk_moisture"`
		SoilWaterTension float64 `json:"soil_water_tension"`
		Depth            float64 `json:"depth"`
	} `json:"data,omitempty"`
	Pagination pagination.Pagination `json:"pagination"`
} // @Name SensorData
