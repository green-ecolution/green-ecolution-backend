package entities

type Entities interface {
	Sensor |
		Flowerbed |
		Image |
		Vehicle |
		TreeCluster |
		Tree
}

type UpdateEntity interface {
	UpdateSensor |
		UpdateFlowerbed |
		UpdateImage |
		UpdateVehicle |
		UpdateTreeCluster |
		UpdateTree
}

type CreateEntity interface {
	CreateSensor |
		CreateFlowerbed |
		CreateImage |
		CreateVehicle |
		CreateTreeCluster |
		CreateTree
}
