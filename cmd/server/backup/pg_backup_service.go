package backup

import (
	"log"
	"sync"
	"time"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/backup"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server/repository"
)

type pgAutoBackup struct {
	interval   time.Duration
	repository repository.IMetricRepository
	storage    repository.IMetricRepository
	*sync.RWMutex
}

// NewPgAutoBackup - инициализация pgAutoBackup, реализующего backup.IAutoBackup, для выгрузки метрик в Postgres базу данных
func NewPgAutoBackup(interval time.Duration, repository repository.IMetricRepository, storage repository.IMetricRepository) backup.IAutoBackup {
	return &pgAutoBackup{
		interval:   interval,
		repository: repository,
		storage:    storage,
		RWMutex:    &sync.RWMutex{},
	}
}

func (b *pgAutoBackup) Start() chan struct{} {
	log.Println("PgAutoBackup activated")
	ticker := time.NewTicker(b.interval)
	done := make(chan struct{})
	go func(one chan struct{}, ticker *time.Ticker) {
		for {
			b.doBackupIfNotCanceled(done, ticker)
		}
	}(done, ticker)
	return done
}

func (b *pgAutoBackup) doBackupIfNotCanceled(done chan struct{}, ticker *time.Ticker) {
	select {
	case <-done:
		log.Println("PgAutoBackup metrics interrupted!")
		ticker.Stop()
		return
	case <-ticker.C:
		log.Println("PgAutoBackup metrics start")
		err := b.Backup()
		if err != nil {
			log.Printf("failed to backup metrics: %v", err)
		}
	}
}

func (b *pgAutoBackup) Restore() error {
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

func (b *pgAutoBackup) Backup() error {
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
	err = b.storage.SaveMetrics(metrics)
	if err != nil {
		return err
	}
	log.Printf("[%v] PgAutoBackup metrics finished, saved '%d' metrics", time.Since(start).String(), len(metrics))
	return nil
}

func (b *pgAutoBackup) readBackup() ([]*model.Metric, error) {
	b.RLock()
	defer b.RUnlock()
	return b.storage.GetMetricsList()
}
