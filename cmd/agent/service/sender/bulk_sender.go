package sender

import (
	"context"
	"encoding/json"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/consts/rest"
	"github.com/go-resty/resty/v2"
	"log"
	"sync"
)

type bulkSender struct {
	client     *resty.Client
	metrics    []*model.Metric
	batchLimit int
	sync.RWMutex
	batchMetricsCh chan []*model.Metric
}

func NewBulkMetricSender(
	msCtx context.Context,
	client *resty.Client,
	batchLimit int,
) MetricSender {
	sender := &bulkSender{
		client:         client,
		batchLimit:     batchLimit,
		metrics:        make([]*model.Metric, 0, batchLimit),
		batchMetricsCh: make(chan []*model.Metric),
	}
	sender.startSendMetricsListener(msCtx)
	return sender
}

func (s *bulkSender) SendMetrics(ctx context.Context, metrics []*model.Metric) {
	s.Lock()
	defer s.Unlock()

	metricsSize := len(metrics)
	for i := 0; i < len(metrics); i = i + s.batchLimit {
		select {
		case <-ctx.Done():
			log.Println("SendMetrics interrupted!")
			return
		default:
			if i+s.batchLimit > metricsSize {
				s.batchMetricsCh <- metrics[i:metricsSize]
			} else {
				s.batchMetricsCh <- metrics[i : i+s.batchLimit]
			}
		}
	}
}

// Looks like useless for current project specific :((((
//func (s *bulkSender) startSenderListener(ctx context.Context, sendInterval time.Duration) {
//	ticker := time.NewTicker(sendInterval)
//	for {
//		var metrics []*model.Metric
//		select {
//		case <-ticker.C:
//			metrics = s.uploadMetrics()
//			log.Printf("Upload '%v' metrics", len(metrics))
//		case metrics = <-s.batchMetricsCh:
//			log.Printf("Handle '%v' metrics", len(metrics))
//		case <-ctx.Done():
//			ticker.Stop()
//			return
//		}
//		if len(metrics) != 0 {
//			log.Printf("Sending '%v' metrics", len(metrics))
//			s.sendMetricsAsync(metrics)
//		}
//	}
//}
//
//func (s *bulkSender) uploadMetrics() []*model.Metric {
//	s.Lock()
//	metrics := s.metrics
//	s.metrics = make([]*model.Metric, 0, s.batchLimit)
//	s.Unlock()
//	return metrics
//}

func (s *bulkSender) startSendMetricsListener(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("SendMetricsListener interrupted!")
				return
			case metrics := <-s.batchMetricsCh:
				log.Printf("Sending bulk %d metrics", len(metrics))
				err := s.sendMetrics(metrics)
				if err != nil {
					log.Printf("failed to send batch metrics %v: %v", metrics, err)
				}
			}
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
