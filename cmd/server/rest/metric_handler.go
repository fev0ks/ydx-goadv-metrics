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
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"

	"github.com/go-chi/chi/v5"
)

// MetricsHandler - хендлер для обработки запросов по метрикам
type MetricsHandler struct {
	Ctx        context.Context
	Repository repository.IMetricRepository
	HashKey    string
}

func NewMetricsHandler(ctx context.Context, repository repository.IMetricRepository, hashKey string) *MetricsHandler {
	return &MetricsHandler{Ctx: ctx, Repository: repository, HashKey: hashKey}
}

// ReceptionMetricHandler - сохранение состояния метрики,
// которое может быть представлено как в виде Json тела запроса,
// так и в виде URL параметров
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

// GetMetricsHandler - получение состояние метрик в виде html таблицы
// @Tags Таблица метрик
// @Summary Запрос состояния метрик
// @Produce json
// @Success 200 {string} html страницы, в виде таблицы метрик
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router / [get]
func (mh *MetricsHandler) GetMetricsHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Get metrics")
		metrics, err := mh.Repository.GetMetrics()
		if err != nil {
			http.Error(writer, fmt.Sprintf("failed to get metrics: %v", err), http.StatusInternalServerError)
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

// GetMetricHandler - получение состояние метрики,
// указанной в виде JSON тела POST запроса
// или в виде URL параметров GET запроса
func (mh *MetricsHandler) GetMetricHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get(rest.ContentType)
		switch contentType {
		case rest.Empty, rest.TextPlain:
			mh.getTextMetricHandler(writer, request)
		case rest.ApplicationJSON:
			mh.getJSONMetricHandler(writer, request)
		default:
			err := fmt.Errorf("Content-Type: '%s' - is not supported", contentType)
			log.Printf("failed to save metric: %v", err)
			http.Error(writer, err.Error(), http.StatusNotImplemented)
		}
	}
}

// getTextMetricHandler - получение состояние метрики
// @Tags Получение метрики
// @Summary Запрос на получение метрики
// @Produce json
// @Param mType path string true "тип метрики"
// @Param id path string true "имя метрики"
// @Success 200 {object} model.Metric
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 404 {string} string "Метрика не найдена"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Failure 501 {string} string "Запрашиваемый тип метрики не поддерживается"
// @Router /value/{mType}/{id} [get]
func (mh *MetricsHandler) getTextMetricHandler(writer http.ResponseWriter, request *http.Request) {
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
		log.Printf("failed to get metric %s from repo: %v", id, err)
		http.Error(writer, fmt.Sprintf("failed to get metric %s: %v", id, err), http.StatusInternalServerError)
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

// getJSONMetricHandler - получение состояние метрики
// @Tags Получение метрики
// @Summary Запрос на получение метрики
// @Accept  json
// @Produce json
// @Success 200 {object} model.Metric
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 404 {string} string "Метрика не найдена"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Failure 501 {string} string "Запрашиваемый тип метрики не поддерживается"
// @Router /value/ [post]
func (mh *MetricsHandler) getJSONMetricHandler(writer http.ResponseWriter, request *http.Request) {
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
		log.Printf("failed to get metric %s from repo: %v", metricToFind.ID, err)
		http.Error(writer, fmt.Sprintf("failed to get metric %s: %v", metricToFind.ID, err), http.StatusInternalServerError)
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

// ReceptionMetricsHandler - сохранение состояния списка метрик
// @Tags Обновление метрик
// @Summary Запрос на обновление списка метрик
// @Accept  json
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Failure 501 {string} string "Тип метрики не поддерживается"
// @Router /updates/ [post]
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
