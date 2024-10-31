package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

func RegisterRoutes(svc service.AuthService) *fiber.App {
	app := fiber.New()

	app.Post("/", Register(svc))
	app.Get("/", GetAllUsers(svc))
	app.Get("/:id", GetUserByID(svc))
	app.Put("/:id", UpdateUserByID(svc))
	app.Delete("/:id", DeleteUserByID(svc))
	app.Get("/:id/roles", GetUserRoles(svc))

	return app
}

func RegisterPublicRoutes(svc service.AuthService) *fiber.App {
	app := fiber.New()

	app.Post("/logout", Logout(svc))
	app.Get("/login", Login(svc))
	app.Post("/login/token", RequestToken(svc))
	app.Post("/token/refresh", RefreshToken(svc))

	return app
}
