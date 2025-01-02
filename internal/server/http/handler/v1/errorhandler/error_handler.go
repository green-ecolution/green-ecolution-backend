package errorhandler

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

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

func HandleError(err error, contexts ...string) *fiber.Error {
	if err == nil {
		return nil
	}
	code := fiber.StatusInternalServerError

	// Capture runtime information
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	} else {
		// Trim the path to start from green-ecolution-
		baseMarker := "green-ecolution-"
		if idx := strings.Index(file, baseMarker); idx != -1 {
			file = file[idx:]
		}
	}

	// Use the first context if provided
	context := "No context provided"
	if len(contexts) > 0 {
		context = contexts[0]
	}

	// Classify the error
	var errorType string
	if svcErr, ok := err.(*service.Error); ok {
		errorType = "ServiceError"
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
			slog.Debug("Missing service error code", "code", svcErr.Code, "file", file, "line", line, "context", context)
		}
	} else if isNotFoundError(err) {
		errorType = "NotFoundError"
		code = fiber.StatusNotFound
	} else {
		errorType = "UnexpectedError"
	}

	// Check for specific not found errors
	for _, notFoundErr := range notFoundErrors {
		if errors.Is(err, notFoundErr) {
			code = fiber.StatusNotFound
			break
		}
	}

	// Log the error
	slog.Error(fmt.Sprintf("[%s] %s: %s", errorType, context, err.Error()), "file", file, "line", line)

	// Return error
	return fiber.NewError(code, fmt.Sprintf("[%s] %s: %s", errorType, context, err.Error()))
}

// isNotFoundError check not found errors.
func isNotFoundError(err error) bool {
	for _, notFoundErr := range notFoundErrors {
		if errors.Is(err, notFoundErr) {
			return true
		}
	}
	return false
}
