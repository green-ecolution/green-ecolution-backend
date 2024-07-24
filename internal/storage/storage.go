package storage

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/tree"
)

var (
	ErrIPNotFound            = errors.New("local ip not found")
	ErrIFacesNotFound        = errors.New("cant get interfaces")
	ErrIFacesAddressNotFound = errors.New("cant get interfaces address")
	ErrHostnameNotFound      = errors.New("cant get hostname")
	ErrCannotUpsertData      = errors.New("cannot upsert data")
	ErrDataNotFound          = errors.New("data not found")
)

type InfoRepository interface {
	GetAppInfo(context.Context) (*info.App, error)
}

type SensorRepository interface {
	Insert(ctx context.Context, data *sensor.MqttPayload) (*sensor.MqttPayload, error)
	Get(ctx context.Context, id string) (*sensor.MqttPayload, error)
	GetFirst(ctx context.Context) (*sensor.MqttPayload, error)
	GetAllByTreeID(ctx context.Context, treeID string) ([]*sensor.MqttPayload, error)
	GetLastByTreeID(ctx context.Context, treeID string) (*sensor.MqttPayload, error)
}

type TreeRepository interface {
	Insert(ctx context.Context, data *tree.Tree) error
	Get(ctx context.Context, id string) (*tree.Tree, error)
	GetAll(ctx context.Context) ([]*tree.Tree, error)
}

type SchemaRepository interface {
	SetupSensorTable()
	SetupTreeTable()
	TeardownDatabase()
}
type Repository struct {
	Info   InfoRepository
	Sensor SensorRepository
	Tree   TreeRepository
	Schema SchemaRepository
}
