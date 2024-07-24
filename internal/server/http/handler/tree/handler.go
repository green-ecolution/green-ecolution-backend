package tree

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/tree/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all trees
// @Description	Get all trees
// @Id				get-all-trees
// @Tags			Trees
// @Produce		json
// @Param			sensor_data	query		boolean	false	"Get raw sensorSQL data for each treeSQL"
// @Success		200			{object}	[]treeSQL.TreeSensorDataResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/treeSQL [get]
func GetAllTree(svc service.TreeService) fiber.Handler {
	var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		domainTree, err := svc.GetAllTreesResponse(c.Context(), c.QueryBool("sensor_data"))
		if err != nil {
			return handler.HandleError(err)
		}

		response := mapper.ToTreeSensorDataResponseList(domainTree)
		return c.JSON(response)
	}
}

// @Summary		Get treeSQL by ID
// @Description	Get treeSQL by ID
// @Id				get-treeSQL-by-id
// @Tags			Trees
// @Produce		json
// @Param			treeID		path		string	true	"Tree ID"
// @Param			sensor_data	query		boolean	false	"Get raw sensorSQL data for each treeSQL"
// @Success		200			{object}	treeSQL.TreeSensorDataResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/treeSQL/{treeID} [get]
func GetTreeByID(svc service.TreeService) fiber.Handler {
	var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		domainTree, err := svc.GetTreeByIDResponse(c.Context(), c.Params("id"), c.QueryBool("sensor_data"))
		if err != nil {
			return handler.HandleError(err)
		}

		response := mapper.ToTreeSensorDataResponse(domainTree)
		return c.JSON(response)
	}
}

// @Summary		Get treeSQL prediction by treeSQL ID
// @Description	Get treeSQL prediction by treeSQL ID
// @Id				get-treeSQL-prediction-by-id
// @Tags			Trees
// @Produce		json
// @Param			treeID		path		string	true	"Tree ID"
// @Param			sensor_data	query		boolean	false	"Get raw sensorSQL data for each treeSQL"
// @Success		200			{object}	treeSQL.TreeSensorPredictionResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/treeSQL/{treeID}/prediction [get]
func GetTreePredictions(svc service.TreeService) fiber.Handler {
	var mapper tree.TreeHTTPMapper = &generated.TreeHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		domainTree, err := svc.GetTreePredictionResponse(c.Context(), c.Params("id"), c.QueryBool("sensor_data"))
		if err != nil {
			return handler.HandleError(err)
		}

		response := mapper.ToTreeSensorPredictionResponse(domainTree)
		return c.JSON(response)
	}
}
