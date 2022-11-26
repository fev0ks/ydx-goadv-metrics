package rest

import (
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func HandleMetricRequests(router *mux.Router, mh handlers.MetricsHandler) {
	router.Methods(http.MethodGet).Path("/").HandlerFunc(mh.GetMetricsHandler())
	router.Methods(http.MethodPost).Path("/update/{mType}/{name}/{value}").HandlerFunc(mh.ReceptionMetricsHandler())
	router.Methods(http.MethodGet).Path("/value/{mType}/{name}").HandlerFunc(mh.GetMetricHandler())
}
