package kafka

import (
	"encoding/json"
	"fmt"
	"time"
	"worker_pool/internal/handlers/models"
	"worker_pool/pkg/metrics"

	"github.com/IBM/sarama"
)

type Producer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewProducer(brokers []string, clientID string, topic string) (*Producer, error) {
	scfg := sarama.NewConfig()
	scfg.ClientID = clientID
	scfg.Producer.RequiredAcks = sarama.WaitForAll
	scfg.Producer.Return.Errors = true
	scfg.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(brokers, scfg)
	if err != nil {
		return nil, err
	}

	return &Producer{
		topic:    topic,
		producer: p,
	}, nil
}

func (p *Producer) Close() error {

	return p.producer.Close()
}

func (p *Producer) Publish(key string, value []byte) error {
	start := time.Now()
	defer func() {
		metrics.KafkaProduceDuration.WithLabelValues(p.topic).Observe(time.Since(start).Seconds())
	}()
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(value),
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		metrics.KafkaProduceErrorsTotal.WithLabelValues(p.topic).Inc()
		metrics.KafkaProduceDuration.WithLabelValues(p.topic).Observe(time.Since(start).Seconds())
		return fmt.Errorf("error sending message to kafka: %w", err)
	}
	metrics.KafkaProduceTotal.WithLabelValues(p.topic).Inc()
	metrics.KafkaProduceDuration.WithLabelValues(p.topic).Observe(time.Since(start).Seconds())
	fmt.Printf("message sent to partition %d at offset %d\n", partition, offset)
	return err
}

func (p *Producer) PublishCreateTask(description string) error {
	msg := models.CreateTaskMessage{
		Description: description,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.Publish("create_task", data)

}
