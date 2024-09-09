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

	ErrIDNotFound          = errors.New("entity id not found")
	ErrIDAlreadyExists     = errors.New("entity id already exists")
	ErrEntityNotFound      = errors.New("entity not found")
	ErrSensorNotFound      = errors.New("sensor not found")
	ErrImageNotFound       = errors.New("image not found")
	ErrFlowerbedNotFound   = errors.New("flowerbed not found")
	ErrTreeClusterNotFound = errors.New("treecluster not found")

	ErrUnknowError      = errors.New("unknown error")
	ErrToManyRows       = errors.New("receive more rows then expected")
	ErrConnectionClosed = errors.New("connection is closed")
	ErrTxClosed         = errors.New("transaction closed")
	ErrTxCommitRollback = errors.New("transaction cannot commit or rollback")
)

type BasicCrudRepository[T entities.Entities] interface {
	GetAll(ctx context.Context) ([]*T, error)
	GetByID(ctx context.Context, id int32) (*T, error)

	Create(ctx context.Context, fn ...entities.EntityFunc[T]) (*T, error)
	Update(ctx context.Context, id int32, fn ...entities.EntityFunc[T]) (*T, error)
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
	UpdateGeometry(ctx context.Context, id int32, latitude float64, longitude float64) error
	Archive(ctx context.Context, id int32) error
}

type TreeRepository interface {
	BasicCrudRepository[entities.Tree]
	GetByTreeClusterID(ctx context.Context, id int32) ([]*entities.Tree, error)
	GetAllImagesByID(ctx context.Context, id int32) ([]*entities.Image, error)
	UpdateWithImages(ctx context.Context, id int32, fFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error)
	DeleteAndUnlinkImages(ctx context.Context, id int32) error
	UnlinkAllImages(ctx context.Context, id int32) error
	UnlinkImage(ctx context.Context, flowerbedID, imageID int32) error
	CreateAndLinkImages(ctx context.Context, tcFn ...entities.EntityFunc[entities.Tree]) (*entities.Tree, error)
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

	CreateAndLinkImages(ctx context.Context, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error)
	UpdateWithImages(ctx context.Context, id int32, fFn ...entities.EntityFunc[entities.Flowerbed]) (*entities.Flowerbed, error)
	DeleteAndUnlinkImages(ctx context.Context, id int32) error
	UnlinkAllImages(ctx context.Context, id int32) error
	UnlinkImage(ctx context.Context, flowerbedID, imageID int32) error
	Archive(ctx context.Context, id int32) error
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
