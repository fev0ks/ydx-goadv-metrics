package configs

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultMetricReportInterval = time.Second * 2
	defaultMetricPollInterval   = time.Second * 10

	defaultServerAddress = "localhost:8080"
)

func GetReportInterval() time.Duration {
	reportInterval := os.Getenv("REPORT_INTERVAL")
	if reportInterval == "" {
		return defaultMetricReportInterval
	}
	reportIntervalVal, err := strconv.Atoi(reportInterval)
	if err != nil {
		return defaultMetricReportInterval
	}
	return time.Second * time.Duration(reportIntervalVal)
}

func GetPollInterval() time.Duration {
	reportInterval := os.Getenv("POLL_INTERVAL")
	if reportInterval == "" {
		return defaultMetricPollInterval
	}
	reportIntervalVal, err := strconv.Atoi(reportInterval)
	if err != nil {
		return defaultMetricPollInterval
	}
	return time.Second * time.Duration(reportIntervalVal)
}

func GetServerAddress() string {
	host := os.Getenv("ADDRESS")
	if host == "" {
		return defaultServerAddress
	}
	return host
}
