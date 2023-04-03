package sender

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/rest"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	consts "github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
)

type jsonSender struct {
	client    *resty.Client
	encryptor *rest.Encryptor
}

func NewJSONMetricSender(client *resty.Client, encryptor *rest.Encryptor) Sender {
	sender := &metricsSender{
		&jsonSender{
			client,
			encryptor,
		},
	}
	return sender
}

func (js *jsonSender) SendMetric(metric *model.Metric) error {
	body, err := json.Marshal(*metric)
	if err != nil {
		return err
	}
	encryptedBody, err := js.encryptor.Encrypt(body)
	if err != nil {
		return err
	}
	resp, err := js.client.R().
		SetHeader(consts.ContentType, consts.ApplicationJSON).
		SetBody(encryptedBody).
		Post("/update/")
	if err != nil {
		return err
	}
	return parseSendMetricResponse(resp, metric)
}
