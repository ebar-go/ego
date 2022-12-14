package kafka

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	addresses := []string{"172.27.93.155:9093"}
	topic := "test-topic"

	consumer := NewConsumer(addresses, "consumer-a", func(topic string, msg []byte) {
		log.Printf("[%s]receive: %s\n", topic, string(msg))
	}, topic)

	go consumer.Receive()

	producer := NewProducer(addresses, topic)

	for {
		err := producer.Write(context.Background(), []byte(uuid.NewV4().String()))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 2)
	}

}
