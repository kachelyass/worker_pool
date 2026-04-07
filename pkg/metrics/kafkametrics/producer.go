package kafkametrics

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaProduceTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_produce_total",
			Help: "Total number of Kafka produce attempts.",
		},
		[]string{"topic"},
	)

	KafkaProduceErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_produce_errors_total",
			Help: "Total number of Kafka produce errors.",
		},
		[]string{"topic"},
	)

	KafkaProduceDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_produce_duration_seconds",
			Help:    "Kafka produce duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"})

	KafkaConsumerRebalancesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_consumer_rebalances_total",
			Help: "Total number of Kafka consumer group rebalances.",
		},
	)
)
