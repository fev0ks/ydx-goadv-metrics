package service

import (
	"context"
	"log"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service/sender"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"

	"github.com/go-resty/resty/v2"
)

type commonMetricPoller struct {
	mpCtx    context.Context
	client   *resty.Client
	sender   sender.MetricSender
	mr       agent.MetricRepository
	interval time.Duration
}

func NewCommonMetricPoller(
	ctx context.Context,
	client *resty.Client,
	metricSender sender.MetricSender,
	repository agent.MetricRepository,
	pollInterval time.Duration,
) agent.MetricPoller {
	return &commonMetricPoller{
		mpCtx:    ctx,
		client:   client,
		sender:   metricSender,
		mr:       repository,
		interval: pollInterval,
	}
}

func (cmp *commonMetricPoller) PollMetrics() chan struct{} {
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

func (cmp *commonMetricPoller) sendMetrics(metrics []*model.Metric) {
	log.Printf("Polling %d metrics", len(metrics))
	for _, metric := range metrics {
		select {
		case <-cmp.mpCtx.Done():
			log.Println("Context was cancelled!")
			return
		default:
			err := cmp.sender.SendMetric(metric)
			if err != nil {
				log.Printf("failed to poll metric %v: %v\n", metric, err)
			}
		}
	}
}
