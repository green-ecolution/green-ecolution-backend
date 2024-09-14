package info

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler/v1/errorhandler"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
)

// @Summary		Get info about the app
// @Description	Get info about the app and the server
// @Id				get-app-info
// @Tags			Info
// @Produce		json
// @Success		200	{object}	info.AppInfoResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/info [get]
// @Param			Authorization	header	string	true	"Insert your access token"	default(Bearer <Add access token here>)
func GetAppInfo(svc service.InfoService) fiber.Handler {
	var mapper mapper.InfoHTTPMapper = &generated.InfoHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		domainInfo, err := svc.GetAppInfoResponse(c.Context())
		if err != nil {
			return errorhandler.HandleError(err)
		}

		response := mapper.ToResponse(domainInfo)
		return c.JSON(response)
	}
}
