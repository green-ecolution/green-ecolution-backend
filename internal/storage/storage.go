package storage

import (
	"context"
	"errors"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var (
	ErrIPNotFound            = errors.New("local ip not found")
	ErrIFacesNotFound        = errors.New("cant get interfaces")
	ErrIFacesAddressNotFound = errors.New("cant get interfaces address")
	ErrHostnameNotFound      = errors.New("cant get hostname")
	ErrCannotGetAppURL       = errors.New("cannot get app url")

	ErrIDNotFound          = errors.New("entity id not found")
	ErrIDAlreadyExists     = errors.New("entity id already exists")
	ErrEntityNotFound      = errors.New("entity not found")
	ErrSensorNotFound      = errors.New("sensor not found")
	ErrImageNotFound       = errors.New("image not found")
	ErrFlowerbedNotFound   = errors.New("flowerbed not found")
	ErrTreeClusterNotFound = errors.New("treecluster not found")

	ErrUnknowError      = errors.New("unknown error")
	ErrToManyRows       = errors.New("recieve more rows then expected")
	ErrConnectionClosed = errors.New("connection is closed")
	ErrTxClosed         = errors.New("transaction closed")
	ErrTxCommitRollback = errors.New("transaction cannot commit or rollback")
)

type BasicCrudRepository[T any, C any, U any] interface {
	GetAll(ctx context.Context) ([]*T, error)
	GetByID(ctx context.Context, id int32) (*T, error)

	Create(ctx context.Context, c *C) (*T, error)
	Update(ctx context.Context, u *U) (*T, error)
	Delete(ctx context.Context, id int32) error
}

type InfoRepository interface {
	GetAppInfo(context.Context) (*entities.App, error)
}

type UserRepository interface {
	BasicCrudRepository[entities.User, entities.User, entities.User]
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	UpdatePassword(ctx context.Context, id int32, password string) error
	Deactivate(ctx context.Context, id int32) error
	Activate(ctx context.Context, id int32) error
	AddRole(ctx context.Context, userID, roleID int32) error
	RemoveRole(ctx context.Context, userID, roleID int32) error
}

type RoleRepository interface {
	BasicCrudRepository[entities.Role, entities.Role, entities.Role]
	GetByName(ctx context.Context, name string) (*entities.Role, error)
}

type ImageRepository interface {
	BasicCrudRepository[entities.Image, entities.CreateImage, entities.UpdateImage]
}

type VehicleRepository interface {
	BasicCrudRepository[entities.Vehicle, entities.Vehicle, entities.Vehicle]
	GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error)
}

type TreeClusterRepository interface {
	BasicCrudRepository[entities.TreeCluster, entities.CreateTreeCluster, entities.UpdateTreeCluster]
	GetSensorByTreeClusterID(ctx context.Context, id int32) (*entities.Sensor, error)
  UpdateSoilCondition(ctx context.Context, id int32, soilCondition entities.TreeSoilCondition) error
  UpdateWateringStatus(ctx context.Context, id int32, wateringStatus entities.TreeClusterWateringStatus) error
  UpdateMoistureLevel(ctx context.Context, id int32, moistureLevel float64) error
  UpdateLastWatered(ctx context.Context, id int32, lastWatered time.Time) error
  UpdateGeometry(ctx context.Context, id int32, latitude float64, longitude float64) error
	Archive(ctx context.Context, id int32) error
}

type TreeRepository interface {
	BasicCrudRepository[entities.Tree, entities.CreateTree, entities.UpdateTree]
	GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
}

type SensorRepository interface {
	BasicCrudRepository[entities.Sensor, entities.CreateSensor, entities.UpdateSensor]
	GetStatusByID(ctx context.Context, id int32) (*entities.SensorStatus, error)
	GetSensorByStatus(ctx context.Context, status *entities.SensorStatus) ([]*entities.Sensor, error)
	GetSensorDataByID(ctx context.Context, id int32) ([]*entities.SensorData, error)
	InsertSensorData(ctx context.Context, data []*entities.SensorData) ([]*entities.SensorData, error)
}

type FlowerbedRepository interface {
	BasicCrudRepository[entities.Flowerbed, entities.CreateFlowerbed, entities.UpdateFlowerbed]
	GetSensorByFlowerbedID(ctx context.Context, id int32) (*entities.Sensor, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
	Archive(ctx context.Context, id int32) error
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entities.User, password string, role *[]string) (*entities.User, error)
	RetrospectToken(ctx context.Context, token string) (*entities.IntroSpectTokenResult, error)
	GetAccessTokenFromClientCode(ctx context.Context, code, redirectURL string) (*entities.ClientToken, error)
}

type Repository struct {
	Auth        AuthRepository
	Info        InfoRepository
	Sensor      SensorRepository
	Tree        TreeRepository
	User        UserRepository
	Role        RoleRepository
	Image       ImageRepository
	Vehicle     VehicleRepository
	TreeCluster TreeClusterRepository
	Flowerbed   FlowerbedRepository
}
