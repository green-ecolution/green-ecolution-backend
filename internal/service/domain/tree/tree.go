package tree

import (
	"context"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"log"
)

type TreeService struct {
	treeRepo   storage.TreeRepository
	sensorRepo storage.SensorRepository
}

func NewTreeService(repoTree storage.TreeRepository, repoSensor storage.SensorRepository) *TreeService {
	return &TreeService{
		treeRepo:   repoTree,
		sensorRepo: repoSensor,
	}
}

// func handleError(err error) error {
// 	if errors.Is(err, storage.ErrMongoDataNotFound) {
// 		return service.NewError(service.NotFound, err.Error())
// 	}
//
// 	return service.NewError(service.InternalError, err.Error())
// }

func (s *TreeService) Ready() bool {
	return s.treeRepo != nil && s.sensorRepo != nil
}

// ImportTree takes the rows and creates trees in the database
func (s *TreeService) ImportTree(ctx context.Context, rows [][]string) error {
	//TODO: implement the logic to import the entries of the csv file into the database
	for i, row := range rows {
		log.Printf("Row %d: %v", i+1, row)

		for j, cell := range row {
			log.Printf("Row %d, Column %d: %s", i+1, j+1, cell)
			// Process each cell here
		}
	}

	/*	// Call the repository Create function using functional options
		_, err = s.treeRepo.Create(ctx,
			tree.WithSpecies(row[3]),
			tree.WithTreeNumber(int32(number)),
			tree.WithLatitude(latitude),
			tree.WithLongitude(longitude),
			tree.WithPlantingYear(int32(plantingYear)),
		)*/

	return nil
}
