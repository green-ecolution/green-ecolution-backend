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

	"github.com/green-ecolution/green-ecolution-backend/docs"
	"github.com/green-ecolution/green-ecolution-backend/internal/config"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/logger"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/http"
	"github.com/green-ecolution/green-ecolution-backend/internal/server/mqtt"
	"github.com/green-ecolution/green-ecolution-backend/internal/service"
	"github.com/green-ecolution/green-ecolution-backend/internal/service/domain"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/auth"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/local"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres"
	_ "github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/openrouteservice"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/routing/valhalla"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/s3"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker"
	"github.com/green-ecolution/green-ecolution-backend/internal/worker/subscriber"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"github.com/twpayne/go-geos"
	pgxgeos "github.com/twpayne/pgx-geos"
)

var version = "develop"

//	@title			Green Space Management API
//	@version		develop
//	@description	This is the API for the Green Ecolution Management System.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Green Ecolution
//	@contact.url	https://green-ecolution.de
//	@contact.email	info@green-ecolution.de

//	@license.name	AGPL
//	@license.url	https://raw.githubusercontent.com/green-ecolution/green-ecolution-management/develop/LICENSE

// @securitydefinitions.oauth2.accessCode	Keycloak
// @tokenUrl								https://auth.green-ecolution.de/realms/green-ecolution-dev/protocol/openid-connect/token
// @authorizationUrl						https://auth.green-ecolution.de/realms/green-ecolution-dev/protocol/openid-connect/auth
// @in										header
// @name									Authorization
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

	startAppServices(ctx, cfg)
}

func postgresRepo(ctx context.Context, cfg *config.Config) (repo *storage.Repository, closeFn func()) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Server.Database.Host, cfg.Server.Database.Port, cfg.Server.Database.Username, cfg.Server.Database.Password, cfg.Server.Database.Name)

	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		slog.Error("Error while parsing PostgreSQL connection string", "error", err)
		return nil, nil
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxgeos.Register(ctx, conn, geos.NewContext())
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		slog.Error("Error while connecting to PostgreSQL", "error", err)
		return nil, nil
	}

	return postgres.NewRepository(pool), pool.Close
}

func startAppServices(ctx context.Context, cfg *config.Config) {
	repositories, closeFn, err := initializeRepositories(ctx, cfg)
	if err != nil {
		slog.Error("Failed to initialize repositories", "error", err)
		return
	}
	defer closeFn()

	em := initializeEventManager()

	services := domain.NewService(cfg, repositories, em)
	httpServer := http.NewServer(cfg, services)
	mqttServer := mqtt.NewMqtt(cfg, services)

	runServices(ctx, httpServer, mqttServer, em, services)
}

func initializeRepositories(ctx context.Context, cfg *config.Config) (*storage.Repository, func(), error) {
	localRepo, err := local.NewRepository(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating local repository: %w", err)
	}

	// can be switched between ors and valhalla
	// routingRepo, err := openrouteservice.NewRepository(cfg)
	routingRepo, err := valhalla.NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	keycloakRepo := auth.NewRepository(&cfg.IdentityAuth)

	s3Repos, err := s3.NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	postgresRepo, closeFn := postgresRepo(ctx, cfg)

	repositories := &storage.Repository{
		Auth: keycloakRepo.Auth,
		User: keycloakRepo.User,

		Info:         localRepo.Info,
		Sensor:       postgresRepo.Sensor,
		Tree:         postgresRepo.Tree,
		TreeCluster:  postgresRepo.TreeCluster,
		Vehicle:      postgresRepo.Vehicle,
		Flowerbed:    postgresRepo.Flowerbed,
		Image:        postgresRepo.Image,
		Region:       postgresRepo.Region,
		WateringPlan: postgresRepo.WateringPlan,
		Routing:      routingRepo.Routing,
		GpxBucket:    s3Repos.GpxBucket,
	}

	return repositories, closeFn, nil
}

func initializeEventManager() *worker.EventManager {
	return worker.NewEventManager(
		entities.EventTypeUpdateTree,
		entities.EventTypeUpdateTreeCluster,
		entities.EventTypeCreateTree,
		entities.EventTypeDeleteTree,
		entities.EventTypeNewSensorData,
		entities.EventTypeUpdateWateringPlan,
	)
}

func runServices(ctx context.Context, httpServer *http.Server, mqttServer *mqtt.Mqtt, em *worker.EventManager, services *service.Services) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		mqttServer.RunSubscriber(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		em.Run(ctx)
	}()

	runEventSubscriptions(ctx, &wg, em, services)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Run(ctx); err != nil {
			slog.Error("Error while running HTTP Server", "error", err)
		}
	}()

	wg.Wait()
}

func runEventSubscriptions(ctx context.Context, wg *sync.WaitGroup, em *worker.EventManager, services *service.Services) {
	subscribers := []worker.Subscriber{
		subscriber.NewUpdateTreeSubscriber(services.TreeClusterService),
		subscriber.NewCreateTreeSubscriber(services.TreeClusterService),
		subscriber.NewDeleteTreeSubscriber(services.TreeClusterService),
		subscriber.NewSensorDataSubscriber(services.TreeClusterService, services.TreeService),
		subscriber.NewUpdateWateringPlanSubscriber(services.TreeClusterService),
	}

	for _, sub := range subscribers {
		wg.Add(1)
		go func(sub worker.Subscriber) {
			defer wg.Done()
			if err := em.RunSubscription(ctx, sub); err != nil {
				slog.Error("stop subscription with err", "eventType", sub.EventType(), "err", err)
			}
		}(sub)
	}
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
