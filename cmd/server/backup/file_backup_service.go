package backup

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/server/configs"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"
)

type fileAutoBackup struct {
	interval   time.Duration
	repository repository.IMetricRepository
	storeFile  string
	*sync.RWMutex
}

// NewFileAutoBackup - инициализация fileAutoBackup, реализующего backup.IAutoBackup, для выгрузки метрик в текстоывй файл
func NewFileAutoBackup(repository repository.IMetricRepository, appConfig *configs.AppConfig) backup.IAutoBackup {
	err := initDir(appConfig.StoreFile)
	if err != nil {
		log.Fatalf("failed to create directories for '%s': %v", appConfig.StoreFile, err)
		return nil
	}
	return &fileAutoBackup{
		storeFile:  appConfig.StoreFile,
		interval:   appConfig.StoreInterval,
		repository: repository,
		RWMutex:    &sync.RWMutex{},
	}
}

func initDir(storeFile string) error {
	path := strings.Split(storeFile, "/")
	dir := strings.Join(path[0:len(path)-1], "/")
	return os.MkdirAll(dir, 0755)
}

func (b *fileAutoBackup) Start() chan struct{} {
	log.Println("FileAutoBackup activated")
	ticker := time.NewTicker(b.interval)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				log.Println("FileAutoBackup metrics interrupted!")
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("FileAutoBackup metrics start")
				err := b.Backup()
				if err != nil {
					log.Printf("failed to backup metrics: %v", err)
				}
			}
		}
	}()
	return done
}

func (b *fileAutoBackup) Restore() error {
	start := time.Now()
	metrics, err := b.readBackup()
	if err != nil {
		return err
	}
	for _, m := range metrics {
		err := b.repository.SaveMetric(m)
		if err != nil {
			return err
		}
	}
	log.Printf("[%v] Restore metrics finished, restored '%d' metrics", time.Since(start).String(), len(metrics))
	return nil
}

func (b *fileAutoBackup) Backup() error {
	start := time.Now()
	metrics, err := b.repository.GetMetricsList()
	if err != nil {
		return err
	}
	if len(metrics) == 0 {
		return nil
	}
	b.Lock()
	defer b.Unlock()
	file, err := os.OpenFile(b.storeFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&metrics)
	if err != nil {
		return err
	}
	log.Printf("[%v] FileAutoBackup metrics finished, saved '%d' metrics", time.Since(start).String(), len(metrics))
	return nil
}

func (b *fileAutoBackup) readBackup() ([]*model.Metric, error) {
	b.RLock()
	defer b.RUnlock()
	file, err := os.OpenFile(b.storeFile, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	metrics := make([]*model.Metric, 0)
	decoder := json.NewDecoder(file)
	if decoder.More() {
		err = decoder.Decode(&metrics)
		if err != nil {
			return nil, err
		}
	}

	return metrics, nil
}
