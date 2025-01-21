package region

import (
	"context"
	"errors"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
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
	log := logger.GetLogger(ctx)
	regions, err := s.regionRepo.GetAll(ctx)
	if err != nil {
		log.Error("failed to get region by id", "error", err)
		return nil, service.NewError(service.InternalError, err.Error())
	}

	return regions, nil
}

func (s *RegionService) GetByID(ctx context.Context, id int32) (*domain.Region, error) {
	log := logger.GetLogger(ctx)
	region, err := s.regionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrEntityNotFound) {
			log.Debug("region with requested id does not exits", "error", err, "region_id", id)
			return nil, service.NewError(service.NotFound, storage.ErrRegionNotFound.Error())
		}
		log.Error("failed to get region by id", "error", err, "region_id", id)
		return nil, service.NewError(service.InternalError, err.Error())
	}

	if region == nil {
		return nil, service.NewError(service.NotFound, storage.ErrRegionNotFound.Error())

	}

	return region, nil
}

func (s *RegionService) Ready() bool {
	return s.regionRepo != nil
}
