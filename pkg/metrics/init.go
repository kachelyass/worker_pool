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
}
