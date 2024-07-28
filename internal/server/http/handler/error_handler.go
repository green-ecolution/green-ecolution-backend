package handler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func HandleError(err error) *fiber.Error {
	code := fiber.StatusInternalServerError
	var svcErr service.Error
	if errors.As(err, &svcErr) {
		switch svcErr.Code {
		case service.NotFound:
			code = fiber.StatusNotFound
		case service.BadRequest:
			code = fiber.StatusBadRequest
		case service.Forbidden:
			code = fiber.StatusForbidden
		case service.Unauthorized:
			code = fiber.StatusUnauthorized
		case service.InternalError:
			code = fiber.StatusInternalServerError
		default:
      slog.Debug("missing service error code", "code", svcErr.Code)
		}
	}
	return fiber.NewError(code, err.Error())
}
