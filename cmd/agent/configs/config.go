package configs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

const (
	defaultMetricReportInterval = time.Second * 2
	defaultMetricPollInterval   = time.Second * 10
	defaultServerAddress        = "localhost:8080"
	defaultAgentAddress         = "localhost:8085"
	defaultHashKey              = ""
	defaultBuffBatchLimit       = 10
	defaultUseBuffSendClient    = true
	defaultPublicKeyPath        = "cmd/agent/pubkey.pem"
)

type AppConfig struct {
	ServerAddress       string
	AgentAddress        string
	ReportInterval      time.Duration
	PollInterval        time.Duration
	HashKey             string
	UseBuffSenderClient bool
	BuffBatchLimit      int
	PublicKey           *rsa.PublicKey
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

	agentAddress := getAgentAddress()
	if agentAddress == "" {
		agentAddress = defaultAgentAddress
	}
	publicKey, err := readRsaPublicKey()
	if err != nil {
		log.Fatal(err)
	}
	return &AppConfig{
		ServerAddress:       address,
		AgentAddress:        agentAddress,
		ReportInterval:      reportInterval,
		PollInterval:        pollInterval,
		HashKey:             hashKey,
		UseBuffSenderClient: useBuffSenderClient(),
		BuffBatchLimit:      defaultBuffBatchLimit,
		PublicKey:           publicKey,
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

func getAgentAddress() string {
	return os.Getenv("AGENT_ADDRESS")
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

func readRsaPublicKey() (*rsa.PublicKey, error) {
	cryptoKeyPath := os.Getenv("CRYPTO_KEY")
	if cryptoKeyPath == "" {
		cryptoKeyPath = defaultPublicKeyPath
	}
	pemBytes, err := os.ReadFile(cryptoKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read publicKey by '%s': %v", cryptoKeyPath, err)
	}
	block, _ := pem.Decode(pemBytes)
	key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse publicKey: %v", err)
	}
	return key, nil
}
