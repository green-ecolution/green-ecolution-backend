//go:generate mockery
//go:generate swag fmt
//go:generate swag init --requiredByDefault
//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/...
//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt/entities/...
//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen github.com/green-ecolution/green-ecolution-backend/internal/storage/mongodb/entities/...
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/docs"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/local"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/mongodb"
)

var version = "develop"

func setSwaggerInfo(appURL *url.URL) {
	docs.SwaggerInfo.Title = "Green Ecolution Management API"
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Description = "This is the API for the Green Ecolution Management System."
	docs.SwaggerInfo.Host = appURL.Host
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Green Ecolution
//	@contact.url	https://green-ecolution.de
//	@contact.email	info@green-ecolution.de

// @license.name	GPL-3.0
// @license.url	https://raw.githubusercontent.com/green-ecolution/green-ecolution-management/develop/LICENSE
func main() {
	cfg, err := config.GetAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version: %s\n", version)
	if cfg.Development {
		fmt.Println("Running in dev mode")
		cfg.LogLevel = "debug"
	}

	setSwaggerInfo(cfg.URL)

	logg := logger.CreateLogger(os.Stdout, cfg.LogFormat, cfg.LogLevel)
	slog.SetDefault(logg)

	slog.Info("Starting Green Space Management API")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	localRepo, err := local.NewRepository(cfg)
	if err != nil {
		slog.Error("Error while creating local repository", "error", err)
		return
	}

	dbRepo, err := mongodb.NewRepository(cfg)
	if err != nil {
		slog.Error("Error while creating MongoDB repository", "error", err)
		return
	}

	repositories := &storage.Repository{
		Info:   localRepo.Info,
		Sensor: dbRepo.Sensor,
		Tree:   dbRepo.Tree,
	}

	services := domain.NewService(cfg, repositories)
	httpServer := http.NewServer(cfg, services)
	mqttServer := mqtt.NewMqtt(cfg, services)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		mqttServer.RunSubscriber(ctx)
	}()

	go func() {
		defer wg.Done()
		if err := httpServer.Run(ctx); err != nil {
			slog.Error("Error while running HTTP Server", "error", err)
		}
	}()

	wg.Wait()
}
