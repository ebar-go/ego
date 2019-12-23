package mns

import (
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/ebar-go/ego/helper"
)

// Client mns客户端接口
type Client interface {
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
type client struct {
	// 配置
	accessKeySecret string

	// 阿里云mns实例
	instance ali_mns.MNSClient

	// 队列
	queueItems map[string]Queue

	// 主题
	topicItems map[string]Topic
}

// NewClient 实例化
func NewClient(url, accessKeyId, accessKeySecret string) Client {
	cli := new(client)
	cli.accessKeySecret = accessKeySecret
	cli.instance = ali_mns.NewAliMNSClient(url, accessKeyId, accessKeySecret)
	cli.queueItems = make(map[string]Queue)
	cli.topicItems = make(map[string]Topic)

	return cli
}

// GenerateSign 生成签名
func (cli *client) GenerateSign(str string) string {
	return helper.GetMd5String(str + cli.accessKeySecret)
}

// NewTopic 实例化主题
func (cli *client) AddTopic(name string) {
	t := &topic{Name: name, instance: ali_mns.NewMNSTopic(name, cli.instance)}
	cli.topicItems[name] = t
}

// AddQueue 实例化队列
func (cli *client) AddQueue(name string, handler QueueHandler, waitSecond int) {
	q := new(queue)
	q.Name = name
	q.handler = handler
	q.WaitSecond = waitSecond
	q.instance = ali_mns.NewMNSQueue(name, cli.instance)
	cli.queueItems[name] = q
}

// GetTopic 获取主题
func (cli *client) GetTopic(name string) Topic {
	if _, ok := cli.topicItems[name]; !ok {
		cli.AddTopic(name)
	}

	return cli.topicItems[name]
}

// GetQueue 获取队列
func (cli *client) GetQueue(name string) Queue {
	return cli.queueItems[name]
}

// ListenQueues 监听队列
func (cli *client) ListenQueues() {
	if len(cli.queueItems) == 0 {
		return
	}

	for _, item := range cli.queueItems {
		if item.HasHandler() {
			go item.ReceiveMessage()
		}
	}
}
