package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/rest"
	"github.com/fev0ks/ydx-goadv-metrics/internal"
)

func main() {
	ctx := context.Background()

	sr := repositories.NewCommonRepository()
	mh := rest.NewMetricsHandler(ctx, sr)

	router := rest.NewRouter()
	rest.HandleMetricRequests(router, mh)

	internal.ProperExitDefer(&internal.ExitHandler{})
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":8080", router))
}
