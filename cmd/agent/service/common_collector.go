package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
)

type commonMetricCollector struct {
	mcCtx     context.Context
	mr        agent.MetricRepository
	mf        MetricFactory
	interval  time.Duration
	pollCount model.CounterVT
}

func NewCommonMetricCollector(
	mcCtx context.Context,
	metricRepository agent.MetricRepository,
	metricFactory MetricFactory,
	interval time.Duration,
) agent.MetricCollector {
	return &commonMetricCollector{mcCtx: mcCtx, mr: metricRepository, mf: metricFactory, interval: interval, pollCount: 0}
}

func (cmr *commonMetricCollector) CollectMetrics() chan struct{} {
	ticker := time.NewTicker(cmr.interval)
	done := make(chan struct{})
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
				cmr.mr.SaveMetric(cmr.collectMetrics())
				log.Printf("[%v] Collect metrics finished\n", time.Since(start).String())
			}
		}
	}()
	return done
}

func (cmr *commonMetricCollector) collectMetrics() []*model.Metric {
	metrics := make([]*model.Metric, 0)
	metrics = append(metrics, cmr.getMemStatsMetrics()...)
	metrics = append(metrics, cmr.getPollCounterMetric())
	metrics = append(metrics, cmr.getRandomValueMetric())
	metrics = append(metrics, cmr.getGopsMetrics()...)
	return metrics
}

func (cmr *commonMetricCollector) getPollCounterMetric() *model.Metric {
	cmr.pollCount += 1
	return cmr.mf.NewCounterMetric("PollCount", cmr.pollCount)
}

func (cmr *commonMetricCollector) getRandomValueMetric() *model.Metric {
	rand.Seed(time.Now().Unix())
	randomValue := rand.Float64() * 100
	return cmr.mf.NewGaugeMetric("RandomValue", model.GaugeVT(randomValue))
}

func (cmr *commonMetricCollector) getMemStatsMetrics() []*model.Metric {
	metrics := make([]*model.Metric, 0)
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	metrics = append(metrics, cmr.mf.NewGaugeMetric("Alloc", model.GaugeVT(memStats.Alloc)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("BuckHashSys", model.GaugeVT(memStats.BuckHashSys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("Frees", model.GaugeVT(memStats.Frees)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("GCCPUFraction", model.GaugeVT(memStats.GCCPUFraction)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("GCSys", model.GaugeVT(memStats.GCSys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("HeapAlloc", model.GaugeVT(memStats.HeapAlloc)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("HeapIdle", model.GaugeVT(memStats.HeapIdle)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("HeapInuse", model.GaugeVT(memStats.HeapInuse)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("HeapObjects", model.GaugeVT(memStats.HeapObjects)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("HeapReleased", model.GaugeVT(memStats.HeapReleased)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("HeapSys", model.GaugeVT(memStats.HeapSys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("LastGC", model.GaugeVT(memStats.LastGC)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("Lookups", model.GaugeVT(memStats.Lookups)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("MCacheInuse", model.GaugeVT(memStats.MCacheInuse)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("MCacheSys", model.GaugeVT(memStats.MCacheSys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("MSpanInuse", model.GaugeVT(memStats.MSpanInuse)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("MSpanSys", model.GaugeVT(memStats.MSpanSys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("Mallocs", model.GaugeVT(memStats.Mallocs)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("NextGC", model.GaugeVT(memStats.NextGC)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("NumForcedGC", model.GaugeVT(memStats.NumForcedGC)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("NumGC", model.GaugeVT(memStats.NumGC)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("OtherSys", model.GaugeVT(memStats.OtherSys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("PauseTotalNs", model.GaugeVT(memStats.PauseTotalNs)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("StackInuse", model.GaugeVT(memStats.StackInuse)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("StackSys", model.GaugeVT(memStats.StackSys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("Sys", model.GaugeVT(memStats.Sys)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("TotalAlloc", model.GaugeVT(memStats.TotalAlloc)))
	return metrics
}

func (cmr *commonMetricCollector) getGopsMetrics() []*model.Metric {
	metrics := make([]*model.Metric, 0)
	memoryStat, _ := mem.VirtualMemory()
	metrics = append(metrics, cmr.mf.NewGaugeMetric("TotalMemory", model.GaugeVT(memoryStat.Total)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("FreeMemory", model.GaugeVT(memoryStat.Free)))
	cpuUsed, _ := cpu.Percent(time.Second*10, true)
	for i := range cpuUsed {
		metrics = append(metrics, cmr.mf.NewGaugeMetric(fmt.Sprintf("CPUutilization%d", i+1), model.GaugeVT(cpuUsed[i])))
	}
	return metrics
}
