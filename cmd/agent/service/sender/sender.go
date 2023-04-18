package sender

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"

	"github.com/go-resty/resty/v2"
)

type Sender interface {
	MetricSender
	MetricsSender
}

type MetricSender interface {
	SendMetric(ctx context.Context, metric *model.Metric) error
}

type MetricsSender interface {
	SendMetrics(ctx context.Context, metrics []*model.Metric) error
}

type metricsSender struct {
	MetricSender
}

func (ms *metricsSender) SendMetrics(ctx context.Context, metrics []*model.Metric) error {
	errors := make([]string, 0)
	for _, metric := range metrics {
		select {
		case <-ctx.Done():
			log.Println("Context was cancelled!")
			return nil
		default:
			err := ms.MetricSender.SendMetric(ctx, metric)
			if err != nil {
				errors = append(errors, fmt.Sprintf("{%v: %v}", metric, err))
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("failed to send metrics: %s", strings.Join(errors, "; "))
	}
	return nil
}

func parseSendMetricResponse(resp *resty.Response, metric *model.Metric) error {
	if resp.StatusCode() != http.StatusOK {
		respBody := resp.Body()
		return fmt.Errorf("response status is not OK '%v': %s, body: '%s'", metric, resp.Status(), strings.TrimSpace(string(respBody)))
	}
	log.Printf("metric was succesfully pooled: %v", metric)
	return nil
}

func parseSendMetricsResponse(resp *resty.Response, metrics []*model.Metric) error {
	if resp.StatusCode() != http.StatusOK {
		respBody := resp.Body()
		return fmt.Errorf("response status is not OK: %s, body: '%s'", resp.Status(), strings.TrimSpace(string(respBody)))
	}
	log.Printf("%d metrics was successfully pooled", len(metrics))
	return nil
}
