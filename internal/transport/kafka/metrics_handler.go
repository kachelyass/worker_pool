package kafka

import (
	"context"
	"time"
	"worker_pool/internal/infrastructure/kafka"
	"worker_pool/pkg/metrics/kafka/consumermetrics"
)

type metricsHandler struct {
	topic string
	next  kafka.MessageHandler
}

func newMetricsHandler(topic string, next kafka.MessageHandler) kafka.MessageHandler {
	return &metricsHandler{
		topic: topic,
		next:  next,
	}
}

func (h *metricsHandler) Handle(ctx context.Context, key []byte, value []byte) error {
	start := time.Now()
	defer func() {
		consumermetrics.KafkaConsumeDuration.WithLabelValues(h.topic).
			Observe(time.Since(start).Seconds())
	}()

	err := h.next.Handle(ctx, key, value)
	if err != nil {
		consumermetrics.KafkaConsumeErrorsTotal.WithLabelValues(h.topic).Inc()
		return err
	}

	consumermetrics.KafkaConsumeTotal.WithLabelValues(h.topic).Inc()
	return nil
}
