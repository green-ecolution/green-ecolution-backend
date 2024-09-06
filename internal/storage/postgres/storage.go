package postgres

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed"
	flowerbedMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	sensorMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
	treeMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
	treeClusterMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/vehicle"
	vehicleMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/vehicle/mapper/generated"
	"github.com/jackc/pgx/v5"
)

func NewRepository(conn *pgx.Conn) *storage.Repository {
	store := store.NewStore(conn)
	querier := sqlc.New(conn)

	treeMappers := tree.NewTreeRepositoryMappers(
		&treeMapper.InternalTreeRepoMapperImpl{},
		&imgMapper.InternalImageRepoMapperImpl{},
	)
	treeRepo := tree.NewTreeRepository(store, treeMappers)

	tcMappers := treecluster.NewTreeClusterRepositoryMappers(
		&treeClusterMapper.InternalTreeClusterRepoMapperImpl{},
		&sensorMapper.InternalSensorRepoMapperImpl{},
	)
	treeClusterRepo := treecluster.NewTreeClusterRepository(store, tcMappers)

	imageMappers := image.NewImageRepositoryMappers(
		&imgMapper.InternalImageRepoMapperImpl{},
	)
	imageRepo := image.NewImageRepository(store, imageMappers)

	vehicleMappers := vehicle.NewVehicleRepositoryMappers(
		&vehicleMapper.InternalVehicleRepoMapperImpl{},
	)
	vehicleRepo := vehicle.NewVehicleRepository(querier, vehicleMappers)

	sensorMappers := sensor.NewSensorRepositoryMappers(
		&sensorMapper.InternalSensorRepoMapperImpl{},
	)
	sensorRepo := sensor.NewSensorRepository(querier, sensorMappers)

	flowMappers := flowerbed.NewFlowerbedMappers(
		&flowerbedMapper.InternalFlowerbedRepoMapperImpl{},
		&imgMapper.InternalImageRepoMapperImpl{},
		&sensorMapper.InternalSensorRepoMapperImpl{},
	)
	flowerbedRepo := flowerbed.NewFlowerbedRepository(store, flowMappers)

	return &storage.Repository{
		Tree:        treeRepo,
		TreeCluster: treeClusterRepo,
		Image:       imageRepo,
		Vehicle:     vehicleRepo,
		Sensor:      sensorRepo,
		Flowerbed:   flowerbedRepo,
	}
}
