package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities/tree"
)

var (
	ErrIPNotFound            = errors.New("local ip not found")
	ErrIFacesNotFound        = errors.New("cant get interfaces")
	ErrIFacesAddressNotFound = errors.New("cant get interfaces address")
	ErrHostnameNotFound      = errors.New("cant get hostname")
	ErrValidation            = errors.New("validation error")
)

type Error struct {
	Message string
	Code    ErrorCode
}

func NewError(code ErrorCode, msg string) Error {
	return Error{Code: code, Message: msg}
}

func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

type ErrorCode int

const (
	BadRequest    ErrorCode = 400
	Unauthorized  ErrorCode = 401
	Forbidden     ErrorCode = 403
	NotFound      ErrorCode = 404
	InternalError ErrorCode = 500
)

type InfoService interface {
	Service
	GetAppInfo(context.Context) (*info.App, error)
	GetAppInfoResponse(context.Context) (*info.App, error)
}

type MqttService interface {
	Service
	HandleMessage(ctx context.Context, payload *sensor.MqttPayload) (*sensor.MqttPayload, error)
	SetConnected(bool)
}

type TreeService interface {
	Service
	InsertTree(ctx context.Context, data *tree.Tree) error

	GetAllTreesResponse(ctx context.Context, withSensorData bool) ([]*tree.TreeSensorData, error)
	GetTreeByIDResponse(ctx context.Context, id string, withSensorData bool) (*tree.TreeSensorData, error)
	GetTreePredictionResponse(ctx context.Context, treeID string, withSensorData bool) (*tree.TreeSensorPrediction, error)
}

type AuthService interface {
	Service
	Register(ctx context.Context, user *auth.RegisterUser) (*auth.User, error)
}

type Service interface {
	Ready() bool
}

type Services struct {
	InfoService InfoService
	MqttService MqttService
	TreeService TreeService
	AuthService AuthService
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
