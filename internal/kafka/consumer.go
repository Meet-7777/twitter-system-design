package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokerAddr string, topic string, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddr},
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})
	return &Consumer{reader: reader}
}

func (c *Consumer) Consume(ctx context.Context, handler func([]byte) error) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("kafka fetch error:", err)
			time.Sleep(time.Second)
			continue
		}

		if err := handler(msg.Value); err != nil {
			log.Println("handler error:", err)
			continue
		}

		c.reader.CommitMessages(ctx, msg)
	}
}

func (c *Consumer) Close() {
	c.reader.Close()
}
