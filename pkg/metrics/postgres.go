package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	DBQueryTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_total",
			Help: "Total number of database queries.",
		},
		[]string{"operation"},
	)

	DBQueryErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_errors_total",
			Help: "Total number of database query errors.",
		},
		[]string{"operation"},
	)

	DBQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)

	DBPoolAcquiredConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_acquired_connections",
			Help: "Number of currently acquired connections in the DB pool.",
		},
	)

	DBPoolIdleConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_idle_connections",
			Help: "Number of idle connections in the DB pool.",
		},
	)

	DBPoolTotalConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_total_connections",
			Help: "Total number of connections in the DB pool.",
		},
	)
)
