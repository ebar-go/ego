package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

type Client struct {
	host        string // kafka地址,一般是: $host:9092
	*kafka.Conn        //
}

// Connect 连接kafka服务
func (client *Client) Connect(topic string) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", client.host, topic, 0)
	//conn, err := kafka.Dial("tcp", client.host)
	if err != nil {
		return fmt.Errorf("failed to dial learder: %v", err)
	}

	//_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	client.Conn = conn
	return nil
}

func (client *Client) NewProducer() *Producer {
	return &Producer{client: client}
}

func (client *Client) NewConsumer(groupId string, callback ConsumerCallback, topics ...string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{client.host},
		GroupID:        groupId,
		GroupTopics:    topics,
		Partition:      0,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second * 2,
	})

	return &Consumer{reader: reader, callback: callback}
}

// NewClient 实例化
func NewClient(host string) *Client {
	return &Client{host: host}
}
