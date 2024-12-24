package user

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

var (
	userMapper = generated.UserHTTPMapperImpl{}
)

// @Summary	Request to login
// @Descriptio	Request to login to the system. Returns a Login URL
// @Tags		User
// @Produce	json
// @Param		redirect_url	query		string	true	"Redirect URL"
// @Success	200				{object}	entities.LoginResponse
// @Failure	400				{object}	HTTPError
// @Failure	500				{object}	HTTPError
// @Router		/v1/user/login [get]
func Login(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		redirectURL, err := url.ParseRequestURI(c.Query("redirect_url"))
		if err != nil {
			return errorhandler.HandleError(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse redirect url").Error()))
		}

		req := domain.LoginRequest{
			RedirectURL: redirectURL,
		}

		resp, err := svc.LoginRequest(ctx, &req)
		if err != nil {
			return errorhandler.HandleError(service.NewError(service.InternalError, errors.Wrap(err, "failed to login").Error()))
		}

		response := entities.LoginResponse{
			LoginURL: resp.LoginURL.String(),
		}

		return c.JSON(response)
	}
}

// @Summary	Logout from the system
// @Descriptio	Logout from the system
// @Tags		User
// @Param		body	body		entities.LogoutRequest	true	"Logout information"
// @Success	200		{string}	string					"OK"
// @Failure	400		{object}	HTTPError
// @Failure	500		{object}	HTTPError
// @Router		/v1/user/logout [post]
func Logout(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := entities.LogoutRequest{}
		if err := c.BodyParser(&req); err != nil {
			return errorhandler.HandleError(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		domainReq := domain.Logout{
			RefreshToken: req.RefreshToken,
		}

		err := svc.LogoutRequest(ctx, &domainReq)
		if err != nil {
			return errorhandler.HandleError(service.NewError(service.InternalError, errors.Wrap(err, "failed to logout").Error()))
		}

		return c.SendStatus(fiber.StatusOK)
	}
}

// @Summary	Validate login code and request a access token
// @Descriptio	Validate login code and request a access token
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		body			body		entities.LoginTokenRequest	true	"Callback information"
// @Param		redirect_url	query		string						true	"Redirect URL"
// @Success	200				{object}	entities.ClientTokenResponse
// @Failure	400				{object}	HTTPError
// @Failure	500				{object}	HTTPError
// @Router		/v1/user/login/token [post]
func RequestToken(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := entities.LoginTokenRequest{}
		if err := c.BodyParser(&req); err != nil {
			return errorhandler.HandleError(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		redirectURL, err := parseURL(c.Query("redirect_url"))
		if err != nil {
			return errorhandler.HandleError(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse redirect url").Error()))
		}

		domainReq := domain.LoginCallback{
			Code:        req.Code,
			RedirectURL: redirectURL,
		}

		token, err := svc.ClientTokenCallback(ctx, &domainReq)
		if err != nil {
			return errorhandler.HandleError(service.NewError(service.InternalError, errors.Wrap(err, "failed to get token").Error()))
		}

		response := entities.ClientTokenResponse{
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
// @Param			user	body		entities.UserRegisterRequest	true	"User information"
// @Success		201		{object}	entities.UserResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Router			/v1/user [post]
// @Security		Keycloak
func Register(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := entities.UserRegisterRequest{}
		if err := c.BodyParser(&req); err != nil {
			return errorhandler.HandleError(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		domainUser := domain.RegisterUser{
			User: domain.User{
				Email:     req.Email,
				FirstName: req.FirstName,
				LastName:  req.LastName,
				Username:  req.Username,
			},
			Password: req.Password,
			Roles:    req.Roles,
		}

		u, err := svc.Register(ctx, &domainUser)
		if err != nil {
			return errorhandler.HandleError(err)
		}

		response := userMapper.FromResponse(u)

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func parseURL(rawURL string) (*url.URL, error) {
	return url.ParseRequestURI(rawURL)
}

// @Summary		Get all users
// @Description	Get all users
// @Id				get-all-users
// @Tags			User
// @Produce		json
// @Success		200			{object}	entities.UserListResponse
// @Failure		400			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Param			page		query		string	false	"Page"
// @Param			limit		query		string	false	"Limit"
// @Param			limit		query		string	false	"Limit"
// @Param			user_ids	query		string	false	"User IDs"
// @Router			/v1/user [get]
// @Security		Keycloak
func GetAllUsers(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		var domainData []*domain.User
		var err error

		userIDsParam := c.Query("user_ids")
		if userIDsParam == "" {
			domainData, err = svc.GetAll(ctx)
			if err != nil {
				return errorhandler.HandleError(err)
			}
		} else {
			userIDs := strings.Split(userIDsParam, ",")
			domainData, err = svc.GetByIDs(ctx, userIDs)
			if err != nil {
				return errorhandler.HandleError(err)
			}
		}

		data := make([]*entities.UserResponse, len(domainData))
		for i, domain := range domainData {
			data[i] = userMapper.FromResponse(domain)
		}

		return c.Status(fiber.StatusOK).JSON(entities.UserListResponse{
			Data:       data,
			Pagination: entities.Pagination{}, // TODO: Handle pagination
		})
	}
}

// @Summary		Get users by role
// @Description	Get users by role
// @Tags			User
// @Produce		json
// @Success		200		{object}	entities.UserListResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			page	query		string	false	"Page"
// @Param			limit	query		string	false	"Limit"
// @Param			limit	query		string	false	"Role"
// @Router			/v1/user [get]
// @Security		Keycloak
func GetUsersByRole(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()

		role := strings.Clone(c.Params("Role"))
		if role == "" {
			return errorhandler.HandleError(service.NewError(service.BadRequest, "invalid role format"))
		}

		userRole := domain.ParseUserRole(role)
		if userRole == domain.UserRoleUnknown {
			return errorhandler.HandleError(service.NewError(service.BadRequest, "invalid role name"))
		}

		users, err := svc.GetAllByRole(ctx, domain.Role{Name: userRole})
		if err != nil {
			return errorhandler.HandleError(service.NewError(service.InternalError, errors.Wrap(err, "failed to get users by role").Error()))
		}
		response := make([]entities.UserResponse, len(users))

		for i, user := range users {
			userResponse := userMapper.FromResponse(user)
			response[i] = *userResponse
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// @Summary		Get a user by ID
// @Description	Get a user by ID
// @Tags			User
// @Produce		json
// @Success		200		{object}	entities.UserResponse
// @Failure		400		{object}	HTTPError
// @Failure		404		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Param			user_id	path		string	true	"User ID"
// @Router			/v1/user/{user_id} [get]
// @Security		Keycloak
func GetUserByID(_ service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

		data := make([]*entities.UserResponse, len(domainData))
		for i, domain := range domainData {
			data[i] = userMapper.FromResponse(domain)
		}

		return c.Status(fiber.StatusOK).JSON(entities.UserListResponse{
			Data:       data,
			Pagination: entities.Pagination{}, // TODO: Handle pagination
		})
	}
}

var group singleflight.Group

// @Summary		Refresh token
// @Description	Refresh token
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			body	body		entities.RefreshTokenRequest	true	"Refresh token information"
// @Success		200		{object}	entities.ClientTokenResponse
// @Failure		400		{object}	HTTPError
// @Failure		500		{object}	HTTPError
// @Router			/v1/user/token/refresh [post]
func RefreshToken(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.Context()
		req := entities.RefreshTokenRequest{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		if req.RefreshToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(service.NewError(service.BadRequest, errors.New("refresh token is required").Error()))
		}

		jwtClaims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(req.RefreshToken, jwtClaims, func(token *jwt.Token) (any, error) {
			return token, nil
		})

		var sub string
		if claims, ok := jwtClaims["sub"]; ok {
			if e, ok := claims.(string); ok {
				sub = e
			}
		}

		if sub == "" {
			return c.Status(fiber.StatusBadRequest).JSON(service.NewError(service.BadRequest, errors.Wrap(err, "failed to parse request").Error()))
		}

		data, err, _ := group.Do(sub, func() (any, error) {
			return svc.RefreshToken(ctx, req.RefreshToken)
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(service.NewError(service.InternalError, errors.Wrap(err, "failed to refresh token").Error()))
		}

		token := data.(*domain.ClientToken)
		response := entities.ClientTokenResponse{
			AccessToken:  token.AccessToken,
			ExpiresIn:    token.ExpiresIn,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		}

		return c.JSON(response)
	}
}
