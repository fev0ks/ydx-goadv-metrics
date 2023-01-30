package configs

import (
	"github.com/spf13/pflag"
	"os"
	"strconv"
	"time"
)

const (
	defaultMetricReportInterval = time.Second * 2
	defaultMetricPollInterval   = time.Second * 10
	defaultServerAddress        = "localhost:8080"
	defaultHashKey              = ""
	defaultBuffBatchLimit       = 10
	defaultUseBuffSendClient    = true
)

type AppConfig struct {
	ServerAddress       string
	ReportInterval      time.Duration
	PollInterval        time.Duration
	HashKey             string
	UseBuffSenderClient bool
	BuffBatchLimit      int
}

func InitAppConfig() *AppConfig {
	address := getServerAddress()
	var addressF string
	pflag.StringVarP(&addressF, "a", "a", defaultServerAddress, "Address of the server")

	reportInterval := getReportInterval()
	var reportIntervalF time.Duration
	pflag.DurationVarP(&reportIntervalF, "r", "r", defaultMetricReportInterval, "Report to server interval in sec")

	pollInterval := getPollInterval()
	var pollIntervalF time.Duration
	pflag.DurationVarP(&pollIntervalF, "p", "p", defaultMetricPollInterval, "Pool metrics interval in sec")

	hashKey := getHashKey()
	var hashKeyF string
	pflag.StringVarP(&hashKeyF, "k", "k", defaultHashKey, "Hash key")

	pflag.Parse()

	if address == "" {
		address = addressF
	}
	if reportInterval == 0 {
		reportInterval = reportIntervalF
	}
	if pollInterval == 0 {
		pollInterval = pollIntervalF
	}
	if hashKey == "" {
		hashKey = hashKeyF
	}
	return &AppConfig{
		ServerAddress:       address,
		ReportInterval:      reportInterval,
		PollInterval:        pollInterval,
		HashKey:             hashKey,
		UseBuffSenderClient: useBuffSenderClient(),
		BuffBatchLimit:      defaultBuffBatchLimit,
	}
}

func getReportInterval() time.Duration {
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

func getPollInterval() time.Duration {
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

func getServerAddress() string {
	return os.Getenv("ADDRESS")
}

func getHashKey() string {
	return os.Getenv("KEY")
}

func useBuffSenderClient() bool {
	useBuffSendClient, err := strconv.ParseBool(os.Getenv("USE_BUFF_SEND_CLIENT"))
	if err != nil {
		useBuffSendClient = defaultUseBuffSendClient
	}
	return useBuffSendClient
}
