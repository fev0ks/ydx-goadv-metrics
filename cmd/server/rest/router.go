package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(timerTrace)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	return router
}

func timerTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%v] Request time execution for: %s '%s' \r\n", time.Since(start), r.Method, r.RequestURI)
	})
}

func HandleMetricRequests(router chi.Router, mh *MetricsHandler) {
	router.Get("/", mh.GetMetricsHandler())
	router.Post("/update/{mType}/{id}/{value}", mh.ReceptionMetricsHandler())
	router.Post("/update/", mh.ReceptionMetricsHandler())
	router.Get("/value/{mType}/{id}", mh.GetMetricHandler())
	router.Get("/value", mh.GetMetricsHandler())
}
