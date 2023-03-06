package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/middlewares"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
)

func NewRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middlewares.TimerTrace)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(3, rest.ApplicationJSON, rest.TextPlain, rest.TextHTML))
	router.Use(middlewares.Decompress)
	return router
}

func HandleMetricRequests(router chi.Router, mh *MetricsHandler) {
	router.Get("/", mh.GetMetricsHandler())
	router.Post("/update/{mType}/{id}/{value}", mh.ReceptionMetricHandler())
	router.Post("/update/", mh.ReceptionMetricHandler())
	router.Post("/updates/", mh.ReceptionMetricsHandler())
	router.Get("/value/{mType}/{id}", mh.GetMetricHandler())
	router.Post("/value/", mh.GetMetricHandler())
}

func HandleHeathCheck(router chi.Router, hc HealthChecker) {
	router.Get("/ping", hc.CheckDBHandler())
}

func HandlePprof(router chi.Router) {
	router.Mount("/debug", middleware.Profiler())
}
