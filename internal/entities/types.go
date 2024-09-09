package entities

type Entities interface {
	Sensor |
		Flowerbed |
		Image |
		Vehicle |
		TreeCluster |
		Tree
}

type EntityFunc[T Entities] func(*T)
