package evaluation

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get evaluation
// @Description	Get evaluation values
// @Id				get-evaluation
// @Tags			Evaluation
// @Produce		json
// @Success		200	{object}	entities.EvaluationResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/evaluation [get]
// @Security		Keycloak
func GetEvaluation(svc service.EvaluationService) fiber.Handler {
	var m mapper.EvaluationHTTPMapper = &mapper.EvaluationHTTPMapper{}

	return func(c *fiber.Ctx) error {
		domainInfo, err := svc.GetEvaluation(c.Context())
		if err != nil {
			return errorhandler.HandleError(err)
		}
		return c.JSON(m.FromResponse(domainInfo))
	}
}
