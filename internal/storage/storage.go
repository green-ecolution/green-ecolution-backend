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

	ErrIDNotFound      = errors.New("entity id not found")
	ErrIDAlreadyExists = errors.New("entity id already exists")
  ErrEntityNotFound  = errors.New("entity not found")
	ErrSensorNotFound  = errors.New("sensor not found")
  ErrImageNotFound   = errors.New("image not found")
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
	BasicCrudRepository[entities.Image, entities.Image, entities.Image]
}

type VehicleRepository interface {
	BasicCrudRepository[entities.Vehicle, entities.Vehicle, entities.Vehicle]
	GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error)
}

type TreeClusterRepository interface {
	BasicCrudRepository[entities.TreeCluster, entities.TreeCluster, entities.TreeCluster]
	GetSensorByTreeClusterID(ctx context.Context, id int32) (*entities.Sensor, error)
	Archive(ctx context.Context, id int32) error
}

type TreeRepository interface {
	BasicCrudRepository[entities.Tree, entities.Tree, entities.Tree]
	GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
}

type SensorRepository interface {
	BasicCrudRepository[entities.Sensor, entities.Sensor, entities.Sensor]
	GetStatusByID(ctx context.Context, id int32) (*entities.SensorStatus, error)
	GetSensorByStatus(ctx context.Context, status *entities.SensorStatus) ([]*entities.Sensor, error)
	GetSensorDataByID(ctx context.Context, id int32) ([]*entities.SensorData, error)
	InsertSensorData(ctx context.Context, data []*entities.SensorData) ([]*entities.SensorData, error)
}

type FlowerbedRepository interface {
	BasicCrudRepository[entities.Flowerbed, entities.CreateFlowerbed, entities.UpdateFlowerbed]
	GetSensorByFlowerbedID(ctx context.Context, id int32) (*entities.Sensor, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
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
