package treecluster

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	serverEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

var (
	TestCluster = &entities.TreeCluster{
		ID:             1,
		Name:           "Test Cluster",
		Address:        "456 New St",
		Description:    "Description",
		WateringStatus: entities.WateringStatusBad,
		Region:         &entities.Region{ID: 1, Name: "Region 1"},
		Archived:       false,
		Latitude:       utils.Float64Ptr(9.446741),
		Longitude:      utils.Float64Ptr(54.801539),
		SoilCondition:  entities.TreeSoilConditionSandig,
		Trees: []*entities.Tree{
			{
				ID:           1,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Species:      "Oak",
				Number:       "T001",
				Latitude:     9.446741,
				Longitude:    54.801539,
				Description:  "A mature oak tree",
				PlantingYear: 2023,
				Readonly:     true,
			},
		},
	}

	TestClusterRequest = serverEntities.TreeClusterCreateRequest{
		Name:          "Cluster Request",
		Address:       "123 Main St",
		Description:   "Test description",
		SoilCondition: serverEntities.TreeSoilConditionSandig,
		TreeIDs:       []*int32{utils.PtrInt32(1)},
	}

	TestClusterList = []*entities.TreeCluster{
		TestCluster,
		{
			ID:             2,
			Name:           "Second Cluster",
			Address:        "789 Another St",
			Description:    "Another description",
			WateringStatus: entities.WateringStatusGood,
			Region:         &entities.Region{ID: 2, Name: "Region 2"},
			Archived:       false,
			Latitude:       utils.Float64Ptr(10.123456),
			Longitude:      utils.Float64Ptr(55.654321),
			SoilCondition:  entities.TreeSoilConditionLehmig,
			Trees: []*entities.Tree{},
		},
	}
)
