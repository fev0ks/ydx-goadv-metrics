package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/agent"
)

type commonMetricCollector struct {
	mr        agent.MetricRepository
	mf        MetricFactory
	interval  time.Duration
	pollCount uint64
}

func NewCommonMetricCollector(
	metricRepository agent.MetricRepository,
	metricFactory MetricFactory,
	interval time.Duration,
) agent.MetricCollector {
	return &commonMetricCollector{
		mr:        metricRepository,
		mf:        metricFactory,
		interval:  interval,
		pollCount: 0,
	}
}

func (cmr *commonMetricCollector) CollectMetrics(ctx context.Context, done chan struct{}) {
	colCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-done
		log.Println("Collect metrics interrupted!")
		cancel()
	}()
	cmr.collectCommonMetrics(colCtx)
	cmr.collectGopsMetrics(colCtx)
}

func (cmr *commonMetricCollector) collectCommonMetrics(ctx context.Context) {
	ticker := time.NewTicker(cmr.interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Collect common metrics interrupted!")
				ticker.Stop()
				return
			case <-ticker.C:
				cmr.processPollCounterMetric()
				cmr.processRandomValueMetric()
				cmr.processMemStatsMetrics()
			}
		}
	}()
}

func (cmr *commonMetricCollector) processPollCounterMetric() {
	cmr.mr.SaveMetric(cmr.mf.NewCounterMetric("PollCount", model.CounterVT(atomic.AddUint64(&cmr.pollCount, 1))))
	log.Printf("PollCount metric finished")
}

func (cmr *commonMetricCollector) processRandomValueMetric() {
	rand.Seed(time.Now().Unix())
	randomValue := rand.Float64() * 100
	cmr.mr.SaveMetric(cmr.mf.NewGaugeMetric("RandomValue", model.GaugeVT(randomValue)))
	log.Printf("RandomValue metric finished")
}

func (cmr *commonMetricCollector) processMemStatsMetrics() {
	metrics := make([]*model.Metric, 0, 19)
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
	cmr.mr.SaveMetrics(metrics)
	log.Printf("MemStats metrics finished")
}

func (cmr *commonMetricCollector) collectGopsMetrics(ctx context.Context) {
	ticker := time.NewTicker(cmr.interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Collect gops metrics interrupted!")
				ticker.Stop()
				return
			case <-ticker.C:
				// processGopsMetrics takes 10 sec
				go cmr.processGopsMetrics()
			}
		}
	}()
}

func (cmr *commonMetricCollector) processGopsMetrics() {
	metrics := make([]*model.Metric, 0, 3)
	memoryStat, _ := mem.VirtualMemory()
	metrics = append(metrics, cmr.mf.NewGaugeMetric("TotalMemory", model.GaugeVT(memoryStat.Total)))
	metrics = append(metrics, cmr.mf.NewGaugeMetric("FreeMemory", model.GaugeVT(memoryStat.Free)))
	cpuUsed, _ := cpu.Percent(time.Second*10, true)
	for i := range cpuUsed {
		metrics = append(metrics, cmr.mf.NewGaugeMetric(fmt.Sprintf("CPUutilization%d", i+1), model.GaugeVT(cpuUsed[i])))
	}
	cmr.mr.SaveMetrics(metrics)
	log.Printf("Gops metrics finished")
}
