package reporter

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/rocky2015aaa/ethdefender/internal/config"
	"github.com/rocky2015aaa/ethdefender/internal/repository/postgres"
	"github.com/rocky2015aaa/ethdefender/internal/services"
	"github.com/rocky2015aaa/ethdefender/internal/services/reporter/handlers"
	httplib "github.com/rocky2015aaa/ethdefender/pkg/http"
	"github.com/rocky2015aaa/ethdefender/pkg/service"
)

const (
	EnvReporterPort = "REPORTER_PORT"
)

type App struct {
	server httplib.Server
}

func NewApp() *App {
	services.SetupWithDB()

	db, err := postgres.New(os.Getenv(services.EnvReportDBUri), config.Default.Database.Log)
	if err != nil {
		log.WithError(err).Fatal("Database init error")
	}

	router := handlers.NewRouter(db)
	server := httplib.NewHTTPServer(router, os.Getenv(EnvReporterPort))

	return &App{server: server}
}

func (a *App) Run(ctx context.Context) {
	service.RunWithGracefulShutdown(ctx, a.server.Run)
}
