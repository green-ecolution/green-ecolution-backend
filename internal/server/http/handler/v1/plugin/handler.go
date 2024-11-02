package plugin

import (
	"fmt"
	"log/slog"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	domain "github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

func getPluginFiles(c *fiber.Ctx) error {
	// pluginMutex.RLock()
	plugin := registeredPlugins[c.Params("plugin")]
	// pluginMutex.RUnlock()

	reverseProxy := httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(&plugin.Path)
			r.Out.Host = r.In.Host
			r.Out.URL.Path = strings.Replace(r.In.URL.Path, "/api/v1/plugin/"+plugin.Slug, plugin.Path.String(), 1)
			r.SetXForwarded()
		},
	}

	return adaptor.HTTPHandler(&reverseProxy)(c)
}

//	@Summary		Register a plugin
//	@Description	Register a plugin
//	@Id				register-plugin
//	@Tags			Plugin
//	@Produce		json
//	@Success		200	{object}	entities.ClientTokenResponse
//	@Failure		400	{object}	HTTPError
//	@Failure		401	{object}	HTTPError
//	@Failure		403	{object}	HTTPError
//	@Failure		404	{object}	HTTPError
//	@Failure		500	{object}	HTTPError
//	@Router			/v1/plugin [post]
//
// @Param			body	body	entities.PluginRegisterRequest	true	"Plugin registration request"
func registerPlugin(svc service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.PluginRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if req.Auth.Username == "" || req.Auth.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Username or password is empty")
		}

		// Authenticate the plugin
		token, err := svc.AuthPlugin(c.Context(), &domain.AuthPlugin{
			Username: req.Auth.Username,
			Password: req.Auth.Password,
		})
		if err != nil {
			return err
		}

		path, err := url.Parse(req.Path)
		if err != nil {
			slog.Error("Failed to parse plugin path", "error", err)
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse plugin path")
		}

		// pluginMutex.Lock()
		// defer pluginMutex.Unlock()

		if _, ok := registeredPlugins[req.Slug]; ok {
			return c.Status(fiber.StatusForbidden).SendString("plugin already registered")
		}

		registeredPlugins[req.Slug] = &domain.Plugin{
			Name:          req.Name,
			Path:          *path,
			LastHeartbeat: time.Now(),
			Version:       req.Version,
			Description:   req.Description,
			Slug:          req.Slug,
		}

		slog.Info("Plugin registered", "plugin", req.Name, "version", req.Version, "slug", req.Slug)
		slog.Debug("Plugin registered", "plugin", fmt.Sprintf("%+v", registeredPlugins[req.Slug]))

		response := entities.ClientTokenResponse{
			AccessToken:  token.AccessToken,
			ExpiresIn:    token.ExpiresIn,
			RefreshToken: token.RefreshToken,
			TokenType:    token.TokenType,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}

// @Summary		Heartbeat for a plugin
// @Description	Heartbeat for a plugin
// @Id				plugin-heartbeat
// @Tags			Plugin
// @Produce		json
// @Success		200	{object}	string
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/plugin/{plugin_slug}/heartbeat [post]
// @Param			plugin_slug	path	string	true	"Name of the plugin specified by slug during registration"
// @Security		Keycloak
func pluginHeartbeat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// pluginMutex.Lock()
		// defer pluginMutex.Unlock()

		plugin := registeredPlugins[c.Params("plugin")]
		plugin.LastHeartbeat = time.Now()
		registeredPlugins[c.Params("plugin")] = plugin

		return c.Status(fiber.StatusOK).SendString("Heartbeat received")
	}
}

// @Summary		Get a list of all registered plugins
// @Description	Get a list of all registered plugins
// @Id				get-plugins-list
// @Tags			Plugin
// @Produce		json
// @Success		200	{object}	entities.PluginListResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/plugin [get]
func GetPluginsList() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// pluginMutex.RLock()
		// defer pluginMutex.RUnlock()

		plugins := utils.MapKeysSlice(registeredPlugins, func(_ string, v *domain.Plugin) entities.PluginResponse {
			return entities.PluginResponse{
				Slug:        v.Slug,
				Name:        v.Name,
				Version:     v.Version,
				Description: v.Description,
				HostPath:    v.Path.String(),
			}
		})

		return c.Status(fiber.StatusOK).JSON(entities.PluginListResponse{
			Plugins: plugins,
		})
	}
}

// @Summary		Get a plugin info
// @Description	Get a plugin info
// @Id				get-plugin-info
// @Tags			Plugin
// @Produce		json
// @Success		200	{object}	entities.PluginResponse
// @Failure		400	{object}	HTTPError
// @Failure		401	{object}	HTTPError
// @Failure		403	{object}	HTTPError
// @Failure		404	{object}	HTTPError
// @Failure		500	{object}	HTTPError
// @Router			/v1/plugin/{plugin_slug} [get]
// @Param			plugin_slug	path	string	true	"Slug of the plugin"
func GetPluginInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// pluginMutex.RLock()
		// defer pluginMutex.RUnlock()

		plugin, ok := registeredPlugins[c.Params("plugin")]
		if !ok {
			return c.Status(fiber.StatusNotFound).SendString("plugin not found")
		}

		return c.Status(fiber.StatusOK).JSON(entities.PluginResponse{
			Slug:        plugin.Slug,
			Name:        plugin.Name,
			Version:     plugin.Version,
			Description: plugin.Description,
			HostPath:    plugin.Path.String(),
		})
	}
}
