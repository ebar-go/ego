package kafka

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	addresses := []string{"172.27.90.179:9093", "172.27.90.179:9094", "172.27.90.179:9095"}
	topic := "test-topic"

	consumer := NewConsumer(addresses, "consumer-a", topic)

	go consumer.Polling(func(topic string, msg []byte) {
		log.Printf("[%s]receive: %s\n", topic, string(msg))
	})

	producer := NewProducer(addresses, topic)

	for {
		err := producer.Write(context.Background(), []byte(uuid.NewV4().String()))
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 2)
	}

}
