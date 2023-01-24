package configs

import (
	"os"
	"strconv"
	"time"
)

const (
	DefaultMetricReportInterval = time.Second * 2
	DefaultMetricPollInterval   = time.Second * 10
	DefaultServerAddress        = "localhost:8080"
	DefaultHashKey              = ""
	defaultBuffSendInterval     = "1s"
	defaultBuffBatchLimit       = 10
	defaultUseBuffSendClient    = true
)

func GetReportInterval() time.Duration {
	reportInterval := os.Getenv("REPORT_INTERVAL")
	if reportInterval == "" {
		return 0
	}
	reportIntervalVal, err := strconv.Atoi(reportInterval)
	if err != nil {
		return 0
	}
	duration := time.Second * time.Duration(reportIntervalVal)
	return duration
}

func GetPollInterval() time.Duration {
	reportInterval := os.Getenv("POLL_INTERVAL")
	if reportInterval == "" {
		return 0
	}
	reportIntervalVal, err := strconv.Atoi(reportInterval)
	if err != nil {
		return 0
	}
	duration := time.Second * time.Duration(reportIntervalVal)
	return duration
}

func GetServerAddress() string {
	return os.Getenv("ADDRESS")
}

func GetHashKey() string {
	return os.Getenv("KEY")
}

func UseBuffSenderClient() bool {
	useBuffSendClient, err := strconv.ParseBool(os.Getenv("USE_BUFF_SEND_CLIENT"))
	if err != nil {
		useBuffSendClient = defaultUseBuffSendClient
	}
	return useBuffSendClient
}

func GetBuffBatchLimit() int {
	buffBatchLimit := os.Getenv("BUFF_BATCH_LIMIT")
	if buffBatchLimit == "" {
		return defaultBuffBatchLimit
	}
	buffBatchLimitVal, err := strconv.Atoi(buffBatchLimit)
	if err != nil {
		return defaultBuffBatchLimit
	}
	return buffBatchLimitVal
}

func GetBuffSendInterval() time.Duration {
	buffSendInterval := os.Getenv("BUFF_SEND_INTERVAL")

	if buffSendInterval == "" {
		buffSendInterval = defaultBuffSendInterval
	}

	interval, err := time.ParseDuration(buffSendInterval)
	if err != nil {
		interval, _ = time.ParseDuration(defaultBuffSendInterval)
	}
	return interval
}
