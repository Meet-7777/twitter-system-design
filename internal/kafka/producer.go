package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerAddr string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokerAddr),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  5,
		WriteTimeout: 10 * time.Second,
	}
	return &Producer{
		writer: writer,
	}
}

func (p *Producer) Publish(ctx context.Context, event any) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

func (p *Producer) Close() {
	p.writer.Close()
}
