package consumermetrics

import "github.com/prometheus/client_golang/prometheus"

func Init() {
	prometheus.MustRegister(KafkaConsumeTotal)
	prometheus.MustRegister(KafkaConsumeDuration)
	prometheus.MustRegister(KafkaConsumeErrorsTotal)
	prometheus.MustRegister(KafkaConsumeDuration)
}
