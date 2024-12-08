package storage

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var (
	ErrIPNotFound            = errors.New("local ip not found")
	ErrIFacesNotFound        = errors.New("cant get interfaces")
	ErrIFacesAddressNotFound = errors.New("cant get interfaces address")
	ErrHostnameNotFound      = errors.New("cant get hostname")
	ErrCannotGetAppURL       = errors.New("cannot get app url")

	ErrIDNotFound           = errors.New("entity id not found")
	ErrIDAlreadyExists      = errors.New("entity id already exists")
	ErrEntityNotFound       = errors.New("entity not found")
	ErrSensorNotFound       = errors.New("sensor not found")
	ErrImageNotFound        = errors.New("image not found")
	ErrFlowerbedNotFound    = errors.New("flowerbed not found")
	ErrTreeClusterNotFound  = errors.New("treecluster not found")
	ErrRegionNotFound       = errors.New("region not found")
	ErrTreeNotFound         = errors.New("tree not found")
	ErrVehicleNotFound      = errors.New("vehicle not found")
	ErrWateringPlanNotFound = errors.New("watering plan not found")

	ErrUnknowError      = errors.New("unknown error")
	ErrToManyRows       = errors.New("receive more rows then expected")
	ErrConnectionClosed = errors.New("connection is closed")
	ErrTxClosed         = errors.New("transaction closed")
	ErrTxCommitRollback = errors.New("transaction cannot commit or rollback")

	ErrInvalidLatitude  = errors.New("latitude must be between 90,-90")
	ErrInvalidLongitude = errors.New("longitude must be between 180,-180")
)

type BasicCrudRepository[T entities.Entities] interface {
	// GetAll returns all entities
	GetAll(ctx context.Context) ([]*T, error)
	// GetByID returns one entity by id
	GetByID(ctx context.Context, id int32) (*T, error)
	// Create creates a new entity. It accepts a list of EntityFunc[T] to apply to the new entity
	Create(ctx context.Context, fn ...entities.EntityFunc[T]) (*T, error)
	// Update updates a already existing entity. It accepts a list of EntityFunc[T] to apply to the entity
	Update(ctx context.Context, id int32, fn ...entities.EntityFunc[T]) (*T, error)
	// Delete deletes a entity by id
	Delete(ctx context.Context, id int32) error
}

type InfoRepository interface {
	GetAppInfo(context.Context) (*entities.App, error)
}

type RegionRepository interface {
	BasicCrudRepository[entities.Region]
	GetByName(ctx context.Context, name string) (*entities.Region, error)
	GetByPoint(ctx context.Context, latitude, longitude float64) (*entities.Region, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User, password string, roles []string) (*entities.User, error)
	RemoveSession(ctx context.Context, token string) error
}

type RoleRepository interface {
	GetByName(ctx context.Context, name string) (*entities.Role, error)
}

type ImageRepository interface {
	BasicCrudRepository[entities.Image]
}

type VehicleRepository interface {
	BasicCrudRepository[entities.Vehicle]
	GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error)
}

type WateringPlanRepository interface {
	// GetAll returns all watering plans
	GetAll(ctx context.Context) ([]*entities.WateringPlan, error)
	// GetByID returns one watering plan by id
	GetByID(ctx context.Context, id int32) (*entities.WateringPlan, error)
	// Get vehicles linked to a watering plan by the watering plan id and the vehicle type
	GetLinkedVehicleByIDAndType(ctx context.Context, id int32, vehicleType entities.VehicleType) (*entities.Vehicle, error)
	// Get tree cluster linked to a watering plan by the watering plan id
	GetLinkedTreeClustersByID(ctx context.Context, id int32) ([]*entities.TreeCluster, error)
	// Create creates a new watering plan. It accepts a function that takes a watering plan that can be modified. Any changes made to the plan will be saved in the storage. If the function returns true, the watering plan will be created, otherwise it will not be created.
	Create(ctx context.Context, fn func(tc *entities.WateringPlan) (bool, error)) (*entities.WateringPlan, error)
	// Update updates a watering plan by id. It takes the id of the watering plan to update and a function that takes a watering plan that can be modified. Any changes made to the plan will be saved updated in the storage. If the function returns true, the watering plan will be updated, otherwise it will not be updated.
	Update(ctx context.Context, id int32, fn func(tc *entities.WateringPlan) (bool, error)) error
	// Delete deletes a watering plan by id
	Delete(ctx context.Context, id int32) error
}

type TreeClusterRepository interface {
	// GetAll returns all tree clusters
	GetAll(ctx context.Context) ([]*entities.TreeCluster, error)
	// GetByID returns one tree cluster by id
	GetByID(ctx context.Context, id int32) (*entities.TreeCluster, error)
	// GetByIDs returns mutiple tree cluster by ids
	GetByIDs(ctx context.Context, ids []int32) ([]*entities.TreeCluster, error)
	// Create creates a new tree cluster. It accepts a function that takes a tree cluster that can be modified. Any changes made to the tree cluster will be saved in the storage. If the function returns true, the tree cluster will be created, otherwise it will not be created.
	Create(ctx context.Context, fn func(tc *entities.TreeCluster) (bool, error)) (*entities.TreeCluster, error)
	// Update updates a tree cluster by id. It takes the id of the tree cluster to update and a function that takes a tree cluster that can be modified. Any changes made to the tree cluster will be saved updated in the storage. If the function returns true, the tree cluster will be updated, otherwise it will not be updated.
	Update(ctx context.Context, id int32, fn func(tc *entities.TreeCluster) (bool, error)) error
	// Delete deletes a tree cluster by id
	Delete(ctx context.Context, id int32) error

	GetRegionByTreeClusterID(ctx context.Context, id int32) (*entities.Region, error)
	GetLinkedTreesByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error)
	Archive(ctx context.Context, id int32) error
	LinkTreesToCluster(ctx context.Context, treeClusterID int32, treeIDs []int32) error
}

type TreeRepository interface {
	BasicCrudRepository[entities.Tree]
	GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
	GetSensorByTreeID(ctx context.Context, id int32) (*entities.Sensor, error)
	GetTreesByIDs(ctx context.Context, ids []int32) ([]*entities.Tree, error)
	GetByCoordinates(ctx context.Context, latitude, longitude float64) (*entities.Tree, error)
	GetBySensorID(ctx context.Context, id string) (*entities.Tree, error)

	UpdateWithImages(ctx context.Context, id int32, fFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error)
	DeleteAndUnlinkImages(ctx context.Context, id int32) error
	UnlinkAllImages(ctx context.Context, id int32) error
	UnlinkTreeClusterID(ctx context.Context, treeClusterID int32) error
	UnlinkSensorID(ctx context.Context, sensorID string) error
	UnlinkImage(ctx context.Context, flowerbedID, imageID int32) error
	CreateAndLinkImages(ctx context.Context, tcFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error)
	GetCenterPoint(ctx context.Context, id []int32) (float64, float64, error)
}

type SensorRepository interface {
	GetAll(ctx context.Context) ([]*entities.Sensor, error)
	GetByID(ctx context.Context, id string) (*entities.Sensor, error)
	Create(ctx context.Context, fn ...entities.EntityFunc[entities.Sensor]) (*entities.Sensor, error)
	Update(ctx context.Context, id string, fn ...entities.EntityFunc[entities.Sensor]) (*entities.Sensor, error)
	Delete(ctx context.Context, id string) error

	GetStatusByID(ctx context.Context, id string) (*entities.SensorStatus, error)
	GetSensorByStatus(ctx context.Context, status *entities.SensorStatus) ([]*entities.Sensor, error)
	GetSensorDataByID(ctx context.Context, id string) ([]*entities.SensorData, error)
	InsertSensorData(ctx context.Context, data []*entities.SensorData, id string) ([]*entities.SensorData, error)
}

type FlowerbedRepository interface {
	BasicCrudRepository[entities.Flowerbed]
	GetSensorByFlowerbedID(ctx context.Context, id int32) (*entities.Sensor, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
	GetRegionByFlowerbedID(ctx context.Context, id int32) (*entities.Region, error)

	CreateAndLinkImages(ctx context.Context, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error)
	UpdateWithImages(ctx context.Context, id int32, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error)
	DeleteAndUnlinkImages(ctx context.Context, id int32) error
	UnlinkAllImages(ctx context.Context, id int32) error
	UnlinkImage(ctx context.Context, flowerbedID, imageID int32) error
	UnlinkSensorID(ctx context.Context, sensorID string) error
	Archive(ctx context.Context, id int32) error
}

type AuthRepository interface {
	RetrospectToken(ctx context.Context, token string) (*entities.IntroSpectTokenResult, error)
	GetAccessTokenFromClientCode(ctx context.Context, code, redirectURL string) (*entities.ClientToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (*entities.ClientToken, error)
	GetAccessTokenFromClientCredentials(ctx context.Context, clientID, clientSecret string) (*entities.ClientToken, error)
}

type Repository struct {
	Auth         AuthRepository
	Info         InfoRepository
	Sensor       SensorRepository
	Tree         TreeRepository
	User         UserRepository
	Role         RoleRepository
	Image        ImageRepository
	Vehicle      VehicleRepository
	TreeCluster  TreeClusterRepository
	Flowerbed    FlowerbedRepository
	Region       RegionRepository
	WateringPlan WateringPlanRepository
}
