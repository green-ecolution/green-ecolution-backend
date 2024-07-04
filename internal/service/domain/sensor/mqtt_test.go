package sensor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	storageMock "github.com/SmartCityFlensburg/green-space-management/internal/storage/_mock"
)

func TestNewSensorService(t *testing.T) {
  repo := storageMock.NewMockSensorRepository(t)
	t.Run("should create a new service", func(t *testing.T) {
		svc := NewSensorService(repo)
		assert.NotNil(t, svc)
	})
}

// other test cases
