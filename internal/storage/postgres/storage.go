package postgres

import (
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed"
	mapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/vehicle"
	"github.com/jackc/pgx/v5"
)

func NewRepository(conn *pgx.Conn) *storage.Repository {
	store := store.NewStore(conn)
	querier := sqlc.New(conn)

	treeMappers := tree.NewTreeRepositoryMappers(
		&mapper.InternalTreeRepoMapperImpl{},
		&mapper.InternalImageRepoMapperImpl{},
    &mapper.InternalSensorRepoMapperImpl{},
    &mapper.InternalTreeClusterRepoMapperImpl{},
	)
	treeRepo := tree.NewTreeRepository(store, treeMappers)

	tcMappers := treecluster.NewTreeClusterRepositoryMappers(
		&mapper.InternalTreeClusterRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
	)
	treeClusterRepo := treecluster.NewTreeClusterRepository(store, tcMappers)

	imageMappers := image.NewImageRepositoryMappers(
		&mapper.InternalImageRepoMapperImpl{},
	)
	imageRepo := image.NewImageRepository(store, imageMappers)

	vehicleMappers := vehicle.NewVehicleRepositoryMappers(
		&mapper.InternalVehicleRepoMapperImpl{},
	)
	vehicleRepo := vehicle.NewVehicleRepository(querier, vehicleMappers)

	sensorMappers := sensor.NewSensorRepositoryMappers(
		&mapper.InternalSensorRepoMapperImpl{},
	)
	sensorRepo := sensor.NewSensorRepository(store, sensorMappers)

	flowMappers := flowerbed.NewFlowerbedMappers(
		&mapper.InternalFlowerbedRepoMapperImpl{},
		&mapper.InternalImageRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
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
