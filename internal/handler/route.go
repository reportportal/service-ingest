package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// validate is a singleton validator instance used for struct validation across handlers.

func NewRouter(basePath string, launchHandler LaunchHandler, itemHandler ItemHandler, logHandler LogHandler) chi.Router {
	initValidatorOnce()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route(basePath, func(r chi.Router) {
		r.Mount("/info", infoRouter())
		r.Mount("/health", healthRouter())
		r.Mount("/", apiRouter(launchHandler, itemHandler, logHandler))
	})

	return r
}

func apiRouter(launchHandler LaunchHandler, itemHandler ItemHandler, logHandler LogHandler) chi.Router {
	r := chi.NewRouter()

	r.Route("v1/{projectName}", func(r chi.Router) {
		r.Mount("/launch", launchHandler.routesV1())
		r.Mount("/item", itemHandler.routesV1())
		r.Get("/settings", RespondNotImplemented)
	})

	r.Route("v2/{projectName}", func(r chi.Router) {
		r.Mount("/launch", launchHandler.routesV2())
		r.Mount("/item", itemHandler.routesV2())
		r.Mount("/log", logHandler.routes())
	})

	return r
}
