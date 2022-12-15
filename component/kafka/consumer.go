package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

type ConsumerCallback func(topic string, val []byte)
type Consumer struct {
	reader *kafka.Reader
	done   chan struct{}
	once   sync.Once
}

func NewConsumer(addresses []string, groupId string, topics ...string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addresses,
		GroupID:        groupId,
		GroupTopics:    topics,
		Partition:      0,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second * 2,
	})

	return &Consumer{reader: reader}
}

// Polling 接收
func (consumer *Consumer) Polling(callback ConsumerCallback) {
	if callback == nil {
		panic("polling with nil callback")
	}
	// 保持监听
	for {
		select {
		case <-consumer.done:
			return
		default:
		}

		ctx := context.Background()
		message, err := consumer.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("unable to read message:", err)
			return
		}

		callback(message.Topic, message.Value)
	}

}

func (consumer *Consumer) Stop() {
	consumer.once.Do(func() {
		close(consumer.done)
	})
}
