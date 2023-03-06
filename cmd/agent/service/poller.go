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
	client   *resty.Client
	sender   sender.Sender
	mr       agent.IMetricRepository
	interval time.Duration
}

func NewCommonMetricPoller(
	client *resty.Client,
	metricSender sender.Sender,
	repository agent.IMetricRepository,
	pollInterval time.Duration,
) agent.MetricPoller {
	return &commonMetricPoller{
		client:   client,
		sender:   metricSender,
		mr:       repository,
		interval: pollInterval,
	}
}

func (cmp *commonMetricPoller) PollMetrics(ctx context.Context, done chan struct{}) {
	ticker := time.NewTicker(cmp.interval)
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
				metrics := cmp.mr.GetMetrics()
				err := cmp.sendMetrics(ctx, metrics)
				if err != nil {
					log.Printf("[%v] Poll metrics finished with errors: %v", time.Since(start).String(), err)
				} else {
					log.Printf("[%v] Poll metrics finished", time.Since(start).String())
				}
			}
		}
	}()
}

func (cmp *commonMetricPoller) sendMetrics(ctx context.Context, metrics []*model.Metric) error {
	log.Printf("Polling %d metrics", len(metrics))
	return cmp.sender.SendMetrics(ctx, metrics)
}
