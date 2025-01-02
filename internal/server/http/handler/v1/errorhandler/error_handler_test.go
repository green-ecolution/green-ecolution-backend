package errorhandler

import (
	"errors"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	t.Run("HandleError with nil error", func(t *testing.T) {
		// when
		err := HandleError(nil)
		// validate
		assert.Nil(t, err)
	})

	t.Run("HandleError with NotFoundError and context", func(t *testing.T) {
		// when
		err := HandleError(storage.ErrEntityNotFound, "Fetching entity by ID")
		// validate
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "NotFoundError")
		assert.Contains(t, err.Error(), "Fetching entity by ID")
	})

	t.Run("HandleError with ServiceError - NotFound and context", func(t *testing.T) {
		// when
		serviceErr := &service.Error{Code: service.NotFound, Message: "Item not found"}
		err := HandleError(serviceErr, "Service operation")
		// validate
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "ServiceError")
		assert.Contains(t, err.Error(), "Service operation")
		assert.Contains(t, err.Error(), "Item not found")
	})

	t.Run("HandleError with ServiceError - BadRequest", func(t *testing.T) {
		// when
		serviceErr := &service.Error{Code: service.BadRequest, Message: "Invalid input"}
		err := HandleError(serviceErr, "Invalid request")
		// validate
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "ServiceError")
		assert.Contains(t, err.Error(), "Invalid request")
		assert.Contains(t, err.Error(), "Invalid input")
	})

	t.Run("HandleError with ServiceError - Unauthorized", func(t *testing.T) {

		serviceErr := &service.Error{Code: service.Unauthorized, Message: "Unauthorized access"}
		err := HandleError(serviceErr, "Accessing restricted resource")
		// validate
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "ServiceError")
		assert.Contains(t, err.Error(), "Accessing restricted resource")
		assert.Contains(t, err.Error(), "Unauthorized access")
	})

	t.Run("HandleError with UnexpectedError", func(t *testing.T) {
		// when
		unexpectedErr := errors.New("unexpected issue occurred")
		err := HandleError(unexpectedErr, "Processing operation")
		// validate
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "UnexpectedError")
		assert.Contains(t, err.Error(), "Processing operation")
		assert.Contains(t, err.Error(), "unexpected issue occurred")
	})

	t.Run("HandleError with empty context", func(t *testing.T) {
		// when
		unexpectedErr := errors.New("unexpected issue occurred")
		err := HandleError(unexpectedErr)
		// validate
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "UnexpectedError")
		assert.Contains(t, err.Error(), "No context provided")
		assert.Contains(t, err.Error(), "unexpected issue occurred")
	})

	t.Run("HandleError with NotFoundError in list", func(t *testing.T) {
		// when
		err := HandleError(storage.ErrSensorNotFound, "Fetching sensor data")
		// validate
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "NotFoundError")
		assert.Contains(t, err.Error(), "Fetching sensor data")
	})

}
