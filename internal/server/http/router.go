package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/fileimport"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/region"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/sensor"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/tree"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/treecluster"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/user"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/middleware"
)

func (s *Server) privateRoutes(app *fiber.App) {
	grp := app.Group("/api/v1")

	grp.Mount("/info", info.RegisterRoutes(s.services.InfoService))
	grp.Mount("/cluster", treecluster.RegisterRoutes(s.services.TreeClusterService))
	grp.Mount("/tree", tree.RegisterRoutes(s.services.TreeService))
	grp.Mount("/sensor", sensor.RegisterRoutes(s.services.MqttService))
	grp.Mount("/user", user.RegisterRoutes(s.services.AuthService))
	grp.Mount("/role", user.RegisterRoutes(s.services.AuthService))
	grp.Mount("/region", region.RegisterRoutes(s.services.RegionService))
	grp.Mount("/import", fileimport.RegisterRoutes(s.services.TreeService))
}

func (s *Server) publicRoutes(app *fiber.App) {
	app.Use("/", middleware.HealthCheck(s.services))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	grp := app.Group("/api/v1")
	grp.Get("/swagger/*", swagger.HandlerDefault)
	grp.Post("/user/logout", user.Logout(s.services.AuthService))
	grp.Get("/user/login", user.Login(s.services.AuthService))
	grp.Post("/user/login/token", user.RequestToken(s.services.AuthService))
}
