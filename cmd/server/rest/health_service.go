package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
)

type HealthChecker struct {
	ctx  context.Context
	repo server.IMetricRepository
}

func NewHealthChecker(ctx context.Context, repo server.IMetricRepository) HealthChecker {
	return HealthChecker{ctx, repo}
}

func (hc *HealthChecker) CheckDBHandler() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := hc.repo.HealthCheck(hc.ctx)
		if err != nil {
			log.Printf("failed db health check: %v", err)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
