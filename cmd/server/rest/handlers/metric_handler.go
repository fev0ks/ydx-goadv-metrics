package handlers

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type MetricsHandler struct {
	Ctx        context.Context
	Repository server.MetricRepository
}

func (mh *MetricsHandler) ReceptionMetricsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		_ = context.WithValue(mh.Ctx, "execCtxId", "rm")
		name := mux.Vars(request)["name"]
		mType := mux.Vars(request)["mType"]
		value := mux.Vars(request)["value"]
		log.Printf("request vars - name: '%s', type: '%s', value: '%s'", name, mType, value)
		metric, err := model.NewMetric(name, model.MTypeValueOf(mType), value)
		if err != nil {
			log.Printf("failed to parse metric request: %v\n", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}
		err = mh.Repository.SaveMetric(metric)
		if err != nil {
			log.Printf("failed to save metric: %v\n", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		writer.WriteHeader(http.StatusOK)
		return
	}
}
