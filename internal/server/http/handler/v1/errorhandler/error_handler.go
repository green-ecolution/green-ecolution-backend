package errorhandler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
)

var notFoundErrors = []error{
	storage.ErrRegionNotFound,
	storage.ErrIDNotFound,
	storage.ErrEntityNotFound,
	storage.ErrSensorNotFound,
	storage.ErrImageNotFound,
	storage.ErrFlowerbedNotFound,
	storage.ErrTreeClusterNotFound,
	storage.ErrTreeNotFound,
	storage.ErrVehicleNotFound,
	storage.ErrWateringPlanNotFound,
}

func HandleError(err error) *fiber.Error {
	code := fiber.StatusInternalServerError

	// Check if the error is of type service.Error
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

	// Check for specific "not found" errors
	for _, notFoundErr := range notFoundErrors {
		if errors.Is(err, notFoundErr) {
			code = fiber.StatusNotFound
			break
		}
	}

	return fiber.NewError(code, err.Error())
}
