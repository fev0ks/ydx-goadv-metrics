package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/pages"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts"
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
		contentType := request.Header.Get(consts.ContentType)
		switch contentType {
		case consts.TextPlain:
			mh.receptionTextMetricsHandler(writer, request)
		case consts.ApplJson:
			mh.receptionJsonMetricsHandler(writer, request)
		default:
			err := fmt.Errorf("Content-Type: '%s' - is not supported", contentType)
			log.Printf("failed to save metric: %v\n", err)
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
	metric, err := model.NewMetric(id, model.MTypeValueOf(mType), value)
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

func (mh *MetricsHandler) receptionJsonMetricsHandler(writer http.ResponseWriter, request *http.Request) {
	var metric *model.Metric
	body, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	err := json.Unmarshal(body, &metric)
	if err != nil {
		log.Printf("failed to parse metric request: %v\n", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if metric.MType == model.NanType {
		err = fmt.Errorf("type '%s' is not supported", metric.MType)
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
		contentType := request.Header.Get(consts.ContentType)
		switch contentType {
		case consts.TextPlain:
			mh.GetTextMetricHandler(writer, request)
		case consts.ApplJson:
			mh.receptionJsonMetricsHandler(writer, request)
		default:
			err := fmt.Errorf("Content-Type: '%s' - is not supported", contentType)
			log.Printf("failed to save metric: %v\n", err)
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
	metric := mh.Repository.GetMetric(id)
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

func (mh *MetricsHandler) GetJsonMetricHandler(writer http.ResponseWriter, request *http.Request) {
	var metricToFind *model.Metric
	body, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	err := json.Unmarshal(body, metricToFind)
	if err != nil {
		log.Printf("failed to parse metric request: %v\n", err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if metricToFind.MType == model.NanType {
		err = fmt.Errorf("type '%s' is not supported", metricToFind.MType)
		log.Printf("failed to get metric: %v\n", err)
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	metric := mh.Repository.GetMetric(metricToFind.ID)
	if metric == nil {
		http.Error(writer, fmt.Sprintf("metric was not found: %s", metricToFind.ID), http.StatusNotFound)
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
