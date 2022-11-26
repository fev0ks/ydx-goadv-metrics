package service

import (
	"context"
	"fmt"
	"github.com/fev0ks/ydx-goadv-metrics/cmd/agent/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost"
)

// TODO looks very bad...use some handler stubs instead of httpmock?
func TestPollMetrics(t *testing.T) {
	testCases := []struct {
		name      string
		metrics   []*model.Metric
		pollCount int
	}{
		{
			name: "Poller should send request and stop working properly",
			metrics: []*model.Metric{
				model.NewCounterMetric("poolMetric1", 1),
				model.NewCounterMetric("poolMetric2", 2),
			},
			pollCount: 1,
		},
		{
			name: "Poller should not send request by NaN metric type",
			metrics: []*model.Metric{
				model.NewNanMetric("Nan"),
			},
			pollCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client := resty.New().SetBaseURL(baseURL)
			httpmock.ActivateNonDefault(client.GetClient())
			repo := repositories.CommonMetricRepository{}
			urls := getUrlsOfMetrics(baseURL, tc.metrics)
			if tc.pollCount != 0 {
				for _, url := range urls {
					httpmock.RegisterResponder(
						"POST",
						url,
						httpmock.NewStringResponder(http.StatusOK, ""))
				}
			}

			repo.SaveMetric(tc.metrics)
			metricPoller := NewCommonMetricPoller(ctx, client, &repo, 2*time.Second)
			stopChannel := metricPoller.PollMetrics()
			time.Sleep(3 * time.Second)
			stopChannel <- struct{}{}
			time.Sleep(3 * time.Second)
			callCountInfo := httpmock.GetCallCountInfo()
			for _, url := range urls {
				count := callCountInfo[fmt.Sprintf("POST %s", url)]
				assert.Equal(t, tc.pollCount, count)
			}
		})
	}
}

func TestPollMetricsWithCancelContext(t *testing.T) {
	testCases := []struct {
		name    string
		metrics []*model.Metric
	}{
		{
			name: "Poller should not send request if context Done",
			metrics: []*model.Metric{
				model.NewCounterMetric("context_done", 1),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			client := resty.New().SetBaseURL(baseURL)
			httpmock.ActivateNonDefault(client.GetClient())
			repo := repositories.CommonMetricRepository{}
			repo.SaveMetric(tc.metrics)

			urls := getUrlsOfMetrics(baseURL, tc.metrics)
			metricPoller := NewCommonMetricPoller(ctx, client, &repo, 2*time.Second)
			cancel()
			_ = metricPoller.PollMetrics()
			time.Sleep(3 * time.Second)
			callCountInfo := httpmock.GetCallCountInfo()
			for _, url := range urls {
				count := callCountInfo[fmt.Sprintf("POST %s", url)]
				assert.Equal(t, 0, count)
			}
		})
	}
}

func getUrlsOfMetrics(baseURL string, metrics []*model.Metric) []string {
	urls := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		url := fmt.Sprintf("%s/update/%s/%s/%s",
			baseURL, metric.MType, metric.Name, metric.GetValue())
		urls = append(urls, url)
	}
	return urls
}

func TestSendMetric(t *testing.T) {
	type want struct {
		code int
		msg  string
		url  string
	}
	type metric struct {
		name  string
		mType string
		value string
	}
	testCases := []struct {
		name   string
		metric metric
		want   want
	}{
		{
			name:   "OK response for Counter type",
			metric: metric{"send_metric1", string(model.CounterType), "1"},
			want: want{
				code: http.StatusOK,
				msg:  "",
			},
		},
		{
			name:   "OK response for Gauge type",
			metric: metric{"send_metric2", string(model.GaugeType), "1.1"},
			want: want{
				code: http.StatusOK,
				msg:  "",
			},
		},
		{
			name:   "Not OK response",
			metric: metric{"send_metric3", string(model.CounterType), "1"},
			want: want{
				code: http.StatusInternalServerError,
				msg:  "response status is not OK 'Name: send_metric3, Type: counter, Value: 1': 500, body: ''",
			},
		},
		{
			name:   "Bad request",
			metric: metric{"", string(model.GaugeType), "1123"},
			want: want{
				code: http.StatusBadRequest,
				msg:  "response status is not OK 'Name: , Type: gauge, Value: 1123.000000': 400, body: ''",
			},
		},
		{
			name:   "NaN metric value",
			metric: metric{"send_metric4", string(model.NanType), ""},
			want:   want{msg: "metric type 'NaN' is not supported"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client := resty.New().SetBaseURL(baseURL)
			httpmock.ActivateNonDefault(client.GetClient())
			metric, err := model.NewMetric(tc.metric.name, model.MetricType(tc.metric.mType), tc.metric.value)
			require.NoError(t, err)
			url := fmt.Sprintf("%s/update/%s/%s/%s",
				client.BaseURL, metric.MType, metric.Name, metric.GetValue())
			httpmock.RegisterResponder(
				"POST",
				url,
				httpmock.NewStringResponder(tc.want.code, ""))
			metricPoller := NewCommonMetricPoller(ctx, client, nil, 2*time.Second)
			err = metricPoller.SendMetric(metric)
			if tc.want.msg != "" {
				assert.Equal(t, tc.want.msg, strings.TrimSpace(err.Error()))
			} else {
				assert.NoError(t, err)
			}
			callCountInfo := httpmock.GetCallCountInfo()
			count := callCountInfo[fmt.Sprintf("POST %s", url)]
			if tc.want.code != 0 {
				assert.Equal(t, 1, count)
			}
		})
	}
}
