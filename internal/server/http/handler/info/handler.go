package info

import (
	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/info"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/info/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/handler"
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
// @Router			/info [get]
func GetAppInfo(svc service.InfoService) fiber.Handler {
	var mapper info.InfoHTTPMapper = &generated.InfoHTTPMapperImpl{}

	return func(c *fiber.Ctx) error {
		info, err := svc.GetAppInfoResponse(c.Context())
		if err != nil {
			return handler.HandleError(err)
		}

		data := mapper.ToResponse(info)
		return c.JSON(data)
	}
}
