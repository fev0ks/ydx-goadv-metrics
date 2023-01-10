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
	mpCtx      context.Context
	client     *resty.Client
	withBuffer bool
	buffer     []*model.Metric
}

func NewJSONMetricSender(mpCtx context.Context, client *resty.Client, withBuffer bool) MetricSender {
	sender := &jsonSender{
		mpCtx:      mpCtx,
		client:     client,
		withBuffer: withBuffer,
	}
	if withBuffer {
		sender.buffer = make([]*model.Metric, 0, 20)
	}
	return sender
}

// SendMetric TODO tmp solution - add buffer listener
func (js *jsonSender) SendMetric(metric *model.Metric) error {
	if js.withBuffer {
		return js.sendWithBuffer(metric)
	} else {
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
}

func (js *jsonSender) sendWithBuffer(metric *model.Metric) error {
	js.buffer = append(js.buffer, metric)

	if cap(js.buffer) == len(js.buffer) {
		err := js.SendMetrics(js.buffer)
		if err != nil {
			return fmt.Errorf("cannot send batch of metrics: %v", err)
		}
		js.buffer = js.buffer[:0]
	}
	return nil
}

func (js *jsonSender) SendMetrics(metrics []*model.Metric) error {
	body, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	resp, err := js.client.R().
		SetHeader(rest.ContentType, rest.ApplicationJSON).
		SetBody(body).
		Post("/updates/")
	if err != nil {
		return err
	}
	return parseSendMetricsResponse(resp, metrics)
}
