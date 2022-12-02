package service

import (
	"context"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
)

type commonMetricCollector struct {
	mcCtx    context.Context
	mr       agent.MetricRepository
	interval time.Duration
}

func NewCommonMetricCollector(
	mcCtx context.Context,
	mr agent.MetricRepository,
	interval time.Duration,
) *commonMetricCollector {
	return &commonMetricCollector{mcCtx: mcCtx, mr: mr, interval: interval}
}

func (cmr *commonMetricCollector) CollectMetrics() chan struct{} {
	ticker := time.NewTicker(cmr.interval)
	done := make(chan struct{})

	var pollCount model.CounterVT
	pollCount = 0
	go func() {
		for {
			select {
			case <-done:
				log.Println("Collect metrics interrupted!")
				ticker.Stop()
				return
			case <-ticker.C:
				start := time.Now()
				log.Println("Collect metrics start")
				pollCount += 1
				metrics := cmr.getMemStatsMetrics()
				metrics = append(metrics, cmr.getPollCounterMetric(pollCount))
				metrics = append(metrics, cmr.getRandomValueMetric(100))
				cmr.mr.SaveMetric(metrics)
				log.Printf("[%v] Collect metrics finished\n", time.Since(start).String())
			}
		}
	}()
	return done
}

func (cmr *commonMetricCollector) getPollCounterMetric(pollCount model.CounterVT) *model.Metric {
	return model.NewCounterMetric("PollCount", pollCount)
}

func (cmr *commonMetricCollector) getRandomValueMetric(dim float64) *model.Metric {
	rand.Seed(time.Now().Unix())
	randomValue := rand.Float64() * dim
	return model.NewGaugeMetric("RandomValue", model.GaugeVT(randomValue))
}

func (cmr *commonMetricCollector) getMemStatsMetrics() []*model.Metric {
	metrics := make([]*model.Metric, 0)
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	metrics = append(metrics, model.NewGaugeMetric("Alloc", model.GaugeVT(memStats.Alloc)))
	metrics = append(metrics, model.NewGaugeMetric("BuckHashSys", model.GaugeVT(memStats.BuckHashSys)))
	metrics = append(metrics, model.NewGaugeMetric("Frees", model.GaugeVT(memStats.Frees)))
	metrics = append(metrics, model.NewGaugeMetric("GCCPUFraction", model.GaugeVT(memStats.GCCPUFraction)))
	metrics = append(metrics, model.NewGaugeMetric("GCSys", model.GaugeVT(memStats.GCSys)))
	metrics = append(metrics, model.NewGaugeMetric("HeapAlloc", model.GaugeVT(memStats.HeapAlloc)))
	metrics = append(metrics, model.NewGaugeMetric("HeapIdle", model.GaugeVT(memStats.HeapIdle)))
	metrics = append(metrics, model.NewGaugeMetric("HeapInuse", model.GaugeVT(memStats.HeapInuse)))
	metrics = append(metrics, model.NewGaugeMetric("HeapObjects", model.GaugeVT(memStats.HeapObjects)))
	metrics = append(metrics, model.NewGaugeMetric("HeapReleased", model.GaugeVT(memStats.HeapReleased)))
	metrics = append(metrics, model.NewGaugeMetric("HeapSys", model.GaugeVT(memStats.HeapSys)))
	metrics = append(metrics, model.NewGaugeMetric("LastGC", model.GaugeVT(memStats.LastGC)))
	metrics = append(metrics, model.NewGaugeMetric("Lookups", model.GaugeVT(memStats.Lookups)))
	metrics = append(metrics, model.NewGaugeMetric("MCacheInuse", model.GaugeVT(memStats.MCacheInuse)))
	metrics = append(metrics, model.NewGaugeMetric("MCacheSys", model.GaugeVT(memStats.MCacheSys)))
	metrics = append(metrics, model.NewGaugeMetric("MSpanInuse", model.GaugeVT(memStats.MSpanInuse)))
	metrics = append(metrics, model.NewGaugeMetric("MSpanSys", model.GaugeVT(memStats.MSpanSys)))
	metrics = append(metrics, model.NewGaugeMetric("Mallocs", model.GaugeVT(memStats.Mallocs)))
	metrics = append(metrics, model.NewGaugeMetric("NextGC", model.GaugeVT(memStats.NextGC)))
	metrics = append(metrics, model.NewGaugeMetric("NumForcedGC", model.GaugeVT(memStats.NumForcedGC)))
	metrics = append(metrics, model.NewGaugeMetric("NumGC", model.GaugeVT(memStats.NumGC)))
	metrics = append(metrics, model.NewGaugeMetric("OtherSys", model.GaugeVT(memStats.OtherSys)))
	metrics = append(metrics, model.NewGaugeMetric("PauseTotalNs", model.GaugeVT(memStats.PauseTotalNs)))
	metrics = append(metrics, model.NewGaugeMetric("StackInuse", model.GaugeVT(memStats.StackInuse)))
	metrics = append(metrics, model.NewGaugeMetric("StackSys", model.GaugeVT(memStats.StackSys)))
	metrics = append(metrics, model.NewGaugeMetric("Sys", model.GaugeVT(memStats.Sys)))
	metrics = append(metrics, model.NewGaugeMetric("TotalAlloc", model.GaugeVT(memStats.TotalAlloc)))
	return metrics
}
