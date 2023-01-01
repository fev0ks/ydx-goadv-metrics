package sender

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"

	"github.com/go-resty/resty/v2"
)

type MetricSender interface {
	SendMetric(metrics *model.Metric) error
}

func parseSendMetricResponse(resp *resty.Response, metric *model.Metric) error {
	if resp.StatusCode() != http.StatusOK {
		respBody := resp.Body()
		return fmt.Errorf("response status is not OK '%v': %s, body: '%s'", metric, resp.Status(), strings.TrimSpace(string(respBody)))
	} else {
		log.Printf("metric was succesfully pooled: %v\n", metric)
		return nil
	}
}
