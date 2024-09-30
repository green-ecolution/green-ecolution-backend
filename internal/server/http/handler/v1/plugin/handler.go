package plugin

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	csvimport "github.com/green-ecolution/green-ecolution-backend/plugin/csv_import"
)


func getPluginFiles() *fiber.App { 
  app := fiber.New()
  dir, err := csvimport.F.ReadDir("dist") 
  if err != nil {
    fmt.Printf("Error: %v\n", err)
  }
  fmt.Printf("csvimport.F: %v\n", dir)
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(csvimport.F),
    PathPrefix: "dist",
    Browse: true,
	}))

  return app
}
