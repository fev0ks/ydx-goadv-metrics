package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/pages"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type MetricsHandler struct {
	Ctx        context.Context
	Repository server.MetricRepository
}

func (mh *MetricsHandler) ReceptionMetricsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		name := chi.URLParam(request, "name")
		mType := chi.URLParam(request, "mType")
		value := chi.URLParam(request, "value")
		log.Printf("request vars - name: '%s', type: '%s', value: '%s'", name, mType, value)
		metric, err := model.NewMetric(name, model.MTypeValueOf(mType), value)
		if err != nil {
			log.Printf("failed to parse metric request: %v\n", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if metric.MType == model.NanType {
			err = fmt.Errorf("type is not supported: %v", metric)
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

		if metric := mh.Repository.GetMetric(name); metric == nil {
			http.Error(writer, fmt.Sprintf("metric was not found: %s", name), http.StatusNotFound)
		} else {
			res, err := json.Marshal(metric.GetValue())
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
}
