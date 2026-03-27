package kafka

import "github.com/IBM/sarama"

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
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(value),
	}
	if key != "" {
		msg.Key = sarama.StringEncoder(key)
	}
	_, _, err := p.producer.SendMessage(msg)
	return err
}
