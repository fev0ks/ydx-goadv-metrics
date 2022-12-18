package rest

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter() chi.Router {
	router := chi.NewRouter()
	return router
}
