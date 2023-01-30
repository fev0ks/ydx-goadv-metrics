package sender

import (
	"encoding/json"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
	"github.com/go-resty/resty/v2"
)

type jsonSender struct {
	MetricSender
	client *resty.Client
}

func NewJSONMetricSender(client *resty.Client) MetricSender {
	sender := &jsonSender{
		&abstractMetricSender{},
		client,
	}
	return sender
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
