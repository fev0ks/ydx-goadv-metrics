package main

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	sr := repositories.GetCommonRepository()
	mh := rest.MetricsHandler{Ctx: ctx, Repository: sr}

	router := rest.NewRouter()
	rest.HandleMetricRequests(router, mh)

	internal.ProperExitDefer(&internal.ExitHandler{})
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":8080", router))
}
