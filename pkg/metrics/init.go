package metrics

import (
	"worker_pool/pkg/metrics/httpm"

	"github.com/prometheus/client_golang/prometheus"
)

func Init() {
	prometheus.MustRegister(httpm.HTTPRequestsTotal)
	prometheus.MustRegister(httpm.HTTPRequestDuration)
	prometheus.MustRegister(httpm.HTTPInFlight)
}
