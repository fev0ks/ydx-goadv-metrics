package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	cmpInitOnce sync.Once
	cmpInstance CommonMetricPoller
)

type CommonMetricPoller struct {
	mpCtx    context.Context
	client   *resty.Client
	mr       agent.MetricRepository
	interval time.Duration
}

func NewCommonMetricPoller(
	ctx context.Context,
	client *resty.Client,
	repository agent.MetricRepository,
	pollInterval time.Duration,
) agent.MetricPoller {
	cmpInitOnce.Do(func() {
		cmpInstance = CommonMetricPoller{
			mpCtx:    ctx,
			client:   client,
			mr:       repository,
			interval: pollInterval,
		}
	})
	return &cmpInstance
}

func (cmp *CommonMetricPoller) PollMetrics() chan bool {
	ticker := time.NewTicker(cmp.interval)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				log.Println("Poll metrics interrupted!")
				ticker.Stop()
				return
			case <-ticker.C:
				start := time.Now()
				log.Println("Poll metrics start")
				metrics := cmp.mr.GetMetricsList()
				cmp.sendMetricsAsUrl(metrics)
				log.Printf("[%v] Poll metrics finished\n", time.Now().Sub(start).String())
			}
		}
	}()
	return done
}

func (cmp *CommonMetricPoller) sendMetricsAsUrl(metrics []*model.Metric) {
	log.Printf("Polling %d metrics", len(metrics))
	for _, metric := range metrics {
		select {
		case <-cmp.mpCtx.Done():
			log.Println("Context was cancelled!")
			return
		default:
			value := metric.GetValue()
			if value == model.NanVal {
				log.Printf("failed to send metric %v - metric type is not supported", metric)
				continue
			}
			resp, err := cmp.client.R().
				SetHeader("Content-type", "text/plain").
				SetPathParams(map[string]string{
					"mType": string(metric.MType),
					"name":  metric.Name,
					"value": value,
				}).
				Post("/update/{mType}/{name}/{value}")
			if err != nil {
				log.Printf("failed to poll metric %v: %v\n", metric, err)
				continue
			}
			parseSendMetricsResponse(resp, metric)
		}
	}
}

func parseSendMetricsResponse(resp *resty.Response, metric *model.Metric) {
	if resp.StatusCode() != http.StatusOK {
		respBody := resp.Body()
		log.Printf("response status is not OK %v: %s, %s\n", metric, resp.Status(), string(respBody))
	} else {
		log.Printf("metric was succesfully pooled: %v\n", metric)
	}
}
