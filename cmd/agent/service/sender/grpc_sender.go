package sender

import (
	"context"
	"strings"
	"sync"

	pb "github.com/fev0ks/ydx-goadv-metrics/internal/grpc"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
)

type grpcSender struct {
	client pb.MetricsClient
	sync.RWMutex
}

func NewGrpcMetricSender(
	client pb.MetricsClient,
) Sender {
	sender := &grpcSender{
		client: client,
	}
	return sender
}

func (s *grpcSender) SendMetric(ctx context.Context, metric *model.Metric) error {
	return s.SendMetrics(ctx, []*model.Metric{metric})
}

func (s *grpcSender) SendMetrics(ctx context.Context, metrics []*model.Metric) error {
	_, err := s.client.SaveMetrics(ctx, s.getMetricsRequest(metrics))
	if err != nil {
		return err
	}
	return nil
}

func (s *grpcSender) getMetricsRequest(metrics []*model.Metric) *pb.MetricsRequest {
	pbMetrics := make([]*pb.Metric, 0, len(metrics))
	for _, metric := range metrics {
		hash := metric.Hash
		mType := metric.MType.String()
		pbMetrics = append(pbMetrics, &pb.Metric{
			Id:    metric.ID,
			MType: pb.MetricTypes(pb.MetricTypes_value[strings.ToUpper(mType)]),
			Delta: (*uint64)(metric.Delta),
			Value: (*float64)(metric.Value),
			Hash:  &hash,
		})
	}
	return &pb.MetricsRequest{Metrics: pbMetrics}
}
