package nsq

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	topic = "nsq-test"
	channel = "channel-test"
	address = "192.168.0.222:4150"
)

type Consumer struct {
}

//处理消息
func (*Consumer) HandleMessage(msg *nsq.Message) error {
	fmt.Println("receive", msg.NSQDAddress, "message:", string(msg.Body))
	return nil
}
func TestPublish(t *testing.T) {
	var err error

	client := Client{Address:address}
	producer, err := client.NewProducer()
	assert.Nil(t, err)

	err = producer.Publish(topic, []byte("hello, world"))
	assert.Nil(t, err)
}

func TestReceive(t *testing.T)  {
	client := Client{Address:address}
	err := client.Listen(Param{
		Topic:               topic,
		Channel:             channel,
		Debug:               false,
		LookupdPollInterval: 3,
	}, &Consumer{})
	assert.Nil(t, err)
	for {
		fmt.Println("running...")
		time.Sleep(time.Second * 10)
	}
}
