package tree_test

import (
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	httpEntities "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
)

var (
	testLatitude          = 9.446741
	testLongitude         = 54.801539
	TestTreeUpdateRequest = (*httpEntities.TreeUpdateRequest)(getMockTreeRequest("Updated description"))
	TestTreeCreateRequest = getMockTreeRequest("Created description")
	TestTrees             = []*entities.Tree{
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
		},
		{
			ID:           2,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Species:      "Pine",
			Number:       "T002",
			Latitude:     testLatitude,
			Longitude:    testLongitude,
			Description:  "A young pine tree",
			PlantingYear: 2023,
		},
	}
)

func getMockTreeRequest(description string) *httpEntities.TreeCreateRequest {
	return &httpEntities.TreeCreateRequest{
		TreeClusterID: nil,
		PlantingYear:  2023,
		Species:       "Oak",
		Number:        "T001",
		Latitude:      testLatitude,
		Longitude:     testLongitude,
		SensorID:      nil,
		Description:   description,
	}
}
