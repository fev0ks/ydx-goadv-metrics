package repositories

import (
	"context"
	"database/sql"
	"errors"
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

type pgRepository struct {
	db *sql.DB
}

func InitPgRepository(dbConfig string) server.MetricRepository {
	if dbConfig == "" {
		log.Println("Postgres DB config is empty")
		return nil
	}
	log.Printf("Trying to connect: %s\n", dbConfig)
	db, err := sql.Open(dbDialect, dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to Postgres DB: %v", err)
	}

	pgRep := &pgRepository{db}
	//pgRep.migrationUp()
	return pgRep
}

func (p *pgRepository) HealthCheck(ctx context.Context) error {
	if err := p.db.PingContext(ctx); err != nil {
		log.Printf("failed to check connection to Postgres DB: %v\n", err)
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

func (p *pgRepository) SaveMetric(metric *model.Metric) error {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepository) GetMetrics() map[string]*model.Metric {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepository) GetMetricsList() []*model.Metric {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepository) GetMetric(name string) *model.Metric {
	//TODO implement me
	panic("implement me")
}

func (p *pgRepository) Clear() {
	//TODO implement me
	panic("implement me")
}
