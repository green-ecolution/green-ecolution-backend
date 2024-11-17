package postgres

import (
	"log/slog"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/flowerbed"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image"
	mapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/region"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/vehicle"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepository(conn *pgxpool.Pool) *storage.Repository {
	s, err := store.NewStore(conn, sqlc.New(conn))
  if err != nil {
    slog.Error("failed to create store", "error", err)
    panic(err)
  }

	treeMappers := tree.NewTreeRepositoryMappers(
		&mapper.InternalTreeRepoMapperImpl{},
		&mapper.InternalImageRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
		&mapper.InternalTreeClusterRepoMapperImpl{},
	)
	treeRepo := tree.NewTreeRepository(s, treeMappers)

	tcMappers := treecluster.NewTreeClusterRepositoryMappers(
		&mapper.InternalTreeClusterRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
		&mapper.InternalRegionRepoMapperImpl{},
		&mapper.InternalTreeRepoMapperImpl{},
	)
	treeClusterRepo := treecluster.NewTreeClusterRepository(s, tcMappers)

	imageMappers := image.NewImageRepositoryMappers(
		&mapper.InternalImageRepoMapperImpl{},
	)
	imageRepo := image.NewImageRepository(s, imageMappers)

	vehicleMappers := vehicle.NewVehicleRepositoryMappers(
		&mapper.InternalVehicleRepoMapperImpl{},
	)
	vehicleRepo := vehicle.NewVehicleRepository(s, vehicleMappers)

	sensorMappers := sensor.NewSensorRepositoryMappers(
		&mapper.InternalSensorRepoMapperImpl{},
	)
	sensorRepo := sensor.NewSensorRepository(s, sensorMappers)

	flowMappers := flowerbed.NewFlowerbedMappers(
		&mapper.InternalFlowerbedRepoMapperImpl{},
		&mapper.InternalImageRepoMapperImpl{},
		&mapper.InternalSensorRepoMapperImpl{},
		&mapper.InternalRegionRepoMapperImpl{},
	)
	flowerbedRepo := flowerbed.NewFlowerbedRepository(s, flowMappers)

	regionMappers := region.NewRegionMappers(
		&mapper.InternalRegionRepoMapperImpl{},
	)
	regionRepo := region.NewRegionRepository(s, regionMappers)

	return &storage.Repository{
		Tree:        treeRepo,
		TreeCluster: treeClusterRepo,
		Image:       imageRepo,
		Vehicle:     vehicleRepo,
		Sensor:      sensorRepo,
		Flowerbed:   flowerbedRepo,
		Region:      regionRepo,
	}
}
