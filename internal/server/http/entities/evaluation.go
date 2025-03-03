package entities

type EvaluationResponse struct {
	TreeCount                     int64                        `json:"tree_count"`
	TreeClusterCount              int64                        `json:"treecluster_count"`
	SensorCount                   int64                        `json:"sensor_count"`
	WateringPlanCount             int64                        `json:"watering_plan_count"`
	TotalWaterConsumption         int64                        `json:"total_water_consumption"`
	VehicleCouVehicleEvaluationnt []*VehicleEvaluationResponse `json:"vehicle_evaluation"`
} // @Name Evaluation

type VehicleEvaluationResponse struct {
	NumberPlate       string `json:"number_plate"`
	WateringPlanCount int64  `json:"watering_plan_count"`
} // @Name VehicleEvaluation
