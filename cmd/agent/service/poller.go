package service

import (
	"context"
	"log"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/service/sender"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"github.com/fev0ks/ydx-goadv-metrics/internal/shutdown"

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

func (cmp *commonMetricPoller) PollMetrics(ctx context.Context, eh *shutdown.ExitHandler, done chan struct{}) {
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
				err := cmp.forceSendMetrics(ctx, eh, metrics)
				if err != nil {
					log.Printf("[%v] Poll metrics finished with errors: %v", time.Since(start).String(), err)
				} else {
					log.Printf("[%v] Poll metrics finished", time.Since(start).String())
				}
			}
		}
	}()
}

func (cmp *commonMetricPoller) forceSendMetrics(ctx context.Context, eh *shutdown.ExitHandler, metrics []*model.Metric) error {
	alias := "sendMetrics"
	if eh.IsNewFuncExecutionAllowed() {
		eh.AddFuncInProcessing(alias)
		defer eh.FuncFinished(alias)
		return cmp.sendMetrics(ctx, metrics)
	} else {
		log.Println("System is going to shutdown, new func execution are rejected!")
	}
	return nil
}

func (cmp *commonMetricPoller) sendMetrics(ctx context.Context, metrics []*model.Metric) error {
	log.Printf("Polling %d metrics", len(metrics))
	return cmp.sender.SendMetrics(ctx, metrics)
}
