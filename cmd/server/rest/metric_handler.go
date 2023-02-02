package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/pages"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"

	"github.com/go-chi/chi/v5"
)

type MetricsHandler struct {
	Ctx        context.Context
	Repository server.MetricRepository
	HashKey    string
}

func NewMetricsHandler(ctx context.Context, repository server.MetricRepository, hashKey string) *MetricsHandler {
	return &MetricsHandler{Ctx: ctx, Repository: repository, HashKey: hashKey}
}

func (mh *MetricsHandler) ReceptionMetricHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get(rest.ContentType)
		switch contentType {
		case rest.Empty, rest.TextPlain:
			mh.receptionTextMetricsHandler(writer, request)
		case rest.ApplicationJSON:
			mh.receptionJSONMetricsHandler(writer, request)
		default:
			err := fmt.Errorf("Content-Type: '%s' - is not supported", contentType)
			log.Printf("failed to save metric: %v", err)
			http.Error(writer, err.Error(), http.StatusNotImplemented)
		}
	}
}

func (mh *MetricsHandler) receptionTextMetricsHandler(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	mType := chi.URLParam(request, "mType")
	value := chi.URLParam(request, "value")
	log.Printf("request vars - id: '%s', type: '%s', value: '%s'", id, mType, value)
	if id == "" {
		http.Error(writer, "metric 'id' must be specified", http.StatusBadRequest)
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
	metric, err := model.ParseMetric(id, model.MTypeValueOf(mType), value)
	if err != nil {
		log.Printf("failed to parse metric request: %v", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if metric.MType == model.NanType {
		err = fmt.Errorf("type '%s' is not supported", mType)
		log.Printf("failed to save metric: %v", err)
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	err = mh.Repository.SaveMetric(metric)
	if err != nil {
		log.Printf("failed to save metric: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (mh *MetricsHandler) receptionJSONMetricsHandler(writer http.ResponseWriter, request *http.Request) {
	var metric *model.Metric
	body, _ := io.ReadAll(request.Body)
	defer request.Body.Close()

	err := json.Unmarshal(body, &metric)
	if err != nil {
		log.Printf("failed to parse metric request: %v", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if metric.MType == model.NanType {
		err = fmt.Errorf("type '%s' is not supported", metric.MType)
		log.Printf("failed to save metric: %v", err)
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	err = metric.CheckHash(mh.HashKey)
	if err != nil {
		log.Printf("failed to check metric hash: %v", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = mh.Repository.SaveMetric(metric)
	if err != nil {
		log.Printf("failed to save metric: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (mh *MetricsHandler) GetMetricsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Get metrics")
		metrics, err := mh.Repository.GetMetrics()
		if err != nil {
			http.Error(writer, fmt.Sprintf("failed to get metrics: %v", err), http.StatusNotFound)
			return
		}
		for _, metric := range metrics {
			metric.UpdateHash(mh.HashKey)
		}
		page := pages.GetMetricsPage(metrics)
		writer.Header().Add(rest.ContentType, rest.TextHTML)
		_, err = writer.Write([]byte(page))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (mh *MetricsHandler) GetMetricHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get(rest.ContentType)
		switch contentType {
		case rest.Empty, rest.TextPlain:
			mh.GetTextMetricHandler(writer, request)
		case rest.ApplicationJSON:
			mh.GetJSONMetricHandler(writer, request)
		default:
			err := fmt.Errorf("Content-Type: '%s' - is not supported", contentType)
			log.Printf("failed to save metric: %v", err)
			http.Error(writer, err.Error(), http.StatusNotImplemented)
		}
	}
}

func (mh *MetricsHandler) GetTextMetricHandler(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")
	mType := chi.URLParam(request, "mType")
	log.Printf("Get metric: request vars - id: '%s', type: '%s'", id, mType)
	if id == "" {
		http.Error(writer, "metric 'id' must be specified", http.StatusBadRequest)
		return
	}
	if mType == "" {
		http.Error(writer, "metric 'mType' must be specified", http.StatusBadRequest)
		return
	}
	metric, err := mh.Repository.GetMetric(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("failed to get metric %s: %v", id, err), http.StatusNotFound)
		return
	}
	if metric == nil {
		http.Error(writer, fmt.Sprintf("metric was not found: %s", id), http.StatusNotFound)
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

func (mh *MetricsHandler) GetJSONMetricHandler(writer http.ResponseWriter, request *http.Request) {
	var metricToFind *model.Metric
	body, _ := io.ReadAll(request.Body)
	defer request.Body.Close()

	err := json.Unmarshal(body, &metricToFind)
	if err != nil {
		log.Printf("failed to parse metric request '%s': %v", body, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if metricToFind.MType == model.NanType {
		err = fmt.Errorf("type '%s' is not supported", metricToFind.MType)
		log.Printf("failed to get metric: %v", err)
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	metric, err := mh.Repository.GetMetric(metricToFind.ID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("failed to get metric %s: %v", metricToFind.ID, err), http.StatusNotFound)
		return
	}
	if metric == nil {
		http.Error(writer, fmt.Sprintf("metric was not found: %s", metricToFind.ID), http.StatusNotFound)
		return
	}
	metric.UpdateHash(mh.HashKey)
	res, err := json.Marshal(metric)
	if err != nil {
		log.Printf("failed to marshal metric: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Add(rest.ContentType, rest.ApplicationJSON)
	_, err = writer.Write(res)
	if err != nil {
		log.Printf("failed to write metric response: %v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (mh *MetricsHandler) ReceptionMetricsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var metrics []*model.Metric
		body, _ := io.ReadAll(request.Body)
		defer request.Body.Close()

		err := json.Unmarshal(body, &metrics)
		if err != nil {
			log.Printf("failed to parse metric request: %v", err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		for _, metric := range metrics {
			if metric.MType == model.NanType {
				err = fmt.Errorf("type '%s' is not supported", metric.MType)
				log.Printf("failed to save metric '%s': %v", metric.ID, err)
				http.Error(writer, err.Error(), http.StatusNotImplemented)
				return
			}
			err = metric.CheckHash(mh.HashKey)
			if err != nil {
				log.Printf("failed to check metric '%s' hash: %v", metric.ID, err)
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
			err = mh.Repository.SaveMetric(metric)
			if err != nil {
				log.Printf("failed to save metric '%s': %v", metric.ID, err)
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			writer.WriteHeader(http.StatusOK)
		}
	}
}
