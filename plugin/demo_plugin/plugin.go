package demoplugin

import (
	"context"
	"embed"
	"log/slog"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

// Embed a single file
//
//go:embed dist/*
var f embed.FS

var (
  username = "demo_plugin"
  password = "demo_plugin"
)

func DemoPluginStart(ctx context.Context) {
	app := fiber.New()
	app.Mount("/", getPluginFiles())

	go func() {
		<-ctx.Done()
		app.Shutdown()
	}()

	go func() {
		app.Listen(":3020")
	}()

	pluginPath, err := url.Parse("http://localhost:3020/")
	if err != nil {
		panic(err)
	}

	hostPath, err := url.Parse("http://localhost:3000/")
	if err != nil {
		panic(err)
	}

	worker, err := NewPluginWorker(
		WithPluginName("demo_plugin"),
		WithPluginPath(pluginPath),
		WithHost(hostPath),
	)
  if err != nil {
    panic(err)
  }
  
  if err := worker.Register(ctx, username, password); err != nil {
    panic(err)
  }

  if err := worker.Run(ctx); err != nil {
    slog.Error("Failed to send heartbeat", "error", err)
  }
}
