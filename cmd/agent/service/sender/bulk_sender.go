package sender

import (
	"context"
	"encoding/json"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
	"github.com/go-resty/resty/v2"
	"log"
	"sync"
	"time"
)

type bulkSender struct {
	msCtx          context.Context
	client         *resty.Client
	metrics        []*model.Metric
	batchLimit     int
	mutex          sync.RWMutex
	batchMetricsCh chan []*model.Metric
}

func NewBulkMetricSender(
	msCtx context.Context,
	client *resty.Client,
	batchLimit int,
	sendInterval time.Duration,
) MetricSender {
	sender := &bulkSender{
		msCtx:          msCtx,
		client:         client,
		batchLimit:     batchLimit,
		metrics:        make([]*model.Metric, 0, batchLimit),
		batchMetricsCh: make(chan []*model.Metric),
	}
	go sender.startSenderListener(msCtx, sendInterval)
	return sender
}

func (s *bulkSender) SendMetric(metric *model.Metric) error {
	s.mutex.Lock()
	s.metrics = append(s.metrics, metric)
	if len(s.metrics) >= s.batchLimit {
		log.Printf("Flush batch of '%v' metrics\n", len(s.metrics))
		s.batchMetricsCh <- s.metrics
		s.metrics = make([]*model.Metric, 0, s.batchLimit)
	}
	s.mutex.Unlock()
	return nil
}

func (s *bulkSender) startSenderListener(ctx context.Context, sendInterval time.Duration) {
	ticker := time.NewTicker(sendInterval)
	for {
		var metrics []*model.Metric
		select {
		case <-ticker.C:
			metrics = s.uploadMetrics()
			log.Printf("Upload '%v' metrics\n", len(metrics))
		case metrics = <-s.batchMetricsCh:
			log.Printf("Handle '%v' metrics\n", len(metrics))
		case <-ctx.Done():
			ticker.Stop()
			return
		}
		if len(metrics) != 0 {
			log.Printf("Sending '%v' metrics\n", len(metrics))
			s.sendMetricsAsync(metrics)
		}
	}
}

func (s *bulkSender) sendMetricsAsync(metrics []*model.Metric) {
	go func() {
		err := s.sendMetrics(metrics)
		if err != nil {
			log.Printf("failed to send metrics %v: %v\n", metrics, err)
		}
	}()
}

func (s *bulkSender) sendMetrics(metrics []*model.Metric) error {
	body, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	resp, err := s.client.R().
		SetHeader(rest.ContentType, rest.ApplicationJSON).
		SetBody(body).
		Post("/updates/")
	if err != nil {
		return err
	}
	return parseSendMetricsResponse(resp, metrics)
}

func (s *bulkSender) uploadMetrics() []*model.Metric {
	s.mutex.Lock()
	metrics := s.metrics
	s.metrics = make([]*model.Metric, 0, s.batchLimit)
	s.mutex.Unlock()
	return metrics
}
