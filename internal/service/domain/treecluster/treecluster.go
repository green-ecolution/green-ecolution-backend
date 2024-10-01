package treecluster

import (
	"context"
	"errors"

	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/treecluster"
)

type TreeClusterService struct {
	treeClusterRepo storage.TreeClusterRepository
	treeRepo        storage.TreeRepository
	regionRepo      storage.RegionRepository
}

func NewTreeClusterService(treeClusterRepo storage.TreeClusterRepository, treeRepo storage.TreeRepository, regionRepo storage.RegionRepository) service.TreeClusterService {
	return &TreeClusterService{
		treeClusterRepo: treeClusterRepo,
		treeRepo:        treeRepo,
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

func (s *TreeClusterService) Create(ctx context.Context, tc *domain.TreeClusterCreate) (*domain.TreeCluster, error) {
	treeIDs := make([]int32, len(tc.TreeIDs))
	fn := make([]domain.EntityFunc[domain.TreeCluster], 0)

	if len(tc.TreeIDs) > 0 {
		for i, id := range tc.TreeIDs {
			treeIDs[i] = *id
		}
		geomFn, err := s.prepareGeom(ctx, treeIDs)
		if err != nil {
			return nil, err
		}
		fn = append(fn, geomFn...)
	}

	fn = append(fn, treecluster.WithName(tc.Name), treecluster.WithAddress(tc.Address), treecluster.WithDescription(tc.Description))

	c, err := s.treeClusterRepo.Create(ctx, fn...)
	if err != nil {
		return nil, handleError(err)
	}

	if err = s.treeRepo.UpdateTreeClusterID(ctx, treeIDs, &c.ID); err != nil {
		return nil, handleError(err)
	}

	trees, err := s.treeRepo.GetByTreeClusterID(ctx, c.ID)
	if err != nil {
		return nil, handleError(err)
	}

	c.Trees = trees

	return c, nil
}

func (s *TreeClusterService) Update(ctx context.Context, id int32, tc *domain.TreeClusterUpdate) (*domain.TreeCluster, error) {
	treeIDs := make([]int32, len(tc.TreeIDs))
	fn := make([]domain.EntityFunc[domain.TreeCluster], 0)

	// TODO: Add a transaction to undo this change if an error occurs.
	if err := s.treeRepo.UnlinkTreeClusterID(ctx, id); err != nil {
		return nil, handleError(err)
	}

	if len(tc.TreeIDs) > 0 {
		for i, id := range tc.TreeIDs {
			treeIDs[i] = *id
		}
		geomFn, err := s.prepareGeom(ctx, treeIDs)
		if err != nil {
			return nil, err
		}
		fn = append(fn, geomFn...)
	}

	fn = append(fn,
		treecluster.WithName(tc.Name),
		treecluster.WithAddress(tc.Address),
		treecluster.WithDescription(tc.Description),
		treecluster.WithArchived(tc.Archived),
		treecluster.WithSoilCondition(tc.SoilCondition),
	)

	c, err := s.treeClusterRepo.Update(ctx, id, fn...)
	if err != nil {
		return nil, handleError(err)
	}

	if err = s.treeRepo.UpdateTreeClusterID(ctx, treeIDs, &c.ID); err != nil {
		return nil, handleError(err)
	}

	return c, nil
}

func (s *TreeClusterService) Delete(ctx context.Context, id int32) error {
	_, err := s.treeClusterRepo.GetByID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	// TODO: Add a transaction to undo this change if an error occurs.
	err = s.treeRepo.UnlinkTreeClusterID(ctx, id)
	if err != nil {
		return handleError(err)
	}

	err = s.treeClusterRepo.Delete(ctx, id)
	if err != nil {
		return handleError(err)
	}

	return nil
}

func (s *TreeClusterService) Ready() bool {
	return s.treeClusterRepo != nil
}

func handleError(err error) error {
	if errors.Is(err, storage.ErrEntityNotFound) {
		return service.NewError(service.NotFound, err.Error())
	}

	return service.NewError(service.InternalError, err.Error())
}

func (s *TreeClusterService) prepareGeom(ctx context.Context, treeIDs []int32) ([]domain.EntityFunc[domain.TreeCluster], error) {
	lat, long, err := s.calculateCenterPoint(ctx, treeIDs)
	if err != nil {
		return nil, err
	}

	region, err := s.getRegionByPoint(ctx, lat, long)
	if err != nil {
		return nil, err
	}

	fn := []domain.EntityFunc[domain.TreeCluster]{
		treecluster.WithLatitude(&lat),
		treecluster.WithLongitude(&long),
		treecluster.WithRegion(region),
	}

	return fn, nil
}

func (s *TreeClusterService) calculateCenterPoint(ctx context.Context, treeIDs []int32) (lat, long float64, err error) {
	lat, long, err = s.treeRepo.GetCenterPoint(ctx, treeIDs)
	if err != nil {
		return 0, 0, err
	}

	return lat, long, nil
}

func (s *TreeClusterService) getRegionByPoint(ctx context.Context, lat, long float64) (*domain.Region, error) {
	region, err := s.regionRepo.GetByPoint(ctx, lat, long)
	if err != nil {
		return nil, handleError(err)
	}

	return region, nil
}
