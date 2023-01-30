package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
	"github.com/go-resty/resty/v2"
	"log"
	"strings"
	"sync"
)

type bulkSender struct {
	client     *resty.Client
	batchLimit int
	sync.RWMutex
}

func NewBulkMetricSender(
	client *resty.Client,
	batchLimit int,
) Sender {
	sender := &bulkSender{
		client:     client,
		batchLimit: batchLimit,
	}
	return sender
}

func (s *bulkSender) SendMetrics(ctx context.Context, metrics []*model.Metric) error {
	errors := make([]string, 0)
	var err error
	metricsSize := len(metrics)
	for i := 0; i < len(metrics); i = i + s.batchLimit {
		select {
		case <-ctx.Done():
			log.Println("SendMetrics interrupted!")
			return nil
		default:
			if i+s.batchLimit > metricsSize {
				err = s.sendMetrics(metrics[i:metricsSize])
			} else {
				err = s.sendMetrics(metrics[i : i+s.batchLimit])
			}
			if err != nil {
				errors = append(errors, fmt.Sprintf("{%s}", err.Error()))
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("failed to send bulk metrics: %s", strings.Join(errors, "; "))
	}
	return nil
}

func (s *bulkSender) SendMetric(metric *model.Metric) error {
	return s.sendMetrics([]*model.Metric{metric})
}

func (s *bulkSender) sendMetrics(metrics []*model.Metric) error {
	body, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	resp, err := s.client.R().
		SetHeader(rest.ContentType, rest.ApplicationJSON).
		SetBody(body).
		Post("/updates/")
	if err != nil {
		return err
	}
	return parseSendMetricsResponse(resp, metrics)
}
