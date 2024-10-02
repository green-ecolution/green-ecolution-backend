package demoplugin

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func getPluginFiles() *fiber.App {
	app := fiber.New()

	app.Use(filesystem.New(filesystem.Config{
		Root:       http.FS(f),
		PathPrefix: "dist",
		Browse:     true,
	}))

	return app
}
