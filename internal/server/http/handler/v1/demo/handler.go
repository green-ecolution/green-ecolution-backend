package demo

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	demoplugin "github.com/green-ecolution/green-ecolution-backend/plugin/demo_plugin"
)


func getPluginFiles() *fiber.App { 
  app := fiber.New()

	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(demoplugin.F),
    PathPrefix: "dist",
    Browse: true,
	}))

  return app
}
