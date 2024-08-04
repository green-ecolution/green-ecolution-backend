package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/error_handler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

//	@Summary		Get all trees
//	@Description	Get all trees
//	@Id				get-all-trees
//	@Tags			Trees,v1
//	@Produce		json
//	@Param			sensor_data	query		boolean	false	"Get raw sensor data for each tree"
//	@Success		200			{object}	[]tree.TreeSensorDataResponse
//	@Failure		400			{object}	HTTPError
//	@Failure		401			{object}	HTTPError
//	@Failure		403			{object}	HTTPError
//	@Failure		404			{object}	HTTPError
//	@Failure		500			{object}	HTTPError
//	@Router			/v1/tree [get]
func GetAllTree(svc service.TreeService) fiber.Handler {
	var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		domainTree, err := svc.GetAllTreesResponse(c.Context(), c.QueryBool("sensor_data"))
		if err != nil {
			return error_handler.HandleError(err)
		}

		response := mapper.ToTreeSensorDataResponseList(domainTree)
		return c.JSON(response)
	}
}

//	@Summary		Get tree by ID
//	@Description	Get tree by ID
//	@Id				get-tree-by-id
//	@Tags			Trees,v1
//	@Produce		json
//	@Param			treeID		path		string	true	"Tree ID"
//	@Param			sensor_data	query		boolean	false	"Get raw sensor data for each tree"
//	@Success		200			{object}	tree.TreeSensorDataResponse
//	@Failure		400			{object}	HTTPError
//	@Failure		401			{object}	HTTPError
//	@Failure		403			{object}	HTTPError
//	@Failure		404			{object}	HTTPError
//	@Failure		500			{object}	HTTPError
//	@Router			/v1/tree/{treeID} [get]
func GetTreeByID(svc service.TreeService) fiber.Handler {
	var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		domainTree, err := svc.GetTreeByIDResponse(c.Context(), c.Params("id"), c.QueryBool("sensor_data"))
		if err != nil {
			return error_handler.HandleError(err)
		}

		response := mapper.ToTreeSensorDataResponse(domainTree)
		return c.JSON(response)
	}
}

//	@Summary		Get tree prediction by tree ID
//	@Description	Get tree prediction by tree ID
//	@Id				get-tree-prediction-by-id
//	@Tags			Trees,v1
//	@Produce		json
//	@Param			treeID		path		string	true	"Tree ID"
//	@Param			sensor_data	query		boolean	false	"Get raw sensor data for each tree"
//	@Success		200			{object}	tree.TreeSensorPredictionResponse
//	@Failure		400			{object}	HTTPError
//	@Failure		401			{object}	HTTPError
//	@Failure		403			{object}	HTTPError
//	@Failure		404			{object}	HTTPError
//	@Failure		500			{object}	HTTPError
//	@Router			/v1/tree/{treeID}/prediction [get]
func GetTreePredictions(svc service.TreeService) fiber.Handler {
	var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		domainTree, err := svc.GetTreePredictionResponse(c.Context(), c.Params("id"), c.QueryBool("sensor_data"))
		if err != nil {
			return error_handler.HandleError(err)
		}

		response := mapper.ToTreeSensorPredictionResponse(domainTree)
		return c.JSON(response)
	}
}
