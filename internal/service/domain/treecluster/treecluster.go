package treecluster

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

type TreeClusterService struct {
	treeClusterRepo storage.TreeClusterRepository
	regionRepo      storage.RegionRepository
	treeRepo        storage.TreeRepository
}

func NewTreeClusterService(treeClusterRepo storage.TreeClusterRepository, regionRepo storage.RegionRepository, treeRepo storage.TreeRepository) service.TreeClusterService {
	return &TreeClusterService{
		treeClusterRepo: treeClusterRepo,
		regionRepo:      regionRepo,
		treeRepo:        treeRepo,
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
	region, err := s.fetchRegionByID(ctx, req.Region)
	if err != nil {
		return nil, err
	}

	for _, tree := range req.TreeIDs {
		fmt.Printf("%v", *tree.ID)
	}

	entityFunc := func(tc *domain.TreeCluster) {
		tc.WateringStatus = domain.TreeClusterWateringStatus(req.WateringStatus)
		tc.Address = req.Address
		tc.Region = region
		tc.Description = req.Description
		tc.SoilCondition = domain.TreeSoilCondition(req.SoilCondition)
		//tc.Trees = trees
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

func (s *TreeClusterService) fetchRegionByID(ctx context.Context, regionID string) (*domain.Region, error) {
	id, err := strconv.Atoi(regionID)
	if err != nil {
		return nil, service.NewError(service.BadRequest, "invalid region ID")
	}

	region, err := s.regionRepo.GetByID(ctx, int32(id))
	if err != nil {
		return nil, service.NewError(service.BadRequest, err.Error())
	}

	if region == nil {
		return nil, service.NewError(service.NotFound, storage.ErrRegionNotFound.Error())
	}

	return region, nil
}

func (s *TreeClusterService) Ready() bool {
	return s.treeClusterRepo != nil
}
