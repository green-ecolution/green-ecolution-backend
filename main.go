//go:generate mockery
//go:generate swag fmt
//go:generate swag init --requiredByDefault
//go:generate goverter gen github.com/green-ecolution/green-ecolution-backend/internal/server/http/entities/...
//go:generate goverter gen github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt/entities/...
//go:generate goverter gen github.com/green-ecolution/green-ecolution-backend/internal/storage/mongodb/entities/...
//go:generate goverter gen github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/...
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/green-ecolution/green-ecolution-backend/config"
	"github.com/green-ecolution/green-ecolution-backend/docs"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/auth/keycloak"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/local"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

var version = "develop"

//	@title			Green Space Management API
//	@version		develop
//	@description	This is the API for the Green Ecolution Management System.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Green Ecolution
//	@contact.url	https://green-ecolution.de
//	@contact.email	info@green-ecolution.de

// @license.name	GPL-3.0
// @license.url	https://raw.githubusercontent.com/green-ecolution/green-ecolution-management/develop/LICENSE
func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Version: %s\n", version)

	fmt.Println("Server Port: ", viper.GetInt("server.port"))

	if cfg.Server.Development {
		fmt.Println("Running in dev mode")
		cfg.Server.Logs.Level = "debug"
	}

	logg := logger.CreateLogger(os.Stdout, cfg.Server.Logs.Format, cfg.Server.Logs.Level)
	slog.SetDefault(logg)

	setSwaggerInfo(cfg.Server.AppURL)

	slog.Info("Starting Green Space Management API")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := pgx.Connect(ctx, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Server.Database.Host, cfg.Server.Database.Port, cfg.Server.Database.Username, cfg.Server.Database.Password, cfg.Server.Database.Name))
	if err != nil {
		slog.Error("Error while connecting to PostgreSQL", "error", err)
		return
	}
	defer conn.Close(context.Background())
	postgresRepo := postgres.NewRepository(conn)

	localRepo, err := local.NewRepository(cfg)
	if err != nil {
		slog.Error("Error while creating local repository", "error", err)
		return
	}

	//	dbRepo, err := mongodb.NewRepository(cfg)
	if err != nil {
		slog.Error("Error while creating MongoDB repository", "error", err)
		return
	}

	keycloakRepo := keycloak.NewKeycloakRepository(&cfg.IdentityAuth)

	repositories := &storage.Repository{
		Auth:        keycloakRepo,
		Info:        localRepo.Info,
		Sensor:      postgresRepo.Sensor,
		Tree:        postgresRepo.Tree,
		TreeCluster: postgresRepo.TreeCluster,
		Vehicle:     postgresRepo.Vehicle,
		Flowerbed:   postgresRepo.Flowerbed,
		Image:       postgresRepo.Image,
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

func setSwaggerInfo(appURL string) {
	slog.Info("Setting Swagger info")

	var schemes []string
	var trimmedAppURL string
	if strings.HasPrefix(appURL, "http://") {
		schemes = []string{"http"}
		trimmedAppURL = strings.TrimPrefix(appURL, "http://")
	} else {
		trimmedAppURL = strings.TrimPrefix(appURL, "https://")
		schemes = []string{"https"}
	}

	docs.SwaggerInfo.Title = "Green Ecolution Management API"
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Description = "This is the API for the Green Ecolution Management System."
	docs.SwaggerInfo.Host = trimmedAppURL
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = schemes
}
