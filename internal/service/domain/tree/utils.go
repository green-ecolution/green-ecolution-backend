package tree

import (
	"time"

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
			Latitude:      float64Ptr(testLatitude),
			Longitude:     float64Ptr(testLongitude),
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
			Latitude:     9.446741,
			Longitude:    54.801539,
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
)

func float64Ptr(f float64) *float64 {
	return &f
}

func ptrToInt32(i int32) *int32 {
	return &i
}
