package nsq

import (
	"github.com/ebar-go/ego/utils/number"
	"github.com/nsqio/go-nsq"
	"time"
)

type Client struct {
	Address string
}

type Param struct {
	Topic               string
	Channel             string
	Debug               bool
	LookupdPollInterval int // reconnect time interval
}

// NewClient
func NewClient(address string) *Client {
	return &Client{Address: address}
}

// Listen make topic listener
func (client Client) Listen(param Param, handler nsq.Handler) error {
	cfg := nsq.NewConfig()

	cfg.LookupdPollInterval = time.Second * time.Duration(number.DefaultInt(param.LookupdPollInterval, 3)) //设置重连时间
	c, err := nsq.NewConsumer(param.Topic, param.Channel, cfg)                                             // 新建一个消费者
	if err != nil {
		return err
	}
	if !param.Debug {
		c.SetLogger(nil, 0) //屏蔽系统日志
	}

	c.AddHandler(handler) // 添加消费者接口

	// 建立一个nsqd连接
	return c.ConnectToNSQD(client.Address)
}

// NewProducer return nsq producer instance with error
func (client Client) NewProducer() (*nsq.Producer, error) {
	return nsq.NewProducer(client.Address, nsq.NewConfig())
}
