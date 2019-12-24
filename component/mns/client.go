package mns

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/ebar-go/ego/component/log"
	"github.com/ebar-go/ego/component/trace"
	"github.com/ebar-go/ego/helper"
)

// Client mns客户端接口
type Client interface {
	// 生成params的sign字段
	GenerateSign(str string) string

	// 添加队列
	AddQueue(name string, handler QueueHandler, waitSecond int)

	// 获取队列
	GetQueue(name string) Queue

	// 监听队列
	ListenQueues()

	PublishMessage(topicName string, params Params, filterTag string) (*ali_mns.MessageSendResponse, error)
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

	// log manager
	logManager log.Manager
}

// NewClient 实例化
func NewClient(url, accessKeyId, accessKeySecret string, manager log.Manager) Client {
	cli := new(client)
	cli.accessKeySecret = accessKeySecret
	cli.instance = ali_mns.NewAliMNSClient(url, accessKeyId, accessKeySecret)
	cli.queueItems = make(map[string]Queue)
	cli.topicItems = make(map[string]Topic)
	cli.logManager = manager

	return cli
}

// GenerateSign 生成签名
func (cli *client) GenerateSign(str string) string {
	return helper.GetMd5String(str + cli.accessKeySecret)
}

// AddQueue 实例化队列
func (cli *client) AddQueue(name string, handler QueueHandler, waitSecond int) {
	q := Queue{}
	q.Name = name
	q.handler = handler
	q.WaitSecond = waitSecond
	q.instance = ali_mns.NewMNSQueue(name, cli.instance)
	cli.queueItems[name] = q
}

// GetTopic 获取主题
func (cli *client) getTopic(name string) Topic {
	if _, ok := cli.topicItems[name]; !ok {
		cli.topicItems[name] = Topic{Name: name, Instance: ali_mns.NewMNSTopic(name, cli.instance)}
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
			go cli.ReceiveMessage(item.Name)
		}
	}
}

// ReceiveMessage 接收消息并处理
func (cli *client) ReceiveMessage(queueName string) {
	q := cli.GetQueue(queueName)
	if q.WaitSecond == 0 {
		q.WaitSecond = 30
	}
	endChan := make(chan int)
	respChan := make(chan ali_mns.MessageReceiveResponse)
	errChan := make(chan error)
	go func() {
		select {
		case resp := <-respChan:
			{
				var params Params

				// 解析消息
				if err := json.Unmarshal([]byte(helper.DecodeBase64Str(resp.MessageBody)), &params); err != nil {
					cli.logManager.System().Error("invalidMessageBody", log.Context{
						"err":   err.Error(),
						"trace": helper.Trace(),
					})
				} else {

					cli.logManager.Mq().Info("receiveMessage", log.Context{
						"receiveTime": helper.GetTimeStr(),
						"queue_name":  q.Name,
						"messageBody": params.Content,
						"tag":         params.Tag,
						"trace_id":    params.TraceId,
					})

					if err := q.handler(params); err != nil {
						cli.logManager.System().Warn("processMessageFailed", log.Context{
							"err":   err.Error(),
							"trace": helper.Trace(),
						})

					} else {
						// 处理成功，删除消息
						err := q.DeleteMessage(resp.ReceiptHandle)
						fmt.Println(err)
						cli.logManager.Mq().Info("deleteMessage", log.Context{
							"receiveTime": helper.GetTimeStr(),
							"queue_name":  q.Name,
							"messageBody": params.Content,
							"tag":         params.Tag,
							"trace_id":    params.TraceId,
							"err":         err,
						})

						endChan <- 1
					}
				}

			}
		case err := <-errChan:
			{
				cli.logManager.System().Info("receiveMessageFailed", log.Context{
					"err":   err.Error(),
					"trace": helper.Trace(),
				})
				fmt.Println(err)
				endChan <- 1
			}
		}
	}()

	// 通过chan去接收数据
	q.instance.ReceiveMessage(respChan, errChan, int64(q.WaitSecond))
	<-endChan
}

// PublishMessage 发布消息
func (cli *client) PublishMessage(topicName string, params Params, filterTag string) (*ali_mns.MessageSendResponse, error) {
	params.TraceId = helper.DefaultString(params.TraceId, trace.GetTraceId())
	params.Sign = helper.DefaultString(params.Sign, cli.GenerateSign(params.TraceId))
	bytes, err := helper.JsonEncode(params)
	if err != nil {
		return nil, err
	}

	topic := cli.getTopic(topicName)
	request := ali_mns.MessagePublishRequest{
		MessageBody: base64.StdEncoding.EncodeToString([]byte(bytes)),
		MessageTag:  filterTag,
	}
	resp, err := topic.Instance.PublishMessage(request)
	if err != nil {
		return nil, err
	}

	cli.logManager.Mq().Info("publishMessage", log.Context{
		"action":          "publishMessage",
		"publish_time":    helper.GetTimeStr(),
		"msectime":        helper.GetTimeStampFloatStr(),
		"message_id":      resp.MessageId,
		"status_code":     resp.Code,
		"topic_name":      topic.Name,
		"message_tag":     params.Tag,
		"global_trace_id": helper.NewTraceId(),
		"trace_id":        params.TraceId,
		"filter_tag":      filterTag,
		"sign":            params.Sign,
	})

	return &resp, nil
}
