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
	defaultAddress             = "localhost:8080"
	defaultMetricStoreInterval = 300 * time.Second
	defaultStoreFile           = "/tmp/devops-metrics-db.json"
	defaultDoRestore           = true
	defaultHashKey             = ""
	defaultDBConfig            = ""
	//defaultPrivateKeyPath      = "cmd/server/privkey.pem"
)

type AppConfig struct {
	// ServerAddress - Адресс сервиса
	ServerAddress string
	// StoreInterval - Временной интервал для беккапа метрик
	StoreInterval time.Duration
	// DoRestore - Восстанавливать ли метрики в память из беккапа при старте сервиса
	DoRestore *bool
	// StoreFile - Имя файла при беккапе метрик в файл
	StoreFile string
	// HashKey - Любое текстовое значение,
	// обязательно должно совпадать с аналогичным параметров в Агент сервисе для архивации//разархивации сообщений
	HashKey string
	// DBConfig - Конфиг подключения к базе
	DBConfig       string
	PrivateKey     *rsa.PrivateKey
	privateKeyPath string
}

func (cfg *AppConfig) UnmarshalJSON(data []byte) (err error) {
	cfgIn := struct {
		ServerAddress  string `json:"address"`
		DoRestore      *bool  `json:"restore"`
		StoreInterval  string `json:"store_interval"`
		StoreFile      string `json:"store_file"`
		DBConfig       string `json:"database_dsn"`
		PrivateKeyPath string `json:"crypto_key"`
	}{}
	if err = json.Unmarshal(data, &cfgIn); err != nil {
		return err
	}
	cfg.ServerAddress = cfgIn.ServerAddress
	cfg.DoRestore = cfgIn.DoRestore
	cfg.StoreFile = cfgIn.StoreFile
	cfg.DBConfig = cfgIn.DBConfig
	cfg.privateKeyPath = cfgIn.PrivateKeyPath
	if cfgIn.StoreInterval != "" {
		if cfg.StoreInterval, err = time.ParseDuration(cfgIn.StoreInterval); err != nil {
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
	if doRestore := getDoReStore(); doRestore != nil {
		cfg.DoRestore = doRestore
	}
	if storeInterval := getStoreInterval(); storeInterval != 0 {
		cfg.StoreInterval = storeInterval
	}
	if hashKey := getHashKey(); hashKey != "" {
		cfg.HashKey = hashKey
	}
	if storeFile := getStoreFile(); storeFile != "" {
		cfg.StoreFile = storeFile
	}
	if dbConfig := getDBConfig(); dbConfig != "" {
		cfg.DBConfig = dbConfig
	}
	if cryptoKeyPath := getCryptoKeyPath(); cryptoKeyPath != "" {
		cfg.privateKeyPath = cryptoKeyPath
	}
}

func setupConfigByFlags(cfg *AppConfig) {
	var serverAddressF string
	pflag.StringVarP(&serverAddressF, "a", "a", defaultAddress, "Address of the server")

	var restoreF bool
	pflag.BoolVarP(&restoreF, "r", "r", defaultDoRestore, "Do autoBackup restore?")

	var storeIntervalF time.Duration
	pflag.DurationVarP(&storeIntervalF, "i", "i", defaultMetricStoreInterval, "Backup interval in sec")

	var storeFileF string
	pflag.StringVarP(&storeFileF, "f", "f", defaultStoreFile, "Path of Backup store file")

	var hashKeyF string
	pflag.StringVarP(&hashKeyF, "k", "k", defaultHashKey, "Hash key")

	var dbDsnF string
	pflag.StringVarP(&dbDsnF, "d", "d", defaultDBConfig, "Postgres DB DSN")

	var cryptoKeyF string
	pflag.StringVarP(&cryptoKeyF, "crypto-key", "c", "", "Path to private key")

	pflag.Parse()

	if serverAddressF != "" {
		cfg.ServerAddress = serverAddressF
	}
	if cfg.DoRestore == nil {
		cfg.DoRestore = &restoreF
	}
	if storeIntervalF != 0 {
		cfg.StoreInterval = storeIntervalF
	}
	if storeFileF != "" {
		cfg.StoreFile = storeFileF
	}
	if hashKeyF != "" {
		cfg.HashKey = hashKeyF
	}
	if dbDsnF != "" {
		cfg.DBConfig = dbDsnF
	}
	if cryptoKeyF != "" {
		cfg.privateKeyPath = cryptoKeyF
	}
}

func setupRSAKey(config *AppConfig) error {
	if config.privateKeyPath != "" {
		key, err := readRsaPrivateKey(config.privateKeyPath)
		if err != nil {
			return err
		}
		config.PrivateKey = key
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

func getServerAddress() string {
	return os.Getenv("ADDRESS")
}

func getStoreInterval() time.Duration {
	storeInterval := os.Getenv("STORE_INTERVAL")
	if storeInterval == "" {
		return 0
	}
	storeIntervalVal, err := strconv.Atoi(storeInterval)
	if err != nil {
		return 0
	}
	duration := time.Duration(storeIntervalVal) * time.Second
	return duration
}

func getStoreFile() string {
	return os.Getenv("STORE_FILE")
}

func getHashKey() string {
	return os.Getenv("KEY")
}

func getDoReStore() *bool {
	doReStore := os.Getenv("RESTORE")
	if doReStore == "" {
		return nil
	}
	doReStoreVal, err := strconv.ParseBool(doReStore)
	if err != nil {
		return nil
	}
	return &doReStoreVal
}

func getDBConfig() string {
	return os.Getenv("DATABASE_DSN")
}

func getCryptoKeyPath() string {
	return os.Getenv("CRYPTO_KEY")
}

func readRsaPrivateKey(cryptoKeyPath string) (*rsa.PrivateKey, error) {
	pemBytes, err := os.ReadFile(cryptoKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read publicKey by '%s': %v", cryptoKeyPath, err)
	}
	block, _ := pem.Decode(pemBytes)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse publicKey: %v", err)
	}
	return key, nil
}
