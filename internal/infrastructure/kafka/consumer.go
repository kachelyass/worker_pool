package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"
	"worker_pool/pkg/metrics"

	"github.com/IBM/sarama"
)

type MessageHandler interface {
	Handle(ctx context.Context, key []byte, value []byte) error
}

type Consumer struct {
	group   sarama.ConsumerGroup
	topics  []string
	handler MessageHandler
}

func NewConsumer(
	brokers []string,
	groupID string,
	clientID string,
	topics []string,
	handler MessageHandler,
) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.ClientID = clientID
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate sarama config: %w", err)
	}

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, fmt.Errorf("create consumer group: %w", err)
	}

	return &Consumer{
		group:   group,
		topics:  topics,
		handler: handler,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context) error {
	h := &consumerGroupHandler{
		handler: c.handler,
	}

	for {
		if err := c.group.Consume(ctx, c.topics, h); err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("consume messages: %w", err)
		}

		if ctx.Err() != nil {
			return nil
		}
	}
}

func (c *Consumer) Close() error {
	return c.group.Close()
}

type consumerGroupHandler struct {
	handler MessageHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	metrics.KafkaConsumerRebalancesTotal.Inc()
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		start := time.Now()

		err := h.handler.Handle(session.Context(), msg.Key, msg.Value)
		metrics.KafkaConsumeDuration.WithLabelValues(msg.Topic).Observe(time.Since(start).Seconds())

		if err != nil {
			metrics.KafkaConsumeErrorsTotal.WithLabelValues(msg.Topic).Inc()
			return err
		}

		session.MarkMessage(msg, "")
		metrics.KafkaConsumeTotal.WithLabelValues(msg.Topic).Inc()
	}

	return nil
}
