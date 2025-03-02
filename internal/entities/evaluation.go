package entities

type Evaluation struct {
	TreeCount             int64
	TreeClusterCount      int64
	SensorCount           int64
	WateringPlanCount     int64
	TotalWaterConsumption int64
}

type VehicleEvaluation struct {
	NumberPlate       string
	WateringPlanCount int64
}
