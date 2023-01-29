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

type MetricSender interface {
	SendMetrics(ctx context.Context, metrics []*model.Metric)
}

func parseSendMetricResponse(resp *resty.Response, metric *model.Metric) error {
	if resp.StatusCode() != http.StatusOK {
		respBody := resp.Body()
		return fmt.Errorf("response status is not OK '%v': %s, body: '%s'", metric, resp.Status(), strings.TrimSpace(string(respBody)))
	} else {
		log.Printf("metric was succesfully pooled: %v", metric)
		return nil
	}
}

func parseSendMetricsResponse(resp *resty.Response, metrics []*model.Metric) error {
	if resp.StatusCode() != http.StatusOK {
		respBody := resp.Body()
		return fmt.Errorf("response status is not OK: %s, body: '%s'", resp.Status(), strings.TrimSpace(string(respBody)))
	} else {
		log.Printf("%d metrics was successfully pooled", len(metrics))
		return nil
	}
}
