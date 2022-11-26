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
	router.Use(timerTraceMdw)
	router.Use(recoverMdw)
	return router
}

func recoverMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered. Error:%v\n", r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func timerTraceMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%v] Request time execution for: %s '%s' \r\n", time.Since(start), r.Method, r.RequestURI)
	})
}
