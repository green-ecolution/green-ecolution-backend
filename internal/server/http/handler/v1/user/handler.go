package user

import (
	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

// @Summary		Register a new user
// @Description	Register a new user
// @Tags			User
// @Accept			json
// @Produce		json
// @Success		201	{object}	auth.UserResponse
// @Failure		400	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/user [post]
func Register(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := auth.RegisterUserRequest{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errors.Wrap(err, "failed to parse request").Error(),
			})
		}

		domainUser := domain.RegisterUser{
			User: domain.User{
				Email:     req.User.Email,
				FirstName: req.User.FirstName,
				LastName:  req.User.LastName,
				Username:  req.User.Username,
			},
			Password: req.Password,
			Roles:    req.Roles,
		}

		user, err := svc.Register(ctx, &domainUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": errors.Wrap(err, "failed to register user").Error(),
			})
		}

		response := auth.UserResponse{
			ID:            user.ID,
			CreatedAt:     user.CreatedAt,
			Email:         user.Email,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Username:      user.Username,
			EmployeeID:    user.EmployeeID,
			PhoneNumber:   user.PhoneNumber,
			EmailVerified: user.EmailVerified,
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}
