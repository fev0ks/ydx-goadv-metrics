package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func NewRouter(_ context.Context) *chi.Mux {
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
