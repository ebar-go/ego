package mns

import (
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/ebar-go/ego/helper"
)

// Conf 阿里云MNS 配置项
type Conf struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

type IClient interface {
	// 生成params的sign字段
	GenerateSign(str string) string

	// 添加主题
	AddTopic(name string)

	// 获取主题
	GetTopic(name string) Topic

	// 添加队列
	AddQueue(name string, handler QueueHandler, waitSecond int)

	// 获取队列
	GetQueue(name string) Queue

	// 监听队列
	ListenQueues()
}

// Client MNS客户端
type Client struct {
	// 配置
	conf Conf

	// 阿里云mns实例
	instance ali_mns.MNSClient

	// 队列
	queueItems map[string]Queue

	// 主题
	topicItems map[string]Topic
}

// NewClient 实例化
func NewClient(conf Conf) IClient {
	cli := new(Client)
	cli.conf = conf
	cli.instance = ali_mns.NewAliMNSClient(conf.Url,
		conf.AccessKeyId,
		conf.AccessKeySecret)
	cli.queueItems = make(map[string]Queue)
	cli.topicItems = make(map[string]Topic)

	return cli
}

// GenerateSign 生成签名
func (cli *Client) GenerateSign(str string) string {
	return  helper.GetMd5String(str + cli.conf.AccessKeySecret)
}

// NewTopic 实例化主题
func (cli *Client) AddTopic(name string) {
	t := &topic{Name:name, instance:ali_mns.NewMNSTopic(name, cli.instance)}
	cli.topicItems[name] = t
}

// AddQueue 实例化队列
func (cli *Client) AddQueue(name string, handler QueueHandler, waitSecond int) {
	q := new(queue)
	q.Name = name
	q.handler = handler
	q.WaitSecond = waitSecond
	q.instance = ali_mns.NewMNSQueue(name, cli.instance)
	cli.queueItems[name] = q
}

// GetTopic 获取主题
func (cli *Client) GetTopic(name string) Topic {
	if _ , ok := cli.topicItems[name]; !ok {
		cli.AddTopic(name)
	}

	return cli.topicItems[name]
}

// GetQueue 获取队列
func (cli *Client) GetQueue(name string) Queue{
	return cli.queueItems[name]
}

// ListenQueues 监听队列
func (cli *Client) ListenQueues() {
	if len(cli.queueItems) == 0 {
		return
	}

	for _, item := range cli.queueItems {
		if item.HasHandler() {
			go item.ReceiveMessage()
		}
	}
}




