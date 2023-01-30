package service

import (
	"testing"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"

	"github.com/stretchr/testify/assert"
)

func TestCollectMetrics(t *testing.T) {
	testCases := []struct {
		name string
	}{{
		"Should collect, save metrics and successfully stopped",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo := repositories.NewCommonMetricsRepository()
			metricFactory := NewMetricFactory("")
			collector := NewCommonMetricCollector(repo, metricFactory, 2*time.Second)
			stopChannel := make(chan struct{})
			collector.CollectMetrics(stopChannel)
			time.Sleep(5 * time.Second)
			stopChannel <- struct{}{}
			metrics := repo.GetMetrics()
			assert.True(t, 29 <= len(metrics))
			for _, metric := range metrics {
				assert.NotNil(t, metric)
			}
		})
	}
}
