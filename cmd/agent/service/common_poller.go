package service

import (
	"context"
	"fmt"
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

func (cmp *CommonMetricPoller) PollMetrics() chan struct{} {
	ticker := time.NewTicker(cmp.interval)
	done := make(chan struct{})
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
				cmp.sendMetrics(metrics)
				log.Printf("[%v] Poll metrics finished\n", time.Since(start).String())
			}
		}
	}()
	return done
}

func (cmp *CommonMetricPoller) sendMetrics(metrics []*model.Metric) {
	log.Printf("Polling %d metrics", len(metrics))
	for _, metric := range metrics {
		select {
		case <-cmp.mpCtx.Done():
			log.Println("Context was cancelled!")
			return
		default:
			err := cmp.SendMetric(metric)
			if err != nil {
				log.Printf("failed to poll metric %v: %v\n", metric, err)
			}
		}
	}
}

func (cmp *CommonMetricPoller) SendMetric(metric *model.Metric) error {
	value := metric.GetValue()
	if value == model.NanVal {
		return fmt.Errorf("metric type '%s' is not supported", metric.MType)
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
		return err
	}
	err = parseSendMetricResponse(resp, metric)
	if err != nil {
		return err
	}
	return nil
}

func parseSendMetricResponse(resp *resty.Response, metric *model.Metric) error {
	if resp.StatusCode() != http.StatusOK {
		respBody := resp.Body()
		return fmt.Errorf("response status is not OK '%v': %s, body: '%s'", metric, resp.Status(), string(respBody))
	} else {
		log.Printf("metric was succesfully pooled: %v\n", metric)
		return nil
	}
}
