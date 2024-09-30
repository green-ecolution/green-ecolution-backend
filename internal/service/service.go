package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
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
	GetAppInfo(context.Context) (*domain.App, error)
	GetAppInfoResponse(context.Context) (*domain.App, error)
}

type MqttService interface {
	Service
	HandleMessage(ctx context.Context, payload *domain.MqttPayload) (*domain.MqttPayload, error)
	SetConnected(bool)
}

type TreeService interface {
	Service
	ImportTree(ctx context.Context, trees []*domain.Tree) error

	GetAll(ctx context.Context) ([]*domain.Tree, error)
	GetByID(ctx context.Context, id int) (*domain.Tree, error)
	// TODO: Implement Functions
}

type AuthService interface {
	Service
	LoginRequest(ctx context.Context, loginRequest *domain.LoginRequest) (*domain.LoginResp, error)
	LogoutRequest(ctx context.Context, logoutRequest *domain.Logout) error
	ClientTokenCallback(ctx context.Context, loginCallback *domain.LoginCallback) (*domain.ClientToken, error)
	Register(ctx context.Context, user *domain.RegisterUser) (*domain.User, error)
	RetrospectToken(ctx context.Context, token string) (*domain.IntroSpectTokenResult, error)
}

type RegionService interface {
	Service
	GetAll(ctx context.Context) ([]*domain.Region, error)
	GetByID(ctx context.Context, id int32) (*domain.Region, error)
}

type Service interface {
	Ready() bool
}

type Services struct {
	InfoService   InfoService
	MqttService   MqttService
	TreeService   TreeService
	AuthService   AuthService
	RegionService RegionService
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
