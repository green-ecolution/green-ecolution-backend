package plugin

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func getPluginFiles(c *fiber.Ctx) error {
	url, err := url.Parse("http://localhost:3000")
	if err != nil {
		return err
	}

	plugin := c.Params("plugin")

	reverseProxy := httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(url)
			r.Out.Host = r.In.Host
			r.Out.URL.Path = strings.Replace(r.In.URL.Path, "/api/v1/plugin/"+plugin, "/api/v1/demo_plugin", 1)
			r.SetXForwarded()
		},
	}

	return adaptor.HTTPHandler(&reverseProxy)(c)
}
