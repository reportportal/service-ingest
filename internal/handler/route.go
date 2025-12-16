package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// validate is a singleton validator instance used for struct validation across handlers.

type Handlers struct {
	Launch *LaunchHandler
	Item   *ItemHandler
	Log    *LogHandler
}

func NewRouter(basePath string, handlers Handlers) chi.Router {
	initValidatorOnce()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route(basePath, func(r chi.Router) {
		r.Mount("/info", infoRouter())
		r.Mount("/health", healthRouter())
		r.Mount("/", apiRouter(handlers))
	})

	return r
}

func apiRouter(h Handlers) chi.Router {
	r := chi.NewRouter()

	r.Route("/v1/{projectName}", func(r chi.Router) {
		r.Mount("/launch", h.Launch.routesV1())
		r.Mount("/item", h.Item.routesV1())
		r.Get("/settings", RespondNotImplemented)
	})

	r.Route("/v2/{projectName}", func(r chi.Router) {
		r.Mount("/launch", h.Launch.routesV2())
		r.Mount("/item", h.Item.routesV2())
		r.Mount("/log", h.Log.routes())
	})

	return r
}
