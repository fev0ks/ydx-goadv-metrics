package servers

import (
	"context"
	"fmt"
	"log"
	// импортируем пакет со сгенерированными protobuf-файлами
	pb "github.com/fev0ks/ydx-goadv-metrics/internal/grpc"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"
)

type Server struct {
	pb.UnimplementedMetricsServer
	Repository repository.IMetricRepository
	HashKey    string
}

func NewGrpcServer(
	repo repository.IMetricRepository,
	hashKey string,
) *Server {
	return &Server{
		Repository: repo,
		HashKey:    hashKey,
	}
}

func (s *Server) SaveMetrics(ctx context.Context, in *pb.MetricsRequest) (*pb.Empty, error) {

	for _, pbMetric := range in.GetMetrics() {

		metric := &model.Metric{
			ID:    pbMetric.Id,
			MType: model.MTypeValueOf(pb.MetricTypes_name[int32(pbMetric.MType)]),
			Delta: (*model.CounterVT)(pbMetric.Delta),
			Value: (*model.GaugeVT)(pbMetric.Value),
		}
		if hash := pbMetric.Hash; hash == nil {
			metric.Hash = ""
		} else {
			metric.Hash = *hash
		}
		if metric.MType == model.NanType {
			err := fmt.Errorf("type '%s' is not supported", metric.MType)
			log.Printf("failed to save metric '%s': %v", metric.ID, err)
			return nil, err
		}

		if err := metric.CheckHash(s.HashKey); err != nil {
			log.Printf("failed to check metric '%s' hash: %v", metric.ID, err)
			return nil, err
		}

		if err := s.Repository.SaveMetric(ctx, metric); err != nil {
			log.Printf("failed to save metric '%s': %v", metric.ID, err)
			return nil, err
		}
	}
	log.Printf("Saved %d metrics throught grpc", len(in.GetMetrics()))
	return &pb.Empty{}, nil
}
