package producermetrics

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaProduceTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_produce_total",
			Help: "Total number of successfully produced Kafka messages.",
		},
		[]string{"topic"},
	)

	KafkaProduceErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_produce_errors_total",
			Help: "Total number of failed Kafka produce operations.",
		},
		[]string{"topic"},
	)

	KafkaProduceDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kafka_produce_duration_seconds",
			Help:    "Kafka produce duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"topic"},
	)
)
