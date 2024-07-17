package info

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	infoResponse "github.com/green-ecolution/green-ecolution-backend/internal/service/entities/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type InfoService struct {
	infoRepository storage.InfoRepository
	mapper         mapper.InfoMapper
}

func NewInfoService(infoRepository storage.InfoRepository) *InfoService {
	return &InfoService{
		infoRepository: infoRepository,
		mapper:         &generated.InfoMapperImpl{},
	}
}

func (s *InfoService) GetAppInfo(ctx context.Context) (*info.App, error) {
	appInfo, err := s.infoRepository.GetAppInfo(ctx)
	if err != nil {
		if errors.Is(err, storage.ErrIpNotFound) {
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

	return s.mapper.FromEntity(appInfo), nil
}

func (s *InfoService) GetAppInfoResponse(ctx context.Context) (*infoResponse.AppInfoResponse, error) {
	appInfo, err := s.GetAppInfo(ctx)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToResponse(appInfo), nil
}

func (s *InfoService) Ready() bool {
	return s.infoRepository != nil
}
