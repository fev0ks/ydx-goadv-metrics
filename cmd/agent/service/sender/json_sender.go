package sender

import (
	"context"
	"encoding/json"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
	"log"

	"github.com/go-resty/resty/v2"
)

type jsonSender struct {
	client *resty.Client
}

func NewJSONMetricSender(client *resty.Client) MetricSender {
	sender := &jsonSender{
		client: client,
	}
	return sender
}

func (js *jsonSender) SendMetrics(ctx context.Context, metrics []*model.Metric) {
	for _, metric := range metrics {
		select {
		case <-ctx.Done():
			log.Println("Context was cancelled!")
			return
		default:
			err := js.sendMetric(metric)
			if err != nil {
				log.Printf("failed to poll metric %v: %v", metric, err)
			}
		}
	}
}

func (js *jsonSender) sendMetric(metric *model.Metric) error {
	body, err := json.Marshal(*metric)
	if err != nil {
		return err
	}
	resp, err := js.client.R().
		SetHeader(rest.ContentType, rest.ApplicationJSON).
		SetBody(body).
		Post("/update/")
	if err != nil {
		return err
	}
	return parseSendMetricResponse(resp, metric)
}
