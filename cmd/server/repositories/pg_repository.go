package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/fev0ks/ydx-goadv-metrics/internal/model"
	"github.com/fev0ks/ydx-goadv-metrics/internal/model/server"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	migrationsDir = "migrations/postgres"
	dbDialect     = "postgres"
)

var (
	saveCounterMetric = "saveCounterMetricStmt"
	saveGaugeMetric   = "saveGaugeMetricStmt"
	getMetric         = "getMetricStmt"
	getMetrics        = "getMetricsStmt"
	truncate          = "truncateStmt"
)

type pgRepository struct {
	db         *sql.DB
	ctx        context.Context
	statements map[string]*sql.Stmt
}

func NewPgRepository(dbConfig string, ctx context.Context) (server.IMetricRepository, error) {
	if dbConfig == "" {
		log.Println("Postgres DB config is empty")
		return nil, errors.New("failed to init pg repository: config is empty")
	}
	log.Printf("Trying to connect: %s", dbConfig)
	db, err := sql.Open(dbDialect, dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to Postgres DB: %v", err)
	}

	pgRep := &pgRepository{db: db, ctx: ctx}
	pgRep.migrationUp()
	err = pgRep.prepareStatements()
	if err != nil {
		return nil, fmt.Errorf("failed to prepareStatements for Postgres DB: %v", err)
	}
	return pgRep, nil
}

func (p *pgRepository) prepareStatements() error {
	log.Println("Prepare statements")
	p.statements = make(map[string]*sql.Stmt)
	if saveCounterMetricStmt, err := p.db.Prepare(
		"insert into metrics(id, type, delta, hash) values($1, $2, $3, $4) " +
			"on conflict (id) do update set type = excluded.type, delta = excluded.delta, hash = excluded.hash",
	); err != nil {
		return err
	} else {
		p.statements[saveCounterMetric] = saveCounterMetricStmt
	}
	if saveGaugeMetricStmt, err := p.db.Prepare(
		"insert into metrics(id, type, value, hash) values($1, $2, $3, $4) " +
			"on conflict (id) do update set type = excluded.type, value = excluded.value, hash = excluded.hash",
	); err != nil {
		return err
	} else {
		p.statements[saveGaugeMetric] = saveGaugeMetricStmt
	}
	if getMetricStmt, err := p.db.Prepare("select * from metrics where id = $1"); err != nil {
		return err
	} else {
		p.statements[getMetric] = getMetricStmt
	}
	if getMetricsStmt, err := p.db.Prepare("select * from metrics"); err != nil {
		return err
	} else {
		p.statements[getMetrics] = getMetricsStmt
	}
	if truncateStmt, err := p.db.Prepare("truncate metrics"); err != nil {
		return err
	} else {
		p.statements[truncate] = truncateStmt
	}

	return nil
}

func (p *pgRepository) HealthCheck(ctx context.Context) error {
	if err := p.db.PingContext(ctx); err != nil {
		log.Printf("failed to check connection to Postgres DB: %v", err)
		return err
	}
	log.Println("Postgres DB connection is active")
	return nil
}

func (p *pgRepository) migrationUp() {
	if p.db == nil {
		log.Fatalln(errors.New("failed to start db migration: db instance is nil"))
	}
	log.Println("migrations are started")
	migration := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	countOfMigrations, err := migrate.Exec(p.db, dbDialect, migration, migrate.Up)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("migrations are finished, total count: %d", countOfMigrations)
}

func (p *pgRepository) SaveMetrics(metrics []*model.Metric) error {
	for _, metric := range metrics {
		err := p.SaveMetric(metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *pgRepository) SaveMetric(metric *model.Metric) error {
	var err error
	switch metric.MType {
	case model.GaugeType:
		_, err = p.statements[saveGaugeMetric].Exec(metric.ID, metric.MType, metric.Value, metric.Hash)
	case model.CounterType:
		var currentMetric *model.Metric
		currentMetric, err = p.GetMetric(metric.ID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if currentMetric != nil {
			newValue := *currentMetric.Delta + *metric.Delta
			metric.Delta = &newValue
		}
		_, err = p.statements[saveCounterMetric].Exec(metric.ID, metric.MType, metric.Delta, metric.Hash)
	}
	return err
}

func (p *pgRepository) GetMetrics() (map[string]*model.Metric, error) {
	metricsList, err := p.GetMetricsList()
	if err != nil {
		return nil, err
	}
	metrics := make(map[string]*model.Metric, len(metricsList))
	for _, metric := range metricsList {
		metrics[metric.ID] = metric
	}
	return metrics, nil
}

func (p *pgRepository) GetMetricsList() ([]*model.Metric, error) {
	metrics := make([]*model.Metric, 0)
	rows, err := p.statements[getMetrics].Query()
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	for rows.Next() {
		metric := &model.Metric{}
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value, &metric.Hash)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (p *pgRepository) GetMetric(name string) (*model.Metric, error) {
	row := p.statements[getMetric].QueryRow(name)
	metric := &model.Metric{}
	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value, &metric.Hash)
	if err != nil {
		return nil, err
	}
	return metric, nil
}

func (p *pgRepository) Clear() error {
	exec, err := p.statements[truncate].Exec()
	if err != nil {
		return err
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("removed %d rows", affected)
	return nil
}

func (p *pgRepository) Close() error {
	for _, stmt := range p.statements {
		err := stmt.Close()
		if err != nil {
			return err
		}
	}
	err := p.db.Close()
	return err
}
