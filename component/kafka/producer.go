package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(addresses []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(addresses...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}
	return &Producer{writer: writer}
}

func (p *Producer) Write(ctx context.Context, msg []byte) error {
	return p.WriteMessages(ctx, msg)
}

func (p *Producer) WriteMessages(ctx context.Context, msg ...[]byte) error {
	messages := make([]kafka.Message, len(msg))
	for i, item := range msg {
		messages[i] = kafka.Message{Value: item}
	}
	return p.writer.WriteMessages(ctx, messages...)
}
