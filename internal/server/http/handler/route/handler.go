package route

import (
	"github.com/SmartCityFlensburg/green-space-management/internal/service"
	"github.com/SmartCityFlensburg/green-space-management/internal/service/domain/route"
	_ "github.com/SmartCityFlensburg/green-space-management/internal/service/entities/tree"
	"github.com/gofiber/fiber/v2"
)

type RouteResponse struct {
	Points        []route.Point `json:"points"`
	TotalDistance float64         `json:"total_distance"`
} //@Name RouteResponse

//	@Summary		Get all routes
//	@Description	Get all routes
//	@Id				get-all-routes
//	@Tags			Route
//	@Produce		json
//	@Success		200	{object}	route.RouteResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/route [get]
func GetRouteAll(svc service.RouteService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		points, totalDistance := svc.NearestNeighbor(svc.DemoPoints())
    response := RouteResponse{
      Points: points,
      TotalDistance: totalDistance,
    }
    return c.JSON(response)
	}
}
