package entities

type Entities interface {
	Sensor |
		Vehicle |
		TreeCluster |
		Tree |
		Region |
		WateringPlan
}

type EntityFunc[T Entities] func(*T)
