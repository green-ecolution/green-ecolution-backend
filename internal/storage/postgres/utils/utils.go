package utils

import (
	"context"
	"errors"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
)

var (
	regionMapper = generated.InternalRegionRepoMapperImpl{}
	treeMapper = generated.InternalTreeRepoMapperImpl{}
)

// This function is required as soon as you want to add data to the TreeCluster object 
// from the database, e.g. the linked region or the linked trees. 
// As this function is required in different repositories, it has been outsourced.
func MapClusterFields(ctx context.Context, s store.Store, tc *entities.TreeCluster) error {
	if err := mapRegion(ctx, s, tc); err != nil {
		return s.HandleError(err)
	}

	if err := mapTrees(ctx, s, tc); err != nil {
		return s.HandleError(err)
	}

	return nil
}

func mapRegion(ctx context.Context, s store.Store, tc *entities.TreeCluster) error {
	region, err := s.GetRegionByTreeClusterID(ctx, tc.ID)
	if err != nil {
		// If region is not found, we can still return the tree cluster
		if !errors.Is(err, storage.ErrRegionNotFound) {
			return err
		}
	}
	tc.Region = regionMapper.FromSql(region)

	return nil
}

func mapTrees(ctx context.Context, s store.Store, tc *entities.TreeCluster) error {
	trees, err := s.GetLinkedTreesByTreeClusterID(ctx, tc.ID)
	if err != nil {
		return err
	}
	tc.Trees = treeMapper.FromSqlList(trees)

	return nil
}
