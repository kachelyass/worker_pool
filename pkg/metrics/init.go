package metrics

import "github.com/prometheus/client_golang/prometheus"

func Init() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDuration)
	prometheus.MustRegister(HTTPInFlight)

	prometheus.MustRegister(KafkaProduceTotal)
	prometheus.MustRegister(KafkaProduceDuration)
	prometheus.MustRegister(KafkaProduceErrorsTotal)

	prometheus.MustRegister(KafkaConsumeTotal)
	prometheus.MustRegister(KafkaConsumeErrorsTotal)
	prometheus.MustRegister(KafkaConsumeDuration)
	prometheus.MustRegister(KafkaConsumerRebalancesTotal)

	prometheus.MustRegister(DBQueryTotal)
	prometheus.MustRegister(DBQueryDuration)
	prometheus.MustRegister(DBQueryErrorsTotal)
	prometheus.MustRegister(DBPoolAcquiredConnections)
	prometheus.MustRegister(DBPoolIdleConnections)
	prometheus.MustRegister(DBPoolTotalConnections)

	prometheus.MustRegister(WorkerPoolWorkers)
	prometheus.MustRegister(WorkerPoolQueueSize)
	prometheus.MustRegister(WorkerPoolJobsEnqueuedTotal)
	prometheus.MustRegister(WorkerPoolJobsStartedTotal)
	prometheus.MustRegister(WorkerPoolJobsCompletedTotal)
	prometheus.MustRegister(WorkerPoolJobsFailedTotal)
	prometheus.MustRegister(WorkerPoolJobDuration)
}
