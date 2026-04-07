package httpm

import "github.com/prometheus/client_golang/prometheus"

func Init() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPRequestDuration)
	prometheus.MustRegister(HTTPInFlight)
}
