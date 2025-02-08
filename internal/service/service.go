package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"reflect"
	"time"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

var (
	ErrIPNotFound            = errors.New("local ip not found")
	ErrIFacesNotFound        = errors.New("cant get interfaces")
	ErrIFacesAddressNotFound = errors.New("cant get interfaces address")
	ErrHostnameNotFound      = errors.New("cant get hostname")
	ErrValidation            = errors.New("validation error")

	ErrPluginRegistered       = NewError(BadRequest, "plugin already registered")
	ErrPluginNotRegistered    = NewError(BadRequest, "plugin not registered")
	ErrVehiclePlateTaken      = NewError(BadRequest, "number plate is already taken")
	ErrVehicleUnsupportedType = NewError(BadRequest, "vehicle type is not supported")
	ErrUserNotCorrectRole     = NewError(BadRequest, "user has an incorrect role")
)

type Error struct {
	Message string
	Code    ErrorCode
}

type ErrorLogMask int

// Bitmask
const (
	ErrorLogNothing        ErrorLogMask = -1
	ErrorLogAll            ErrorLogMask = 0
	ErrorLogEntityNotFound ErrorLogMask = (1 << iota)
	ErrorLogValidation
)

func NewError(code ErrorCode, msg string) Error {
	return Error{Code: code, Message: msg}
}

func (e Error) Error() string {
	return e.Message
}

func MapError(ctx context.Context, err error, errorMask ErrorLogMask) error {
	log := logger.GetLogger(ctx)
	var entityNotFoundErr storage.ErrEntityNotFound
	if errors.As(err, &entityNotFoundErr) {
		if errorMask&ErrorLogEntityNotFound == 0 {
			log.Error("can't find entity", "error", err)
		}
		return NewError(NotFound, entityNotFoundErr.Error())
	}

	if errors.Is(err, ErrValidation) {
		if errorMask&ErrorLogValidation == 0 {
			log.Error("failed to validate struct", "error", err)
		}
		return NewError(BadRequest, err.Error())
	}

	log.Error("an error has occurred", "error", err)
	return NewError(InternalError, err.Error())
}

type ErrorCode int

const (
	BadRequest    ErrorCode = 400
	Unauthorized  ErrorCode = 401
	Forbidden     ErrorCode = 403
	NotFound      ErrorCode = 404
	InternalError ErrorCode = 500
)

type BasicCrudService[T any, CreateType any, UpdateType any] interface {
	GetAll(ctx context.Context, provider string) ([]*T, error)
	GetByID(ctx context.Context, id int32) (*T, error)
	Create(ctx context.Context, createData *CreateType) (*T, error)
	Update(ctx context.Context, id int32, updateData *UpdateType) (*T, error)
	Delete(ctx context.Context, id int32) error
}

type InfoService interface {
	Service
	GetAppInfo(context.Context) (*domain.App, error)
	GetAppInfoResponse(context.Context) (*domain.App, error)
}

type TreeService interface {
	CrudService[domain.Tree, domain.TreeCreate, domain.TreeUpdate]
	ImportTree(ctx context.Context, trees []*domain.TreeImport) error
	GetBySensorID(ctx context.Context, id string) (*domain.Tree, error)
	HandleNewSensorData(context.Context, *domain.EventNewSensorData) error
}

type AuthService interface {
	Service
	LoginRequest(ctx context.Context, loginRequest *domain.LoginRequest) *domain.LoginResp
	LogoutRequest(ctx context.Context, logoutRequest *domain.Logout) error
	ClientTokenCallback(ctx context.Context, loginCallback *domain.LoginCallback) (*domain.ClientToken, error)
	Register(ctx context.Context, user *domain.RegisterUser) (*domain.User, error)
	RetrospectToken(ctx context.Context, token string) (*domain.IntroSpectTokenResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.ClientToken, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*domain.User, error)
	GetAllByRole(ctx context.Context, role domain.UserRole) ([]*domain.User, error)
}

type RegionService interface {
	Service
	GetAll(ctx context.Context) ([]*domain.Region, error)
	GetByID(ctx context.Context, id int32) (*domain.Region, error)
}

type TreeClusterService interface {
	Service
	// TODO: use CrudService as soon as every service has pagination
	// CrudService[domain.TreeCluster, domain.TreeClusterCreate, domain.TreeClusterUpdate]
	GetAll(ctx context.Context, provider string) ([]*domain.TreeCluster, int64, error)
	GetByID(ctx context.Context, id int32) (*domain.TreeCluster, error)
	Create(ctx context.Context, createData *domain.TreeClusterCreate) (*domain.TreeCluster, error)
	Update(ctx context.Context, id int32, updateData *domain.TreeClusterUpdate) (*domain.TreeCluster, error)
	Delete(ctx context.Context, id int32) error

	HandleUpdateTree(context.Context, *domain.EventUpdateTree) error
	HandleCreateTree(context.Context, *domain.EventCreateTree) error
	HandleDeleteTree(context.Context, *domain.EventDeleteTree) error
	HandleNewSensorData(context.Context, *domain.EventNewSensorData) error
	HandleUpdateWateringPlan(context.Context, *domain.EventUpdateWateringPlan) error
}

type SensorService interface {
	Service
	GetAll(ctx context.Context, provider string) ([]*domain.Sensor, error)
	GetByID(ctx context.Context, id string) (*domain.Sensor, error)
	Create(ctx context.Context, createData *domain.SensorCreate) (*domain.Sensor, error)
	Update(ctx context.Context, id string, updateData *domain.SensorUpdate) (*domain.Sensor, error)
	Delete(ctx context.Context, id string) error
	HandleMessage(ctx context.Context, payload *domain.MqttPayload) (*domain.SensorData, error)
	MapSensorToTree(ctx context.Context, sen *domain.Sensor) error
	RunStatusUpdater(ctx context.Context, interval time.Duration)
}

type CrudService[T any, CreateType any, UpdateType any] interface {
	Service
	BasicCrudService[T, CreateType, UpdateType]
}

type VehicleService interface {
	Service
	// TODO: use CrudService as soon as every service has pagination
	// CrudService[domain.Vehicle, domain.VehicleCreate, domain.VehicleUpdate]
	GetAll(ctx context.Context) ([]*domain.Vehicle, int64, error)
	GetByID(ctx context.Context, id int32) (*domain.Vehicle, error)
	Create(ctx context.Context, createData *domain.VehicleCreate) (*domain.Vehicle, error)
	Update(ctx context.Context, id int32, updateData *domain.VehicleUpdate) (*domain.Vehicle, error)
	Delete(ctx context.Context, id int32) error

	GetAllByType(ctx context.Context, vehicleType domain.VehicleType) ([]*domain.Vehicle, int64, error)
	GetByPlate(ctx context.Context, plate string) (*domain.Vehicle, error)
}

type WateringPlanService interface {
	CrudService[domain.WateringPlan, domain.WateringPlanCreate, domain.WateringPlanUpdate]
	PreviewRoute(ctx context.Context, transporterID int32, trailerID *int32, clusterIDs []int32) (*domain.GeoJSON, error)
	GetGPXFileStream(ctx context.Context, objName string) (io.ReadSeekCloser, error)
}

type PluginService interface {
	Service
	Register(ctx context.Context, plugin *domain.Plugin) (*domain.ClientToken, error)
	Get(ctx context.Context, slug string) (domain.Plugin, error)
	GetAll(ctx context.Context) ([]domain.Plugin, []time.Time)
	HeartBeat(ctx context.Context, slug string) error
	Unregister(ctx context.Context, slug string)
	StartCleanup(ctx context.Context) error
}

type Service interface {
	Ready() bool
}

type Services struct {
	InfoService         InfoService
	TreeService         TreeService
	AuthService         AuthService
	RegionService       RegionService
	TreeClusterService  TreeClusterService
	SensorService       SensorService
	VehicleService      VehicleService
	PluginService       PluginService
	WateringPlanService WateringPlanService
}

type ServicesInterface interface {
	AllServicesReady() bool
}

func (s *Services) AllServicesReady() bool {
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		service := v.Field(i).Interface()
		if srv, ok := service.(Service); ok {
			if !srv.Ready() {
				return false
			}
		} else {
			slog.Debug("Service does not implement the Service interface", "service", v.Field(i).Type().Name())
			return false
		}
	}
	return true
}
