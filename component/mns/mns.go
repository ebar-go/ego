package mns

import (
	"github.com/aliyun/aliyun-mns-go-sdk"
)

// Conf 阿里云MNS 配置项
type Conf struct {
	Url             string `json:"url"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

var client *Client

// Client MNS客户端
type Client struct {
	conf Conf
	instance ali_mns.MNSClient
	queueItems map[string]*Queue
	topicItems map[string]*Topic
}


// GetTopic 获取主体
func (client *Client) GetTopic(name string) *Topic {
	if client.topicItems == nil {
		client.topicItems = make(map[string]*Topic)
	}

	if _ , ok := client.topicItems[name]; !ok {
		client.topicItems[name] = &Topic{
			Name: name,
			instance: ali_mns.NewMNSTopic(name, client.instance),
		}
	}

	return client.topicItems[name]
}



// InitClient 初始化客户端
func InitClient(conf Conf) *Client {
	if client == nil {
		client = &Client{
			conf:conf,
		}
		client.queueItems = make(map[string]*Queue)
		client.topicItems = make(map[string]*Topic)
	}

	client.instance = ali_mns.NewAliMNSClient(conf.Url,
		conf.AccessKeyId,
		conf.AccessKeySecret)

	return client
}

// GetClient 获取客户端
func GetClient() *Client {
	return client
}

// ListenQueues 监听队列
func (client *Client) ListenQueues() {
	if len(client.queueItems) == 0 {
		return
	}

	for _, item := range client.queueItems {
		if item.Handler == nil {
			continue
		}

		go item.ReceiveMessage(int64(item.WaitSecond))
	}
}

// AddQueue 添加队列
func (client *Client) AddQueue(queue *Queue) {
	queue.instance = ali_mns.NewMNSQueue(queue.Name, client.instance)
	client.queueItems[queue.Name] = queue
}

// GetQueue 获取队列
func (client *Client) GetQueue(name string) *Queue{
	if client.queueItems == nil {
		client.queueItems = make(map[string]*Queue)
	}

	if _ , ok := client.queueItems[name]; !ok {
		client.queueItems[name] = &Queue{
			Name: name,
			instance: ali_mns.NewMNSQueue(name, client.instance),
			Handler: nil,
		}
	}

	return client.queueItems[name]
}


