package mns

import(
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"fmt"
	"github.com/gogap/logs"
	"github.com/ebar-go/ego/library"
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
	instance ali_mns.MNSClient
	queueItems map[string]*Queue
}

// InitClient 初始化客户端
func InitClient(conf Conf) *Client {
	if client == nil {
		client = &Client{}
		client.queueItems = make(map[string]*Queue)
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

		item.ReceiveMessage(30)
	}
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

// Queue 队列结构体
type Queue struct {
	Name string // 队列名称
	instance ali_mns.AliMNSQueue // 队列实例
	Handler QueueHandler // 处理方式
}

// QueueHandler 队列消息的处理器
type QueueHandler func(messageBody string) error

// SetHandler 设置队列消息处理器
func (queue *Queue) SetHandler(handler QueueHandler)  {
	queue.Handler = handler
}

// SendMessage 发送消息
func (queue *Queue) SendMessage(message string) (ali_mns.MessageSendResponse, error) {
	msg := ali_mns.MessageSendRequest{
		MessageBody:  message,
		DelaySeconds: 0,
		Priority:     8}

	resp, err := queue.instance.SendMessage(msg)
	return resp, err
}

// ReceiveMessage 接收消息并处理
func (queue *Queue) ReceiveMessage(waitSeconds int64) {
	if waitSeconds == 0 {
		waitSeconds = 30
	}
	endChan := make(chan int)
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		select {
		case resp := <-respChan:
			{
				logs.Pretty("response:", resp)

				if err := queue.Handler(resp.MessageBody); err != nil {
					library.Debug(err)

					// TODO ChangeMessageVisibility
				}else {
					// 处理成功，删除消息
					if err := queue.DeleteMessage(resp.ReceiptHandle); err != nil {
						fmt.Println(err)
					}
					endChan <- 1
				}

				// TODO 写日志


			}
		case err := <-errChan:
			{
				library.Debug(err)
				endChan <- 1
			}
		}
	}()

	// 通过chan去接收数据
	queue.instance.ReceiveMessage(respChan, errChan, waitSeconds)
	<-endChan
}

// DeleteMessage 删除消息
func (queue *Queue) DeleteMessage(receiptHandler string ) error{
	library.Debug(receiptHandler)
	return queue.instance.DeleteMessage(receiptHandler)
}
