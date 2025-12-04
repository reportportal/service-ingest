package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

// validate is a singleton validator instance used for struct validation across handlers.
var validate = validator.New()

func NewRouter(basePath string) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount(basePath+"/info", infoRouter())
	r.Mount(basePath+"/health", healthRouter())
	r.Mount(basePath, apiRouter())

	return r
}

func apiRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)

	r.Mount("/v1/{projectName}/launch", launchHandler{}.routesV1())
	r.Mount("/v2/{projectName}/launch", launchHandler{}.routesV2())
	r.Mount("/v1/{projectName}/item", itemHandler{}.routesV1())
	r.Mount("/v2/{projectName}/item", itemHandler{}.routesV2())
	r.Mount("/v2/{projectName}/log", logHandler{}.routes())

	r.Get("/v1/{projectName}/settings", respondNotImplemented)

	return r
}
