package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(basePath string) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount(basePath, infoRouter())
	r.Mount(basePath, healthRouter())
	r.Mount(basePath, apiRouter())

	return r
}

func apiRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Mount("/", launchHandler{}.routes())
	r.Mount("/", itemHandler{}.routes())
	r.Mount("/", logHandler{}.routes())

	r.Get("/v1/{projectName}/settings", notImplemented)

	return r
}
