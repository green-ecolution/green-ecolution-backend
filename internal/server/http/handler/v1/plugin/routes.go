package plugin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(r fiber.Router, svc service.PluginService, middlewares ...fiber.Handler) {
	handlers := append([]fiber.Handler{}, middlewares...)
	handlers = append(handlers, GetPluginsList(svc))
	r.Get("/", handlers...)

	r.Post("/register", registerPlugin(svc))

	handlers = append([]fiber.Handler{}, middlewares...)
	handlers = append(handlers, GetPluginInfo(svc))
	r.Get("/:plugin", handlers...)

	handlers = append([]fiber.Handler{}, middlewares...)
	handlers = append(handlers, pluginHeartbeat(svc))
	r.Post("/:plugin/heartbeat", handlers...)

	handlers = append([]fiber.Handler{}, middlewares...)
	handlers = append(handlers, unregisterPlugin(svc))
	r.Post("/:plugin/unregister", handlers...)

	r.Use("/:plugin", getPluginFiles(svc))
}
