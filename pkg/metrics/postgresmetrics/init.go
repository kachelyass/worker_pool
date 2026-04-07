package postgresmetrics

import "github.com/prometheus/client_golang/prometheus"

func Init() {
	prometheus.MustRegister(DBQueryTotal)
	prometheus.MustRegister(DBQueryErrorsTotal)
	prometheus.MustRegister(DBQueryDuration)
	prometheus.MustRegister(DBPoolAcquiredConnections)
	prometheus.MustRegister(DBPoolIdleConnections)
	prometheus.MustRegister(DBPoolTotalConnections)
}
