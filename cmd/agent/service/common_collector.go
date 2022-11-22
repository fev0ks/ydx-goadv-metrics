package service

import (
	"context"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var (
	cmcInitOnce sync.Once
	cmcInstance CommonMetricCollector
)

type CommonMetricCollector struct {
	mcCtx    context.Context
	mr       agent.MetricRepository
	interval time.Duration
}

func NewCommonMetricCollector(ctx context.Context, metricRepository agent.MetricRepository, interval time.Duration) *CommonMetricCollector {
	cmcInitOnce.Do(func() {
		cmcInstance = CommonMetricCollector{mcCtx: ctx, mr: metricRepository, interval: interval}
	})
	return &cmcInstance
}

func (cmr *CommonMetricCollector) CollectMetrics() chan bool {
	ticker := time.NewTicker(cmr.interval)
	done := make(chan bool)

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
				log.Printf("[%v] Collect metrics finished\n", time.Now().Sub(start).String())
			}
		}
	}()
	return done
}

func (cmr *CommonMetricCollector) getPollCounterMetric(pollCount model.CounterVT) *model.Metric {
	return model.NewCounterMetric("PollCount", pollCount)
}

func (cmr *CommonMetricCollector) getRandomValueMetric(dim float64) *model.Metric {
	rand.Seed(time.Now().Unix())
	randomValue := rand.Float64() * dim
	return model.NewGaugeMetric("RandomValue", model.GaugeVT(randomValue))
}

func (cmr *CommonMetricCollector) getMemStatsMetrics() []*model.Metric {
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
