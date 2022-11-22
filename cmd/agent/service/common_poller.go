package service

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"io"
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
	mr       agent.MetricRepository
	host     string
	port     string
	interval time.Duration
}

func NewCommonMetricPoller(ctx context.Context, repository agent.MetricRepository, host string, port string, pollInterval time.Duration) *CommonMetricPoller {
	cmpInitOnce.Do(func() {
		cmpInstance = CommonMetricPoller{
			mpCtx:    ctx,
			mr:       repository,
			host:     host,
			port:     port,
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
			url, err := cmp.getUrl(metric)
			if err != nil {
				log.Printf("failed to send metric %v: %v\n", metric, err)
				continue
			}
			resp, err := http.Post(url, "text/plain", nil)
			if err != nil {
				log.Printf("failed to poll metric %v: %v\n", metric, err)
				continue
			}
			parseSendMetricsResponse(resp, metric)
		}
	}
}

func parseSendMetricsResponse(resp *http.Response, metric *model.Metric) {
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to read response body: %v\n", err)
		}
		log.Printf("response status is not OK %v: %v, %s\n", metric, resp.StatusCode, string(respBody))
	} else {
		log.Printf("metric was succesfully pooled: %v\n", metric)
	}
}

// "http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>"
func (cmp *CommonMetricPoller) getUrl(metric *model.Metric) (url string, err error) {
	url = fmt.Sprintf("http://%s:%s/update/%s/%s/", cmp.host, cmp.port, metric.MType, metric.Name)
	switch metric.MType {
	case model.GaugeType:
		url += fmt.Sprintf("%v", fmt.Sprintf("%f", metric.Value))
	case model.CounterType:
		url += fmt.Sprintf("%v", metric.Delta)
	default:
		err = fmt.Errorf("metric type is not supported: %v", metric.MType)
	}
	return url, err
}
