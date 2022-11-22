package rest

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(_ context.Context) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}

func HandleMetricRequests(router *mux.Router, mh handlers.MetricsHandler) {
	router.Methods(http.MethodPost).Path("/update/{mType}/{name}/{value}").HandlerFunc(mh.ReceptionMetricsHandler())
}
