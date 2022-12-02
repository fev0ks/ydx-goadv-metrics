package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/pages"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"

	"github.com/go-chi/chi/v5"
)

type MetricsHandler struct {
	Ctx        context.Context
	Repository server.MetricRepository
}

func NewMetricsHandler(ctx context.Context, repository server.MetricRepository) *MetricsHandler {
	return &MetricsHandler{Ctx: ctx, Repository: repository}
}

func (mh *MetricsHandler) ReceptionMetricsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		name := chi.URLParam(request, "name")
		mType := chi.URLParam(request, "mType")
		value := chi.URLParam(request, "value")
		log.Printf("request vars - name: '%s', type: '%s', value: '%s'", name, mType, value)
		if name == "" {
			http.Error(writer, "metric 'name' must be specified", http.StatusBadRequest)
			return
		}
		if mType == "" {
			http.Error(writer, "metric 'mType' must be specified", http.StatusBadRequest)
			return
		}
		if value == "" {
			http.Error(writer, "metric 'value' must be specified", http.StatusBadRequest)
			return
		}
		metric, err := model.NewMetric(name, model.MTypeValueOf(mType), value)
		if err != nil {
			log.Printf("failed to parse metric request: %v\n", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if metric.MType == model.NanType {
			err = fmt.Errorf("type '%s' is not supported", mType)
			log.Printf("failed to save metric: %v\n", err)
			http.Error(writer, err.Error(), http.StatusNotImplemented)
			return
		}
		err = mh.Repository.SaveMetric(metric)
		if err != nil {
			log.Printf("failed to save metric: %v\n", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (mh *MetricsHandler) GetMetricsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Get metrics")
		metrics := mh.Repository.GetMetrics()
		page := pages.GetMetricsPage(metrics)
		_, err := writer.Write([]byte(page))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (mh *MetricsHandler) GetMetricHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		name := chi.URLParam(request, "name")
		mType := chi.URLParam(request, "mType")
		log.Printf("Get metric: request vars - name: '%s', type: '%s'", name, mType)
		if name == "" {
			http.Error(writer, "metric 'name' must be specified", http.StatusBadRequest)
			return
		}
		if mType == "" {
			http.Error(writer, "metric 'mType' must be specified", http.StatusBadRequest)
			return
		}
		metric := mh.Repository.GetMetric(name)
		if metric == nil {
			http.Error(writer, fmt.Sprintf("metric was not found: %s", name), http.StatusNotFound)
			return
		}
		res, err := json.Marshal(metric.GetGenericValue())
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = writer.Write(res)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
