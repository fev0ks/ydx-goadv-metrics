package rest

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func NewRouter(_ context.Context) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(timerTrace)
	router.Use(logging)
	return router
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// timerTrace замеряет время выполнения функции.
func timerTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// напишите код функции
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%v] duration for a request\r\n", time.Since(start))
	})
}
