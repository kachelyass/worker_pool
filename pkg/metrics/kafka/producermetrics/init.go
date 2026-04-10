package producermetrics

import "github.com/prometheus/client_golang/prometheus"

func Init() {
	prometheus.MustRegister(KafkaProduceTotal)
	prometheus.MustRegister(KafkaProduceErrorsTotal)
	prometheus.MustRegister(KafkaProduceDuration)
}
