package backup

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"
)

type AutoBackup interface {
	Start() chan struct{}
	Restore() error
}

type autoBackup struct {
	storeFile  string
	interval   time.Duration
	repository server.MetricRepository
	*sync.RWMutex
}

func NewAutoBackup(storeFile string, interval time.Duration, repository server.MetricRepository) *autoBackup {
	err := initDir(storeFile)
	if err != nil {
		log.Fatalf("failed to create directories for '%s': %v", storeFile, err)
		return nil
	}
	return &autoBackup{
		storeFile:  storeFile,
		interval:   interval,
		repository: repository,
		RWMutex:    &sync.RWMutex{},
	}
}

func initDir(storeFile string) error {
	path := strings.Split(storeFile, "/")
	dir := strings.Join(path[0:len(path)-1], "/")
	return os.MkdirAll(dir, 0755)
}

func (b *autoBackup) Start() chan struct{} {
	log.Println("AutoBackup activated")
	ticker := time.NewTicker(b.interval)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				log.Println("AutoBackup metrics interrupted!")
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("AutoBackup metrics start")
				err := b.Backup()
				if err != nil {
					log.Printf("failed to backup metrics: %v\n", err)
				}
			}
		}
	}()
	return done
}

func (b *autoBackup) Restore() error {
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
	log.Printf("[%v] Restore metrics finished, restored '%d' metrics\n", time.Since(start).String(), len(metrics))
	return nil
}

func (b *autoBackup) Backup() error {
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
	log.Printf("[%v] AutoBackup metrics finished, saved '%d' metrics\n", time.Since(start).String(), len(metrics))
	return nil
}

func (b *autoBackup) readBackup() ([]*model.Metric, error) {
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
