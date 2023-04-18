package servers

import (
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/configs"
	pb "github.com/fev0ks/ydx-goadv-metrics/internal/grpc"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"
)

func StartGrpcServer(
	address string,
	repo repository.IMetricRepository,
	hashKey string,
) *grpc.Server {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	// создаём gRPC-сервер без зарегистрированной службы
	grpcServer := grpc.NewServer()
	// регистрируем сервис
	pb.RegisterMetricsServer(grpcServer, NewGrpcServer(repo, hashKey))

	go func() {
		// получаем запрос gRPC
		log.Printf("Grpc Server started on %s", address)
		if err := grpcServer.Serve(listen); err != nil {
			log.Printf("Grpc Server closed with msg: '%v'", err)
		}
	}()
	return grpcServer
}

func StartHttpServer(appConfig *configs.AppConfig, router chi.Router) *http.Server {
	server := &http.Server{Addr: appConfig.ServerAddress, Handler: router}
	go func() {
		log.Printf("Http Server started on %s", appConfig.ServerAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Http Server closed with msg: '%v'", err)
		}
	}()
	return server
}
