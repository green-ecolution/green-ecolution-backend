package wateringplan

import (
	"fmt"
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

// @Summary		Get all watering plans
// @Description	Get all watering plans
// @Tags			Watering Plan
// @Produce		json
// @Success		200		{object}	entities.WateringPlanListResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			page	query		string	false	"Page"
// @Param			limit	query		string	false	"Limit"
// @Router			/v1/watering-plan [get]
// @Security		Keycloak
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

// @Summary		Get watering plan by ID
// @Description	Get watering plan by ID
// @Id				get-watering-plan-by-id
// @Tags			Watering Plan
// @Produce		json
// @Success		200	{object}	entities.WateringPlanResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/watering-plan/{id} [get]
// @Param			id	path	string	true	"Watering Plan ID"
// @Security		Keycloak
func GetWateringPlanByID(svc service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		fmt.Println("-----")
		fmt.Println(strconv.Atoi(c.Params("id")))
		id, err := strconv.Atoi(c.Params("id"))
		fmt.Print(id)
		fmt.Print(err)
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

// @Summary		Create watering plan
// @Description	Create watering plan
// @Id				create-watering-plan
// @Tags			Watering Plan
// @Produce		json
// @Success		201	{object}	entities.WateringPlanResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/watering-plan [post]
// @Param			body	body	entities.WateringPlanCreateRequest	true	"Watering Plan Create Request"
// @Security		Keycloak
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

// @Summary		Update watering plan
// @Description	Update watering plan
// @Id				update-watering-plan
// @Tags			Watering Plan
// @Produce		json
// @Success		200	{object}	entities.WateringPlanResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/watering-plan/{id} [put]
// @Param			id		path	string								true	"Watering Plan ID"
// @Param			body	body	entities.WateringPlanUpdateRequest	true	"Watering Plan Update Request"
// @Security		Keycloak
func UpdateWateringPlan(svc service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		var req entities.WateringPlanUpdateRequest
		if err = c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		domainReq := wateringPlanMapper.FromUpdateRequest(&req)
		domainData, err := svc.Update(ctx, int32(id), domainReq)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		data := mapWateringPlanToDto(domainData)
		return c.JSON(data)
	}
}

// @Summary		Delete watering plan
// @Description	Delete watering plan
// @Id				delete-watering-plan
// @Tags			Watering Plan
// @Produce		json
// @Success		204
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/watering-plan/{id} [delete]
// @Param			id	path	string	true	"Watering Plan ID"
// @Security		Keycloak
func DeleteWateringPlan(svc service.WateringPlanService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			err := service.NewError(service.BadRequest, "invalid ID format")
			return errorhandler.HandleError(err)
		}

		err = svc.Delete(ctx, int32(id))
		if err != nil {
			return errorhandler.HandleError(err)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func mapWateringPlanToDto(wp *domain.WateringPlan) *entities.WateringPlanResponse {
	dto := wateringPlanMapper.FromResponse(wp)

	return dto
}
