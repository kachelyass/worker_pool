package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	WorkerPoolWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_pool_workers",
			Help: "Current number of active workers in the pool.",
		},
	)

	WorkerPoolQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "worker_pool_queue_size",
			Help: "Current number of jobs waiting in the queue.",
		},
	)

	WorkerPoolJobsEnqueuedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_pool_jobs_enqueued_total",
			Help: "Total number of jobs enqueued into the worker pool.",
		},
	)

	WorkerPoolJobsStartedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_pool_jobs_started_total",
			Help: "Total number of jobs started by workers.",
		},
	)

	WorkerPoolJobsCompletedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_pool_jobs_completed_total",
			Help: "Total number of successfully completed jobs.",
		},
	)

	WorkerPoolJobsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_pool_jobs_failed_total",
			Help: "Total number of failed jobs.",
		},
	)

	WorkerPoolJobDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "worker_pool_job_duration_seconds",
			Help:    "Worker job processing duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
	)
)
