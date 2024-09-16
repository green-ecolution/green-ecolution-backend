package role

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get all user roles
// @Description	Get all user roles
// @Tags			Role
// @Produce		json
// @Success		200		{object}	entities.RoleListResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			page	query		string	false	"Page"
// @Param			limit	query		string	false	"Limit"
// @Router			/v1/role [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAllUserRoles(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

// @Summary		Get a role by ID
// @Description	Get a role by ID
// @Tags			Role
// @Produce		json
// @Success		200	{object}	entities.RoleResponse
// @Failure		400	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Param			id	path		string	true	"Role ID"
// @Router			/v1/role/{id} [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetRoleByID(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}
