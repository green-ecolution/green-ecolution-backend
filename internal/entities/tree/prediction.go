package tree

type PredictedHealth = string

const (
	HealthGood     PredictedHealth = "good"
	HealthModerate                 = "moderate"
	HealthBad                      = "bad"
)

type SensorPrediction struct {
	SensorID string          `json:"sensor_id"`
	Tree     *Tree            `json:"tree"`
	Health   PredictedHealth `json:"predicted_health"`
	Humidity int             `json:"humidity"`
}
