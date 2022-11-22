package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest/handlers"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	var sr server.MetricRepository
	sr = repositories.GetCommonRepository(&ctx)
	mh := handlers.MetricsHandler{Ctx: ctx, Repository: sr}

	router := rest.NewRouter(ctx)
	rest.HandleMetricRequests(router, mh)

	log.Println("Server started")
	go internal.ProperExitDefer(&internal.ExitHandler{})
	log.Fatal(http.ListenAndServe(":8080", router))
}
