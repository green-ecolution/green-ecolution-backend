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
)

type BasicCrudRepository[T any] interface {
	GetAll(ctx context.Context) ([]*T, error)
	GetByID(ctx context.Context, id int32) (*T, error)

	Create(ctx context.Context, t *T) (*T, error)
	Update(ctx context.Context, t *T) (*T, error)
	Delete(ctx context.Context, id int32) error
}

type InfoRepository interface {
	GetAppInfo(context.Context) (*entities.App, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entities.User, password string, roles *[]string) (*entities.User, error)
	GetByAccessToken(ctx context.Context, token string) (*entities.User, error)

	RemoveSession(ctx context.Context, token string) error
}

type RoleRepository interface {
	BasicCrudRepository[entities.Role]
	GetByName(ctx context.Context, name string) (*entities.Role, error)
}

type ImageRepository interface {
	BasicCrudRepository[entities.Image]
}

type VehicleRepository interface {
	BasicCrudRepository[entities.Vehicle]
	GetByPlate(ctx context.Context, plate string) (*entities.Vehicle, error)
}

type TreeClusterRepository interface {
	BasicCrudRepository[entities.TreeCluster]
	GetSensorByTreeClusterID(ctx context.Context, id int32) (*entities.Sensor, error)
	Archive(ctx context.Context, id int32) error
}

type TreeRepository interface {
	BasicCrudRepository[entities.Tree]
	GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
}

type SensorRepository interface {
	BasicCrudRepository[entities.Sensor]
	GetStatusByID(ctx context.Context, id int32) (*entities.SensorStatus, error)
	GetSensorByStatus(ctx context.Context, status *entities.SensorStatus) ([]*entities.Sensor, error)
	GetSensorDataByID(ctx context.Context, id int32) ([]*entities.SensorData, error)
	InsertSensorData(ctx context.Context, data []*entities.SensorData) ([]*entities.SensorData, error)
}

type FlowerbedRepository interface {
	BasicCrudRepository[entities.Flowerbed]
	GetSensorByFlowerbedID(ctx context.Context, id int32) (*entities.Sensor, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
}

type AuthRepository interface {
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
