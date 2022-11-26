package rest

import (
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/handlers"
	"github.com/go-chi/chi/v5"
)

func HandleMetricRequests(router *chi.Mux, mh handlers.MetricsHandler) {
	router.Get("/", mh.GetMetricsHandler())
	router.Post("/update/{mType}/{name}/{value}", mh.ReceptionMetricsHandler())
	router.Get("/value/{mType}/{name}", mh.GetMetricHandler())
}
