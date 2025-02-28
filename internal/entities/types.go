package entities

type Entities interface {
	Sensor |
		Image |
		Vehicle |
		TreeCluster |
		Tree |
		Region |
		WateringPlan |
		Evaluation
}

type EntityFunc[T Entities] func(*T)
