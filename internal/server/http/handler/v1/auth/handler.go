package auth

import (
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
)

// @Summary	Requst to login
// @Descriptio	Requst to login to the system. Returns a Login URL
// @Tags		Login
// @Produce	json
// @Param		redirect_url	query		string	true	"Redirect URL"
// @Success	200				{object}	auth.LoginResponse
// @Failure	400				{object}	HTTPError
// @Failure	500				{object}	HTTPError
// @Router		/v1/login [get]
func Login(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		redirectURLRaw := c.Query("redirect_url")
		redirectURL, err := url.ParseRequestURI(redirectURLRaw)
		log.Println(redirectURL)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse redirect url").Error()))
		}

		req := domain.LoginRequest{
			RedirectURL: redirectURL,
		}

		resp, err := svc.LoginRequest(ctx, &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(service.NewError(service.InternalError, errors.Wrap(err, "failed to handle login request").Error()))
		}

		response := auth.LoginResponse{
			LoginURL: resp.LoginURL.String(),
		}

		return c.JSON(response)
	}
}

// @Summary	Requst to login
// @Descriptio	Requst a access token
// @Tags		Login
// @Accept		json
// @Produce	json
// @Param		body			body		auth.LoginTokenRequest	true	"Callback information"
// @Param		redirect_url	query		string					true	"Redirect URL"
// @Success	200				{object}	auth.ClientTokenResponse
// @Failure	400				{object}	HTTPError
// @Failure	500				{object}	HTTPError
// @Router		/v1/token [post]
func RequestToken(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := auth.LoginTokenRequest{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		redirectURL, err := parseURL(c.Query("redirect_url"))
		if err != nil {
			return err
		}

		domainReq := domain.LoginCallback{
			Code:        req.Code,
			RedirectURL: redirectURL,
		}

		token, err := svc.ClientTokenCallback(ctx, &domainReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(service.NewError(service.InternalError, errors.Wrap(err, "failed to handle token request").Error()))
		}

		response := auth.ClientTokenResponse{
			AccessToken:  token.AccessToken,
			ExpiresIn:    token.ExpiresIn,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		}

		return c.JSON(response)
	}
}

// @Summary		Register a new user
// @Description	Register a new user
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			user	body		auth.RegisterUserRequest	true	"User information"
// @Success		201		{object}	auth.UserResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
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

func parseURL(rawURL string) (*url.URL, error) {
	return url.ParseRequestURI(rawURL)
}
