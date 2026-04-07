package postgre

import (
	"context"
	"fmt"
	"time"
	"worker_pool/pkg/metrics/postgresmetrics"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	cfg.MaxConns = 20
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 10 * time.Minute
	cfg.MaxConnIdleTime = 10 * time.Minute
	cfg.HealthCheckPeriod = 30 * time.Second

	cfg.ConnConfig.Tracer = postgresmetrics.NewMetricsTracer()
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stat := pool.Stat()
				postgresmetrics.DBPoolAcquiredConnections.Set(float64(stat.AcquiredConns()))
				postgresmetrics.DBPoolIdleConnections.Set(float64(stat.IdleConns()))
				postgresmetrics.DBPoolTotalConnections.Set(float64(stat.TotalConns()))
			}
		}
	}()

	return pool, nil
}
