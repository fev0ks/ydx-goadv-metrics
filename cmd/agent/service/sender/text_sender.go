package sender

import (
	"context"
	"fmt"
	"log"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"

	"github.com/go-resty/resty/v2"
)

type textSender struct {
	mpCtx  context.Context
	client *resty.Client
}

func NewTextMetricSender(mpCtx context.Context, client *resty.Client) MetricSender {
	return &textSender{
		mpCtx:  mpCtx,
		client: client,
	}
}

func (ts *textSender) SendMetrics(ctx context.Context, metrics []*model.Metric) {
	for _, metric := range metrics {
		select {
		case <-ctx.Done():
			log.Println("Context was cancelled!")
			return
		default:
			err := ts.sendMetric(metric)
			if err != nil {
				log.Printf("failed to poll metric %v: %v", metric, err)
			}
		}
	}
}

func (ts *textSender) sendMetric(metric *model.Metric) error {
	value := metric.GetValue()
	if value == model.NanVal {
		return fmt.Errorf("metric type '%s' is not supported", metric.MType)
	}
	resp, err := ts.client.R().
		SetHeader(rest.ContentType, rest.TextPlain).
		SetPathParams(map[string]string{
			"mType": string(metric.MType),
			"name":  metric.ID,
			"value": value,
		}).
		Post("/update/{mType}/{name}/{value}")
	if err != nil {
		return err
	}
	return parseSendMetricResponse(resp, metric)
}
