package configs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

const (
	defaultMetricReportInterval = time.Second * 2
	defaultMetricPollInterval   = time.Second * 10
	defaultServerAddress        = "localhost:8080"
	defaultAgentAddress         = "localhost:8089"
	defaultHashKey              = ""
	defaultBuffBatchLimit       = 10
	defaultUseBuffSendClient    = true
	//defaultPublicKeyPath        = "cmd/agent/pubkey.pem"
)

type Duration struct {
	time.Duration
}

type AppConfig struct {
	ServerAddress       string
	AgentAddress        string
	ReportInterval      time.Duration
	PollInterval        time.Duration
	HashKey             string
	UseBuffSenderClient bool
	BuffBatchLimit      int
	PublicKey           *rsa.PublicKey
	publicKeyPath       string
}

func (cfg *AppConfig) UnmarshalJSON(data []byte) (err error) {
	cfgIn := struct {
		ServerAddress  string `json:"address"`
		ReportInterval string `json:"report_interval"`
		PollInterval   string `json:"poll_interval"`
		PublicKeyPath  string `json:"crypto_key"`
	}{}
	if err = json.Unmarshal(data, &cfgIn); err != nil {
		return err
	}
	cfg.ServerAddress = cfgIn.ServerAddress
	cfg.publicKeyPath = cfgIn.PublicKeyPath
	if cfgIn.ReportInterval != "" {
		if cfg.ReportInterval, err = time.ParseDuration(cfgIn.ReportInterval); err != nil {
			return err
		}
	}
	if cfgIn.PollInterval != "" {
		if cfg.PollInterval, err = time.ParseDuration(cfgIn.PollInterval); err != nil {
			return err
		}
	}
	return nil
}

func InitAppConfig(configPath string) (*AppConfig, error) {
	config, err := readConfig(configPath)
	if err != nil {
		return nil, err
	}
	setupConfigByFlags(config)
	setupConfigByEnvVars(config)
	setupConfigByDefaults(config)
	err = setupRSAKey(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func setupConfigByEnvVars(cfg *AppConfig) {
	if serverAddress := getServerAddress(); serverAddress != "" {
		cfg.ServerAddress = serverAddress
	}
	if reportInterval := getReportInterval(); reportInterval != 0 {
		cfg.ReportInterval = reportInterval
	}
	if pollInterval := getPollInterval(); pollInterval != 0 {
		cfg.PollInterval = pollInterval
	}
	if hashKey := getHashKey(); hashKey != "" {
		cfg.HashKey = hashKey
	}
	if agentAddress := getAgentAddress(); agentAddress != "" {
		cfg.AgentAddress = agentAddress
	}
	if cryptoKeyPath := getCryptoKeyPath(); cryptoKeyPath != "" {
		cfg.publicKeyPath = cryptoKeyPath
	}
	cfg.UseBuffSenderClient = useBuffSenderClient()
}

func setupConfigByFlags(cfg *AppConfig) {
	var serverAddressF string
	pflag.StringVarP(&serverAddressF, "a", "a", defaultServerAddress, "Address of the server")

	var reportIntervalF time.Duration
	pflag.DurationVarP(&reportIntervalF, "r", "r", defaultMetricReportInterval, "Report to server interval in sec")

	var pollIntervalF time.Duration
	pflag.DurationVarP(&pollIntervalF, "p", "p", defaultMetricPollInterval, "Pool metrics interval in sec")

	var hashKeyF string
	pflag.StringVarP(&hashKeyF, "k", "k", defaultHashKey, "Hash key")

	var cryptoKeyF string
	pflag.StringVarP(&cryptoKeyF, "crypto-key", "c", "", "Path to public key")

	pflag.Parse()

	if serverAddressF != "" {
		cfg.ServerAddress = serverAddressF
	}
	if reportIntervalF != 0 {
		cfg.ReportInterval = reportIntervalF
	}
	if pollIntervalF != 0 {
		cfg.PollInterval = pollIntervalF
	}
	if hashKeyF != "" {
		cfg.HashKey = hashKeyF
	}
	if cryptoKeyF != "" {
		cfg.publicKeyPath = cryptoKeyF
	}
}

func setupConfigByDefaults(cfg *AppConfig) {
	if cfg.BuffBatchLimit == 0 {
		cfg.BuffBatchLimit = defaultBuffBatchLimit
	}
	cfg.AgentAddress = defaultAgentAddress
}

func setupRSAKey(config *AppConfig) error {
	if config.publicKeyPath != "" {
		key, err := readRsaPublicKey(config.publicKeyPath)
		if err != nil {
			return err
		}
		config.PublicKey = key
	}
	return nil
}

func readConfig(configFilePath string) (*AppConfig, error) {
	if configFilePath == "" {
		return nil, errors.New("failed to init configuration: file path is not specified")
	}
	configBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configFile by '%s': %v", configFilePath, err)
	}
	var config AppConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config json '%s': %v", string(configBytes), err)
	}
	return &config, nil
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

func getCryptoKeyPath() string {
	return os.Getenv("CRYPTO_KEY")
}

func readRsaPublicKey(cryptoKeyPath string) (*rsa.PublicKey, error) {
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
