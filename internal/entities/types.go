package entities

type Entities interface {
	Sensor |
		Flowerbed |
		Image |
		Vehicle |
		TreeCluster |
		Tree |
		Region |
		Departure
}

type EntityFunc[T Entities] func(*T)
