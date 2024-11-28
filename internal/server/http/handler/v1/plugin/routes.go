package plugin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(r fiber.Router, svc service.PluginService) {
	r.Post("/register", registerPlugin(svc))

	r.Post("/:plugin/heartbeat", pluginHeartbeat(svc))
	r.Use("/:plugin", getPluginFiles(svc))
}

func RegisterPrivateRoutes(r fiber.Router, svc service.PluginService) {
	r.Get("/", GetPluginsList(svc))
	r.Get("/:plugin", GetPluginInfo(svc))
}
