package image

import (
	"log/slog"
	"os"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/test"
)

func TestMain(m *testing.M) {
	m.Run()
	close, _, err := test.SetupPostgresContainer()
	if err != nil {
		slog.Error("Error setting up postgres container", "error", err)
		panic(err)
	}
	defer close()

	os.Exit(m.Run())
}

func TestImageRepository(t *testing.T) {

}
