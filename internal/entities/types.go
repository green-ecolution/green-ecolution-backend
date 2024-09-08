package entities

type Entities interface {
	Sensor |
		Flowerbed |
		Image |
		Vehicle |
		TreeCluster |
		Tree
}

type CreateEntity interface {
	CreateSensor |
		CreateFlowerbed |
		CreateImage |
		CreateVehicle |
		CreateTreeCluster |
		CreateTree
}

type EntityFunc[T Entities] func(*T)
