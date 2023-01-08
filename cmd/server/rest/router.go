package rest

import (
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/middlewares"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	router.Post("/update/{mType}/{id}/{value}", mh.ReceptionMetricsHandler())
	router.Post("/update/", mh.ReceptionMetricsHandler())
	router.Get("/value/{mType}/{id}", mh.GetMetricHandler())
	router.Post("/value/", mh.GetMetricHandler())
}

func HandleHeathCheck(router chi.Router, hc HealthChecker) {
	router.Get("/ping", hc.CheckDBHandler())
}
