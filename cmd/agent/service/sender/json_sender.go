package sender

import (
	"context"
	"encoding/json"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"

	"github.com/go-resty/resty/v2"
)

type jsonSender struct {
	msCtx  context.Context
	client *resty.Client
}

func NewJSONMetricSender(msCtx context.Context, client *resty.Client) MetricSender {
	sender := &jsonSender{
		msCtx:  msCtx,
		client: client,
	}
	return sender
}

func (js *jsonSender) SendMetric(metric *model.Metric) error {
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
