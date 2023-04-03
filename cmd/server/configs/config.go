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
	defaultAddress             = "localhost:8080"
	defaultMetricStoreInterval = 300 * time.Second
	defaultStoreFile           = "/tmp/devops-metrics-db.json"
	defaultDoRestore           = true
	defaultHashKey             = ""
	defaultDBConfig            = ""
	defaultPrivateKeyPath      = "cmd/server/privkey.pem"
)

type AppConfig struct {
	// ServerAddress - Адресс сервиса
	ServerAddress string
	// StoreInterval - Временной интервал для беккапа метрик
	StoreInterval time.Duration
	// DoRestore - Восстанавливать ли метрики в память из беккапа при старте сервиса
	DoRestore bool
	// StoreFile - Имя файла при беккапе метрик в файл
	StoreFile string
	// HashKey - Любое текстовое значение,
	// обязательно должно совпадать с аналогичным параметров в Агент сервисе для архивации//разархивации сообщений
	HashKey string
	// DBConfig - Конфиг подключения к базе
	DBConfig   string
	PrivateKey *rsa.PrivateKey
}

func InitAppConfig() *AppConfig {
	address := getAddress()
	var addressF string
	pflag.StringVarP(&addressF, "a", "a", defaultAddress, "Address of the server")

	restore := getDoReStore()
	var restoreF bool
	pflag.BoolVarP(&restoreF, "r", "r", defaultDoRestore, "Do autoBackup restore?")

	storeInterval := getStoreInterval()
	var storeIntervalF time.Duration
	pflag.DurationVarP(&storeIntervalF, "i", "i", defaultMetricStoreInterval, "Backup interval in sec")

	storeFile := getStoreFile()
	var storeFileF string
	pflag.StringVarP(&storeFileF, "f", "f", defaultStoreFile, "Path of Backup store file")

	hashKey := getHashKey()
	var hashKeyF string
	pflag.StringVarP(&hashKeyF, "k", "k", defaultHashKey, "Hash key")

	dbConfig := getDBConfig()
	var dbDsnF string
	pflag.StringVarP(&dbDsnF, "d", "d", defaultDBConfig, "Postgres DB DSN")

	pflag.Parse()

	if address == "" {
		address = addressF
	}
	if restore == nil {
		restore = &restoreF
	}
	if storeInterval == 0 {
		storeInterval = storeIntervalF
	}
	if address == "" {
		address = addressF
	}
	if restore == nil {
		restore = &restoreF
	}
	if storeFile == "" {
		storeFile = storeFileF
	}
	if hashKey == "" {
		hashKey = hashKeyF
	}
	if dbConfig == "" {
		dbConfig = dbDsnF
	}
	privateKey, err := readRsaPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	return &AppConfig{
		ServerAddress: address,
		StoreInterval: storeInterval,
		DoRestore:     *restore,
		StoreFile:     storeFile,
		HashKey:       hashKey,
		DBConfig:      dbConfig,
		PrivateKey:    privateKey,
	}
}

func getAddress() string {
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

func readRsaPrivateKey() (*rsa.PrivateKey, error) {
	cryptoKeyPath := os.Getenv("CRYPTO_KEY")
	if cryptoKeyPath == "" {
		cryptoKeyPath = defaultPrivateKeyPath
	}
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
