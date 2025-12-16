package app

import (
	"net/http"

	"github.com/reportportal/service-ingest/internal/config"
	"github.com/reportportal/service-ingest/internal/handler"
	"github.com/reportportal/service-ingest/internal/service"
)

type App struct {
	server *http.Server
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func New(cfg *config.Config) (*App, error) {
	var launchRepo service.LaunchRepository
	var itemRepo service.ItemRepository
	var logRepo service.LogRepository

	launchService := service.NewLaunchService(launchRepo)
	itemService := service.NewItemService(itemRepo)
	logService := service.NewLogService(logRepo)

	handlers := handler.Handlers{
		Launch: handler.NewLaunchHandler(launchService),
		Item:   handler.NewItemHandler(itemService),
		Log:    handler.NewLogHandler(logService),
	}

	router := handler.NewRouter(cfg.Server.BasePath, handlers)

	server := &http.Server{
		Addr:    cfg.Server.Addr(),
		Handler: router,
	}

	return &App{server: server}, nil
}
