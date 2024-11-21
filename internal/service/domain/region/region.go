package region

import (
	"context"
	"errors"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type RegionService struct {
	regionRepo storage.RegionRepository
}

func NewRegionService(regionRepository storage.RegionRepository) service.RegionService {
	return &RegionService{
		regionRepo: regionRepository,
	}
}

func (s *RegionService) GetAll(ctx context.Context) ([]*domain.Region, error) {
	regions, err := s.regionRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return regions, nil
}

func (s *RegionService) GetByID(ctx context.Context, id int32) (*domain.Region, error) {
	region, err := s.regionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return region, nil
}

func (s *RegionService) Ready() bool {
	return s.regionRepo != nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, storage.ErrRegionNotFound.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}
