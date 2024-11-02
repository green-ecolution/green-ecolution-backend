package tree

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/utils"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
)

var (
	testLatitude     = 9.446741
	testLongitude    = 54.801539
	TestTreeClusters = []*entities.TreeCluster{
		{
			ID:            1,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "Cluster 1",
			Address:       "123 Main St",
			Description:   "Test description",
			SoilCondition: entities.TreeSoilConditionLehmig,
			Archived:      false,
			Latitude:      utils.P(testLatitude),
			Longitude:     utils.P(testLongitude),
			Trees:         TestTreesList,
		},
		{
			ID:            2,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Name:          "Cluster 2",
			Address:       "456 Second St",
			Description:   "Test description",
			SoilCondition: entities.TreeSoilConditionSandig,
			Archived:      false,
			Latitude:      nil,
			Longitude:     nil,
			Trees:         []*entities.Tree{},
		},
	}

	TestTreesList = []*entities.Tree{
		{
			ID:           1,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Species:      "Oak",
			Number:       "T001",
			Latitude:     testLatitude,
			Longitude:    testLongitude,
			Description:  "A mature oak tree",
			PlantingYear: 2023,
			Readonly:     true,
		},
		{
			ID:           2,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Species:      "Pine",
			Number:       "T002",
			Latitude:     9.446700,
			Longitude:    54.801510,
			Description:  "A young pine tree",
			PlantingYear: 2023,
			Readonly:     true,
		},
	}

	TestSensors = []*entities.Sensor{
		{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    entities.SensorStatusUnknown,
			Data:      nil,
		},
		{
			ID:        2,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Status:    entities.SensorStatusUnknown,
			Data:      nil,
		},
	}

	TestTreeCreate = &entities.TreeCreate{
		Species:       "Oak",
		Latitude:      testLatitude,
		Longitude:     testLongitude,
		PlantingYear:  2023,
		Number:        "T001",
		TreeClusterID: utils.PtrInt32(1),
		SensorID:      utils.PtrInt32(1),
	}

	TestTreeImport = &entities.TreeImport{
		Latitude:     testLatitude,
		Longitude:    testLongitude,
		PlantingYear: 2023,
		Species:      "Oak",
		Number:       "T001",
	}

	TestTreeUpdate = &entities.TreeUpdate{
		TreeClusterID: utils.PtrInt32(1),
		SensorID:      utils.PtrInt32(1),
		PlantingYear:  2023,
		Species:       "Oak",
		Number:        "T001",
		Latitude:      testLatitude,
		Longitude:     testLongitude,
		Description:   "Updated description",
	}
)
