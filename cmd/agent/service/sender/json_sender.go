package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"

	"github.com/go-resty/resty/v2"
)

type jsonSender struct {
	mpCtx  context.Context
	client *resty.Client
}

func NewJSONMetricSender(mpCtx context.Context, client *resty.Client) MetricSender {
	return &jsonSender{
		mpCtx:  mpCtx,
		client: client,
	}
}

func (js jsonSender) SendMetric(metric *model.Metric) error {
	value := metric.GetValue()
	if value == model.NanVal {
		return fmt.Errorf("metric type '%s' is not supported", metric.MType)
	}
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
