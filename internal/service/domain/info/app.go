package info

import (
	"context"
	"errors"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type InfoService struct {
	infoRepository storage.InfoRepository
}

func NewInfoService(infoRepository storage.InfoRepository) *InfoService {
	return &InfoService{
		infoRepository: infoRepository,
	}
}

func (s *InfoService) GetAppInfo(ctx context.Context) (*domain.App, error) {
	appInfo, err := s.infoRepository.GetAppInfo(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrIPNotFound) {
			return nil, service.NewError(service.InternalError, err.Error())
		}
		if errors.Is(err, storage.ErrIFacesNotFound) {
			return nil, service.NewError(service.InternalError, err.Error())
		}
		if errors.Is(err, storage.ErrIFacesAddressNotFound) {
			return nil, service.NewError(service.InternalError, err.Error())
		}
		if errors.Is(err, storage.ErrHostnameNotFound) {
			return nil, service.NewError(service.InternalError, err.Error())
		}
	}

	return appInfo, nil
}

func (s *InfoService) GetAppInfoResponse(ctx context.Context) (*domain.App, error) {
	appInfo, err := s.GetAppInfo(ctx)
	if err != nil {
		return nil, err
	}

	return appInfo, nil
}

func (s *InfoService) Ready() bool {
	return s.infoRepository != nil
}
