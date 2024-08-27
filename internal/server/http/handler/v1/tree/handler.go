package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all trees
// @Description	Get all trees
// @Id				get-all-trees
// @Tags			Trees
// @Produce		json
// @Param			sensor_data	query		boolean	false	"Get raw sensor data for each tree"
// @Success		200			{object}	[]tree.TreeSensorDataResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/v1/tree [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllTree(_ service.TreeService) fiber.Handler {
	// var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		// TODO: Implement GetAllTree
		return c.JSON([]tree.TreeSensorDataResponse{})
	}
}

// @Summary		Get tree by ID
// @Description	Get tree by ID
// @Id				get-tree-by-id
// @Tags			Trees
// @Produce		json
// @Param			treeID		path		string	true	"Tree ID"
// @Param			sensor_data	query		boolean	false	"Get raw sensor data for each tree"
// @Success		200			{object}	tree.TreeSensorDataResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/v1/tree/{treeID} [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreeByID(_ service.TreeService) fiber.Handler {
	// var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		// TODO: Implement GetTreeByID
		return c.JSON(tree.TreeSensorDataResponse{})
	}
}

// @Summary		Get tree prediction by tree ID
// @Description	Get tree prediction by tree ID
// @Id				get-tree-prediction-by-id
// @Tags			Trees
// @Produce		json
// @Param			treeID		path		string	true	"Tree ID"
// @Param			sensor_data	query		boolean	false	"Get raw sensor data for each tree"
// @Success		200			{object}	tree.TreeSensorPredictionResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/v1/tree/{treeID}/prediction [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetTreePredictions(_ service.TreeService) fiber.Handler {
	// var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		// TODO: Implement GetTreePredictions
		return c.JSON(tree.TreeSensorPredictionResponse{})
	}
}
