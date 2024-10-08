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
	fn := make([]domain.EntityFunc[domain.TreeCluster], 0)
	treeFn, err := s.prepareTrees(ctx, tc.TreeIDs)
	if err != nil {
		return nil, err
	}

	fn = append(fn, treeFn...)
	fn = append(fn, treecluster.WithName(tc.Name), treecluster.WithAddress(tc.Address), treecluster.WithDescription(tc.Description))

	c, err := s.treeClusterRepo.Create(ctx, fn...)
	if err != nil {
		return nil, handleError(err)
	}

	return c, nil
}

func (s *TreeClusterService) Update(ctx context.Context, id int32, tc *domain.TreeClusterUpdate) (*domain.TreeCluster, error) {
	fn := make([]domain.EntityFunc[domain.TreeCluster], 0)
	treeFn, err := s.prepareTrees(ctx, tc.TreeIDs)
	if err != nil {
		return nil, err
	}

	fn = append(fn, treeFn...)
	fn = append(fn,
		treecluster.WithName(tc.Name),
		treecluster.WithAddress(tc.Address),
		treecluster.WithDescription(tc.Description),
		treecluster.WithSoilCondition(tc.SoilCondition),
	)

	c, err := s.treeClusterRepo.Update(ctx, id, fn...)
	if err != nil {
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

func (s *TreeClusterService) prepareTrees(ctx context.Context, ids []*int32) ([]domain.EntityFunc[domain.TreeCluster], error) {
	fn := make([]domain.EntityFunc[domain.TreeCluster], 0)
	treeIDs := make([]int32, len(ids))
	for i, id := range ids {
		treeIDs[i] = *id
	}

	var err error
	trees, err := s.treeRepo.GetTreesByIDs(ctx, treeIDs)
	if err != nil {
		return nil, err
	}

	fn = append(fn, treecluster.WithTrees(trees))

	if len(trees) > 0 {
		geomFn, err := s.prepareGeom(ctx, treeIDs)
		if err != nil {
			return nil, err
		}
		fn = append(fn, geomFn...)
	} else {
		fn = append(fn, treecluster.WithLatitude(nil), treecluster.WithLongitude(nil), treecluster.WithRegion(nil))
	}

	return fn, nil
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
