package entities

type Entities interface {
	Sensor |
		Flowerbed |
		Image |
		Vehicle |
		TreeCluster |
		Tree |
		Region |
		WateringPlan
}

type EntityFunc[T Entities] func(*T)
