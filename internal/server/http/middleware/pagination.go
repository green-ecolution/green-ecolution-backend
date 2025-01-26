package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

const (
	Page  = "page"
	Limit = "limit"
)

func PaginationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defaultPage := 1
		defaultLimit := -1

		pageParam := c.Query(Page, strconv.Itoa(defaultPage))
		page, err := strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			return service.NewError(service.BadRequest, "invalid page format")
		}

		limitParam := c.Query(Limit, strconv.Itoa(defaultLimit))
		limit, err := strconv.Atoi(limitParam)
		if err != nil || (limit != -1 && limit <= 0) {
			return service.NewError(service.BadRequest, "invalid limit format")
		}

		c.Locals(Page, int32(page))
		c.Locals(Limit, int32(limit))

		return c.Next()
	}
}
