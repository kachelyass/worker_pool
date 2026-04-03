package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaConsumeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_consume_total",
			Help: "Total count of Kafka consume.",
		},
		[]string{"topic"})

	KafkaConsumeErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_consume_error_total",
			Help: "Total number of Kafka consume errors.",
		},
		[]string{"topic"})

	KafkaConsumeDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_consume_duration_seconds",
			Help:    "Kafka consume duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"})
)
