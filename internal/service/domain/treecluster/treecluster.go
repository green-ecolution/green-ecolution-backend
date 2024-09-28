package treecluster

import (
	"context"
	"errors"
	"strconv"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type TreeClusterService struct {
	treeClusterRepo storage.TreeClusterRepository
	regionRepo      storage.RegionRepository
}

func NewTreeClusterService(treeClusterRepo storage.TreeClusterRepository, regionRepo storage.RegionRepository) service.TreeClusterService {
	return &TreeClusterService{
		treeClusterRepo: treeClusterRepo,
		regionRepo:      regionRepo,
	}
}

func (s *TreeClusterService) GetAll(ctx context.Context) ([]*domain.TreeCluster, error) {
	treeClusters, err := s.treeClusterRepo.GetAll(ctx)
	if err != nil {
		return nil, handleError(err)
	}

	return treeClusters, nil
}

func (s *TreeClusterService) GetByID(ctx context.Context, id int32) (*domain.TreeCluster, error) {
	treeCluster, err := s.treeClusterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, handleError(err)
	}

	return treeCluster, nil
}

func (s *TreeClusterService) Create(ctx context.Context, req *entities.TreeClusterCreateRequest) (*domain.TreeCluster, error) {
	id, err := strconv.Atoi(req.Region)
	if err != nil {
		return nil, service.NewError(service.BadRequest, err.Error())
	}

	region, err := s.regionRepo.GetByID(ctx, int32(id))
	if err != nil {
		return nil, service.NewError(service.BadRequest, err.Error())
	}

	entityFunc := func(tc *domain.TreeCluster) {
		tc.WateringStatus = domain.TreeClusterWateringStatus(req.WateringStatus)
		tc.Address = req.Address
		tc.Region = region
		tc.Description = req.Description
		tc.SoilCondition = domain.TreeSoilCondition(req.SoilCondition)
	}

	createdTreeCluster, err := s.treeClusterRepo.Create(ctx, entityFunc)
	if err != nil {
		return nil, handleError(err)
	}

	return createdTreeCluster, nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeClusterService) Ready() bool {
	return s.treeClusterRepo != nil
}
