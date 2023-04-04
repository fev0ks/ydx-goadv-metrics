package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/middlewares"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
)

// NewRouter - инициализация обьекта обработки запросов и настройка посредников
func NewRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middlewares.TimerTrace)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(3, rest.ApplicationJSON, rest.TextPlain, rest.TextHTML))
	router.Use(middlewares.Decompress)
	return router
}

// HandleMetricRequests - настройка хендлеров для работы с метриками
func HandleMetricRequests(router chi.Router, mh *MetricsHandler) {
	router.Get("/", mh.GetMetricsHandler())
	router.Post("/update/{mType}/{id}/{value}", mh.ReceptionMetricHandler())
	router.Get("/value/{mType}/{id}", mh.GetMetricHandler())
	router.Post("/value/", mh.GetMetricHandler())
}

func HandleEncryptedMetricRequests(router chi.Router, mh *MetricsHandler, decrypter *middlewares.Decrypter) {
	router.Group(func(r chi.Router) {
		r.Use(decrypter.Decrypt)
		r.Post("/update/", mh.ReceptionMetricHandler())
		r.Post("/updates/", mh.ReceptionMetricsHandler())
	})
}

// HandleHeathCheck - настройка хендлеров для проверки состояния сервиса
func HandleHeathCheck(router chi.Router, hc HealthChecker) {
	router.Get("/ping", hc.CheckDBHandler())
}

// HandlePprof - настройка хендлеров для работы с профайлером
func HandlePprof(router chi.Router) {
	router.Mount("/debug", middleware.Profiler())
}
