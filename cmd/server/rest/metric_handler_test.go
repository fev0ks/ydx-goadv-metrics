package rest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/repositories"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceptionMetricsHandler(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	testCases := []struct {
		name        string
		requestPath string ///update/{mType}/{name}/{value}
		want        want
	}{
		{
			name:        "Should catch Counter metric",
			requestPath: "/update/counter/case1/1",
			want: want{
				http.StatusOK,
				"",
			},
		},
		{
			name:        "Should catch Gauge metric",
			requestPath: "/update/gauge/case1/1.1",
			want: want{
				http.StatusOK,
				"",
			},
		},
		{
			name:        "Should return 400 if empty name",
			requestPath: "/update/gauge//1.1",
			want: want{
				http.StatusBadRequest,
				"metric 'name' must be specified",
			},
		},
		{
			name:        "Should return 400 if empty mType",
			requestPath: "/update//case1/1.1",
			want: want{
				http.StatusBadRequest,
				"metric 'mType' must be specified",
			},
		},
		{
			name:        "Should return 404 if empty value",
			requestPath: "/update/gauge/case1/",
			want: want{
				http.StatusNotFound,
				"404 page not found",
			},
		},
		{
			name:        "Should return 501 if not supported mType",
			requestPath: "/update/fake/case1/123",
			want: want{
				http.StatusNotImplemented,
				"type 'fake' is not supported",
			},
		},
		{
			name:        "Should return 400 if invalid counter value",
			requestPath: "/update/counter/case1/123.123",
			want: want{
				http.StatusBadRequest,
				"strconv.ParseUint: parsing \"123.123\": invalid syntax",
			},
		},
		{
			name:        "Should return 400 if invalid gauge value",
			requestPath: "/update/gauge/case1/xxx.yyy",
			want: want{
				http.StatusBadRequest,
				"strconv.ParseFloat: parsing \"xxx.yyy\": invalid syntax",
			},
		},
		{
			name:        "Should return 400 if invalid gauge value",
			requestPath: "/update/gauge/case1/xxx.yyy",
			want: want{
				http.StatusBadRequest,
				"strconv.ParseFloat: parsing \"xxx.yyy\": invalid syntax",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cr := repositories.NewCommonRepository()
			mh := NewMetricsHandler(ctx, cr)
			r := NewRouter()
			HandleMetricRequests(r, mh)

			ts := httptest.NewServer(r)
			defer ts.Close()

			actualStatusCode, actualResponse := sendRequest(t, ts, http.MethodPost, tc.requestPath)
			assert.Equal(t, tc.want.statusCode, actualStatusCode)
			assert.Equal(t, tc.want.response, strings.TrimSpace(actualResponse))
		})
	}
}

func TestGetMetricsHandler(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	testCases := []struct {
		name    string
		metrics []*model.Metric
		want    want
	}{
		{
			name: "Should return html page with metrics",
			metrics: []*model.Metric{
				model.NewCounterMetric("counter1", 123),
				model.NewGaugeMetric("gauge1", 123.123),
			},
			want: want{
				http.StatusOK,
				"<html><style>\ntable, th, td {\n  border:1px solid black;\n}\n</style>\n<body><body><h2>Metrics result</h2><table><tr><th>#</th><th>type</th><th>name</th><th>value</th></tr><tr><td>1</td><td>counter</td><td>counter1</td><td>123</td></tr><tr><td>2</td><td>gauge</td><td>gauge1</td><td>123.123000</td></tr></table></body><html>",
			},
		},
		{
			name:    "Should return html page without metrics",
			metrics: []*model.Metric{},
			want: want{
				http.StatusOK,
				"<html><style>\ntable, th, td {\n  border:1px solid black;\n}\n</style>\n<body><body><h2>Metrics result</h2><table><tr><th>#</th><th>type</th><th>name</th><th>value</th></tr></table></body><html>",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cr := repositories.NewCommonRepository()
			defer cr.Clear()
			for _, metric := range tc.metrics {
				_ = cr.SaveMetric(metric)
			}
			mh := NewMetricsHandler(ctx, cr)
			r := NewRouter()
			HandleMetricRequests(r, mh)

			ts := httptest.NewServer(r)
			defer ts.Close()

			actualStatusCode, actualResponse := sendRequest(t, ts, http.MethodGet, "")
			assert.Equal(t, tc.want.statusCode, actualStatusCode)
			assert.Equal(t, tc.want.response, strings.TrimSpace(actualResponse))
		})
	}
}

func TestGetMetricHandler(t *testing.T) {
	type want struct {
		statusCode int
		response   string
	}
	testCases := []struct {
		name        string
		requestPath string // value/{mType}/{name}
		metrics     []*model.Metric
		want        want
	}{
		{
			name:        "Should return gauge metric value",
			requestPath: "/value/gauge/gauge1",
			metrics: []*model.Metric{
				model.NewCounterMetric("counter1", 1),
				model.NewGaugeMetric("gauge1", 1.1),
			},
			want: want{
				http.StatusOK,
				"1.1",
			},
		},
		{
			name:        "Should return total counter value",
			requestPath: "/value/counter/counter1",
			metrics: []*model.Metric{
				model.NewCounterMetric("counter1", 1),
				model.NewCounterMetric("counter1", 2),
			},
			want: want{
				http.StatusOK,
				"3",
			},
		},
		{
			name:        "Should return 404 when metric was not found",
			requestPath: "/value/counter/counter123",
			metrics: []*model.Metric{
				model.NewCounterMetric("counter1", 1),
				model.NewGaugeMetric("gauge1", 1.1),
			},
			want: want{
				http.StatusNotFound,
				"metric was not found: counter123",
			},
		},
		{
			name:        "Should return 404 when metric name is empty",
			requestPath: "/value/counter/",
			metrics: []*model.Metric{
				model.NewCounterMetric("counter1", 1),
				model.NewGaugeMetric("gauge1", 1.1),
			},
			want: want{
				http.StatusNotFound,
				"404 page not found",
			},
		},
		{
			name:        "Should return 400 when metric mType is empty",
			requestPath: "/value//counter1",
			metrics: []*model.Metric{
				model.NewCounterMetric("counter1", 1),
				model.NewGaugeMetric("gauge1", 1.1),
			},
			want: want{
				http.StatusBadRequest,
				"metric 'mType' must be specified",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cr := repositories.NewCommonRepository()
			defer cr.Clear()
			for _, metric := range tc.metrics {
				_ = cr.SaveMetric(metric)
			}
			mh := NewMetricsHandler(ctx, cr)
			r := NewRouter()
			HandleMetricRequests(r, mh)

			ts := httptest.NewServer(r)
			defer ts.Close()

			actualStatusCode, actualResponse := sendRequest(t, ts, http.MethodGet, tc.requestPath)
			assert.Equal(t, tc.want.statusCode, actualStatusCode)
			assert.Equal(t, tc.want.response, strings.TrimSpace(actualResponse))
		})
	}
}

func sendRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}
