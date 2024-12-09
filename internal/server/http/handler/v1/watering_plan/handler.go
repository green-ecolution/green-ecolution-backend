package wateringplan

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

var (
	wateringPlanMapper = generated.WateringPlanHTTPMapperImpl{}
)

//	@Summary		Get all watering plans
//	@Description	Get all watering plans
//	@Tags			Watering Plan
//	@Produce		json
//	@Success		200		{object}	entities.WateringPlanListResponse
//	@Failure		400		{object}	HTTPError
//	@Failure		500		{object}	HTTPError
//	@Param			page	query		string	false	"Page"
//	@Param			limit	query		string	false	"Limit"
//	@Router			/v1/watering-plan [get]
//	@Security		Keycloak
func GetAllWateringPlans(svc service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		domainData, err := svc.GetAll(ctx)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := make([]*entities.WateringPlanResponse, len(domainData))
		for i, domain := range domainData {
			data[i] = mapWateringPlanToDto(domain)
		}

		return c.JSON(entities.WateringPlanListResponse{
			Data:       data,
			Pagination: &entities.Pagination{}, // TODO: Handle pagination
		})
	}
}

//	@Summary		Get a watering plan by ID
//	@Description	Get a watering plan by ID
//	@Tags			Watering Plan
//	@Produce		json
//	@Success		200	{object}	entities.WateringPlanResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Param			id	path		string	true	"Watering plan ID"
//	@Router			/v1/watering-plan/{id} [get]
//	@Security		Keycloak
func GetWateringPlanByID(svc service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("watering_plan_id"))
		if err != nil {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		domainData, err := svc.GetByID(ctx, int32(id))

		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapWateringPlanToDto(domainData)

		return c.JSON(data)
	}
}

//	@Summary		Create watering plan
//	@Description	Create watering plan
//	@Id				create-watering-plan
//	@Tags			Watering Plan
//	@Produce		json
//	@Success		201	{object}	entities.WateringPlanResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/watering-plan [post]
//	@Param			body	body	entities.WateringPlanCreateRequest	true	"Watering Plan Create Request"
//	@Security		Keycloak
func CreateWateringPlan(svc service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		var req entities.WateringPlanCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainReq := wateringPlanMapper.FromCreateRequest(&req)
		domainData, err := svc.Create(ctx, domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapWateringPlanToDto(domainData)
		return c.Status(fiber.StatusCreated).JSON(data)
	}
}

//	@Summary		Update watering plan
//	@Description	Update watering plan
//	@Id				update-watering-plan
//	@Tags			Watering Plan
//	@Produce		json
//	@Success		200	{object}	entities.WateringPlanResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/watering-plan/{watering_plan_id} [put]
//	@Param			watering_plan_id	path	string								true	"Watering Plan ID"
//	@Param			body				body	entities.WateringPlanUpdateRequest	true	"Watering Plan Update Request"
//	@Security		Keycloak
func UpdateWateringPlan(_ service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

//	@Summary		Delete watering plan
//	@Description	Delete watering plan
//	@Id				delete-watering-plan
//	@Tags			Watering Plan
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/watering-plan/{watering_plan_id} [delete]
//	@Param			watering_plan_id	path	string	true	"Watering Plan ID"
//	@Security		Keycloak
func DeleteWateringPlan(_ service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func mapWateringPlanToDto(wp *domain.WateringPlan) *entities.WateringPlanResponse {
	dto := wateringPlanMapper.FromResponse(wp)

	return dto
}
